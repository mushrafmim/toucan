package identity

import (
	"context"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"sync"
	"time"
)

const defaultJWKSCacheTTL = 15 * time.Minute

var defaultDevPrincipal = Principal{
	Subject: "dev-user",
	Email:   "dev@toucan.local",
	Name:    "Development User",
	Roles:   []string{"admin", "instructor", "learner"},
	Scopes:  []string{"dev"},
}

var (
	ErrMissingToken = errors.New("missing bearer token")
	ErrInvalidToken = errors.New("invalid bearer token")
)

type Authenticator struct {
	cfg    Config
	client *http.Client
	now    func() time.Time

	// SyncUser is an optional callback triggered after successful authentication
	// to ensure the user record exists in the local database.
	SyncUser func(context.Context, Principal) error

	mu        sync.Mutex
	keys      map[string]*rsa.PublicKey
	fetchedAt time.Time
}

func NewAuthenticator(cfg Config) (*Authenticator, error) {
	return NewAuthenticatorWithClient(cfg, http.DefaultClient)
}

func NewAuthenticatorWithClient(cfg Config, client *http.Client) (*Authenticator, error) {
	if cfg.JWKSCacheTTL == 0 {
		cfg.JWKSCacheTTL = defaultJWKSCacheTTL
	}
	if cfg.DevPrincipal.Subject == "" {
		cfg.DevPrincipal = defaultDevPrincipal
	}
	if client == nil {
		client = http.DefaultClient
	}
	if cfg.Enabled {
		if strings.TrimSpace(cfg.Issuer) == "" {
			return nil, fmt.Errorf("identity issuer is required when auth is enabled")
		}
		if strings.TrimSpace(cfg.Audience) == "" {
			return nil, fmt.Errorf("identity audience is required when auth is enabled")
		}
		if strings.TrimSpace(cfg.JWKSURL) == "" {
			return nil, fmt.Errorf("identity jwks url is required when auth is enabled")
		}
	}

	return &Authenticator{
		cfg:    cfg,
		client: client,
		now:    time.Now,
		keys:   make(map[string]*rsa.PublicKey),
	}, nil
}

func (a *Authenticator) Middleware(next http.Handler) http.Handler {
	if !a.cfg.Enabled {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			principal := a.cfg.DevPrincipal
			if a.SyncUser != nil {
				if err := a.SyncUser(r.Context(), principal); err != nil {
					http.Error(w, "user synchronization failed", http.StatusInternalServerError)
					return
				}
			}
			ctx := ContextWithPrincipal(r.Context(), principal)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		principal, err := a.Authenticate(r)
		if err != nil {
			writeUnauthorized(w)
			return
		}

		if a.SyncUser != nil {
			if err := a.SyncUser(r.Context(), principal); err != nil {
				http.Error(w, "user synchronization failed", http.StatusInternalServerError)
				return
			}
		}

		next.ServeHTTP(w, r.WithContext(ContextWithPrincipal(r.Context(), principal)))
	})
}

func (a *Authenticator) Authenticate(r *http.Request) (Principal, error) {
	token, ok := bearerToken(r.Header.Get("Authorization"))
	if !ok {
		return Principal{}, ErrMissingToken
	}
	return a.ValidateToken(r.Context(), token)
}

func bearerToken(header string) (string, bool) {
	scheme, token, ok := strings.Cut(strings.TrimSpace(header), " ")
	if !ok || !strings.EqualFold(scheme, "Bearer") || strings.TrimSpace(token) == "" {
		return "", false
	}
	return strings.TrimSpace(token), true
}

func writeUnauthorized(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("WWW-Authenticate", `Bearer realm="toucan"`)
	w.WriteHeader(http.StatusUnauthorized)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized"})
}

func (a *Authenticator) ValidateToken(ctx context.Context, token string) (Principal, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return Principal{}, ErrInvalidToken
	}

	var header tokenHeader
	if err := decodeSegment(parts[0], &header); err != nil {
		return Principal{}, fmt.Errorf("%w: header", ErrInvalidToken)
	}
	if header.Algorithm != "RS256" || header.KeyID == "" {
		return Principal{}, ErrInvalidToken
	}

	var claims tokenClaims
	if err := decodeSegment(parts[1], &claims); err != nil {
		return Principal{}, fmt.Errorf("%w: claims", ErrInvalidToken)
	}

	signature, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return Principal{}, fmt.Errorf("%w: signature", ErrInvalidToken)
	}

	key, err := a.publicKey(ctx, header.KeyID)
	if err != nil {
		return Principal{}, err
	}
	if err := verifyRS256(parts[0]+"."+parts[1], signature, key); err != nil {
		return Principal{}, err
	}
	if err := claims.validate(a.cfg, a.now()); err != nil {
		return Principal{}, err
	}

	return claims.principal(), nil
}

func decodeSegment(segment string, out any) error {
	data, err := base64.RawURLEncoding.DecodeString(segment)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, out)
}

func verifyRS256(signingInput string, signature []byte, key *rsa.PublicKey) error {
	sum := sha256.Sum256([]byte(signingInput))
	if err := rsa.VerifyPKCS1v15(key, crypto.SHA256, sum[:], signature); err != nil {
		return fmt.Errorf("%w: signature", ErrInvalidToken)
	}
	return nil
}

func (a *Authenticator) publicKey(ctx context.Context, kid string) (*rsa.PublicKey, error) {
	key, ok := a.cachedKey(kid)
	if ok {
		return key, nil
	}
	if err := a.refreshKeys(ctx); err != nil {
		return nil, err
	}
	key, ok = a.cachedKey(kid)
	if !ok {
		return nil, fmt.Errorf("%w: unknown key", ErrInvalidToken)
	}
	return key, nil
}

func (a *Authenticator) cachedKey(kid string) (*rsa.PublicKey, bool) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.fetchedAt.IsZero() || a.now().Sub(a.fetchedAt) > a.cfg.JWKSCacheTTL {
		return nil, false
	}
	key, ok := a.keys[kid]
	return key, ok
}

func (a *Authenticator) refreshKeys(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, a.cfg.JWKSURL, nil)
	if err != nil {
		return err
	}
	res, err := a.client.Do(req)
	if err != nil {
		return fmt.Errorf("fetch jwks: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("fetch jwks: unexpected status %d", res.StatusCode)
	}

	var set jwks
	if err := json.NewDecoder(res.Body).Decode(&set); err != nil {
		return fmt.Errorf("decode jwks: %w", err)
	}

	keys := make(map[string]*rsa.PublicKey, len(set.Keys))
	for _, key := range set.Keys {
		publicKey, err := key.rsaPublicKey()
		if err != nil {
			continue
		}
		keys[key.KeyID] = publicKey
	}
	if len(keys) == 0 {
		return fmt.Errorf("decode jwks: no usable rsa keys")
	}

	a.mu.Lock()
	a.keys = keys
	a.fetchedAt = a.now()
	a.mu.Unlock()
	return nil
}

type tokenHeader struct {
	Algorithm string `json:"alg"`
	KeyID     string `json:"kid"`
}

type tokenClaims struct {
	Subject   string      `json:"sub"`
	Issuer    string      `json:"iss"`
	Audience  audience    `json:"aud"`
	ExpiresAt int64       `json:"exp"`
	NotBefore int64       `json:"nbf"`
	Email     string      `json:"email"`
	Name      string      `json:"name"`
	Scope     string      `json:"scope"`
	Scopes    []string    `json:"scp"`
	Roles     stringSlice `json:"roles"`
	Groups    stringSlice `json:"groups"`
}

func (c tokenClaims) validate(cfg Config, now time.Time) error {
	if c.Subject == "" || c.Issuer != cfg.Issuer || !c.Audience.Contains(cfg.Audience) {
		return ErrInvalidToken
	}
	unixNow := now.Unix()
	if c.ExpiresAt == 0 || unixNow >= c.ExpiresAt {
		return fmt.Errorf("%w: expired", ErrInvalidToken)
	}
	if c.NotBefore != 0 && unixNow < c.NotBefore {
		return fmt.Errorf("%w: not active", ErrInvalidToken)
	}
	return nil
}

func (c tokenClaims) principal() Principal {
	scopes := c.Scopes
	if len(scopes) == 0 && strings.TrimSpace(c.Scope) != "" {
		scopes = strings.Fields(c.Scope)
	}

	roles := []string(c.Roles)
	if len(roles) == 0 {
		roles = []string(c.Groups)
	}

	return Principal{
		Subject: c.Subject,
		Email:   c.Email,
		Name:    c.Name,
		Roles:   roles,
		Scopes:  scopes,
	}
}

type audience []string

func (a *audience) UnmarshalJSON(data []byte) error {
	var single string
	if err := json.Unmarshal(data, &single); err == nil {
		*a = []string{single}
		return nil
	}
	var many []string
	if err := json.Unmarshal(data, &many); err != nil {
		return err
	}
	*a = many
	return nil
}

func (a audience) Contains(value string) bool {
	for _, item := range a {
		if item == value {
			return true
		}
	}
	return false
}

type stringSlice []string

func (s *stringSlice) UnmarshalJSON(data []byte) error {
	var single string
	if err := json.Unmarshal(data, &single); err == nil {
		if strings.TrimSpace(single) == "" {
			*s = nil
		} else {
			*s = strings.Fields(single)
		}
		return nil
	}
	var many []string
	if err := json.Unmarshal(data, &many); err != nil {
		return err
	}
	*s = many
	return nil
}

type jwks struct {
	Keys []jwk `json:"keys"`
}

type jwk struct {
	KeyType  string `json:"kty"`
	Use      string `json:"use"`
	KeyID    string `json:"kid"`
	Modulus  string `json:"n"`
	Exponent string `json:"e"`
}

func (j jwk) rsaPublicKey() (*rsa.PublicKey, error) {
	if j.KeyType != "RSA" || j.KeyID == "" || j.Modulus == "" || j.Exponent == "" {
		return nil, fmt.Errorf("invalid rsa jwk")
	}
	if j.Use != "" && j.Use != "sig" {
		return nil, fmt.Errorf("unsupported jwk use")
	}

	modulusBytes, err := base64.RawURLEncoding.DecodeString(j.Modulus)
	if err != nil {
		return nil, err
	}
	exponentBytes, err := base64.RawURLEncoding.DecodeString(j.Exponent)
	if err != nil {
		return nil, err
	}

	exponent := int(new(big.Int).SetBytes(exponentBytes).Int64())
	if exponent == 0 {
		return nil, fmt.Errorf("invalid rsa exponent")
	}

	return &rsa.PublicKey{
		N: new(big.Int).SetBytes(modulusBytes),
		E: exponent,
	}, nil
}
