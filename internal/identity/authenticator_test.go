package identity

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAuthenticatorValidateToken(t *testing.T) {
	privateKey := newTestKey(t)
	kid := "test-key"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"keys": []map[string]string{jwkForKey(kid, &privateKey.PublicKey)},
		})
	}))
	defer server.Close()

	auth, err := NewAuthenticator(Config{
		Enabled:  true,
		Issuer:   "https://api.asgardeo.io/t/example/oauth2/token",
		Audience: "client-id",
		JWKSURL:  server.URL,
	})
	if err != nil {
		t.Fatalf("new authenticator: %v", err)
	}
	auth.now = func() time.Time { return time.Unix(1_700_000_000, 0) }

	token := signedToken(t, privateKey, kid, map[string]any{
		"sub":   "user-123",
		"iss":   "https://api.asgardeo.io/t/example/oauth2/token",
		"aud":   "client-id",
		"exp":   auth.now().Add(time.Hour).Unix(),
		"nbf":   auth.now().Add(-time.Minute).Unix(),
		"email": "user@example.com",
		"scope": "openid profile email",
		"roles": []string{"instructor"},
	})

	principal, err := auth.ValidateToken(context.Background(), token)
	if err != nil {
		t.Fatalf("validate token: %v", err)
	}
	if principal.Subject != "user-123" || principal.Email != "user@example.com" {
		t.Fatalf("unexpected principal: %+v", principal)
	}
	if len(principal.Scopes) != 3 || principal.Scopes[2] != "email" {
		t.Fatalf("unexpected scopes: %+v", principal.Scopes)
	}
	if len(principal.Roles) != 1 || principal.Roles[0] != "instructor" {
		t.Fatalf("unexpected roles: %+v", principal.Roles)
	}
}

func TestAuthenticatorRejectsWrongAudience(t *testing.T) {
	privateKey := newTestKey(t)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"keys": []map[string]string{jwkForKey("test-key", &privateKey.PublicKey)},
		})
	}))
	defer server.Close()

	auth, err := NewAuthenticator(Config{
		Enabled:  true,
		Issuer:   "issuer",
		Audience: "expected",
		JWKSURL:  server.URL,
	})
	if err != nil {
		t.Fatalf("new authenticator: %v", err)
	}
	auth.now = func() time.Time { return time.Unix(1_700_000_000, 0) }

	token := signedToken(t, privateKey, "test-key", map[string]any{
		"sub": "user-123",
		"iss": "issuer",
		"aud": "other",
		"exp": auth.now().Add(time.Hour).Unix(),
	})

	if _, err := auth.ValidateToken(context.Background(), token); err == nil {
		t.Fatal("expected wrong audience to fail")
	}
}

func TestMiddlewareAttachesPrincipal(t *testing.T) {
	privateKey := newTestKey(t)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"keys": []map[string]string{jwkForKey("test-key", &privateKey.PublicKey)},
		})
	}))
	defer server.Close()

	auth, err := NewAuthenticator(Config{
		Enabled:  true,
		Issuer:   "issuer",
		Audience: "audience",
		JWKSURL:  server.URL,
	})
	if err != nil {
		t.Fatalf("new authenticator: %v", err)
	}
	auth.now = func() time.Time { return time.Unix(1_700_000_000, 0) }

	token := signedToken(t, privateKey, "test-key", map[string]any{
		"sub": "user-123",
		"iss": "issuer",
		"aud": []string{"another", "audience"},
		"exp": auth.now().Add(time.Hour).Unix(),
	})

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		principal, ok := PrincipalFromContext(r.Context())
		if !ok || principal.Subject != "user-123" {
			t.Fatalf("missing principal: %+v", principal)
		}
		w.WriteHeader(http.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	res := httptest.NewRecorder()

	auth.Middleware(next).ServeHTTP(res, req)

	if res.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusNoContent)
	}
}

func TestMiddlewareRejectsMissingBearerToken(t *testing.T) {
	auth, err := NewAuthenticator(Config{
		Enabled:  true,
		Issuer:   "issuer",
		Audience: "audience",
		JWKSURL:  "https://example.com/jwks",
	})
	if err != nil {
		t.Fatalf("new authenticator: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()
	auth.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("next handler should not be called")
	})).ServeHTTP(res, req)

	if res.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusUnauthorized)
	}
}

func TestMiddlewareAttachesDevPrincipalWhenAuthDisabled(t *testing.T) {
	auth, err := NewAuthenticator(Config{
		Enabled: false,
		DevPrincipal: Principal{
			Subject: "local-user",
			Email:   "local@example.com",
			Roles:   []string{"instructor"},
			Scopes:  []string{"courses:write"},
		},
	})
	if err != nil {
		t.Fatalf("new authenticator: %v", err)
	}

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		principal, ok := PrincipalFromContext(r.Context())
		if !ok {
			t.Fatal("expected dev principal in request context")
		}
		if principal.Subject != "local-user" || principal.Email != "local@example.com" {
			t.Fatalf("unexpected principal: %+v", principal)
		}
		if len(principal.Roles) != 1 || principal.Roles[0] != "instructor" {
			t.Fatalf("unexpected roles: %+v", principal.Roles)
		}
		w.WriteHeader(http.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()

	auth.Middleware(next).ServeHTTP(res, req)

	if res.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusNoContent)
	}
}

func TestMiddlewareUsesDefaultDevPrincipalWhenAuthDisabled(t *testing.T) {
	auth, err := NewAuthenticator(Config{Enabled: false})
	if err != nil {
		t.Fatalf("new authenticator: %v", err)
	}

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		principal, ok := PrincipalFromContext(r.Context())
		if !ok {
			t.Fatal("expected default dev principal in request context")
		}
		if principal.Subject != "dev-user" || principal.Email != "dev@toucan.local" {
			t.Fatalf("unexpected principal: %+v", principal)
		}
		w.WriteHeader(http.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()

	auth.Middleware(next).ServeHTTP(res, req)

	if res.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusNoContent)
	}
}

func newTestKey(t *testing.T) *rsa.PrivateKey {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	return key
}

func signedToken(t *testing.T, key *rsa.PrivateKey, kid string, claims map[string]any) string {
	t.Helper()
	header := map[string]string{
		"alg": "RS256",
		"typ": "JWT",
		"kid": kid,
	}
	headerSegment := encodeJSONSegment(t, header)
	claimsSegment := encodeJSONSegment(t, claims)
	signingInput := headerSegment + "." + claimsSegment
	hash := sha256.Sum256([]byte(signingInput))
	signature, err := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, hash[:])
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}
	return signingInput + "." + base64.RawURLEncoding.EncodeToString(signature)
}

func encodeJSONSegment(t *testing.T, value any) string {
	t.Helper()
	data, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("marshal segment: %v", err)
	}
	return base64.RawURLEncoding.EncodeToString(data)
}

func jwkForKey(kid string, key *rsa.PublicKey) map[string]string {
	return map[string]string{
		"kty": "RSA",
		"use": "sig",
		"kid": kid,
		"n":   base64.RawURLEncoding.EncodeToString(key.N.Bytes()),
		"e":   base64.RawURLEncoding.EncodeToString(big.NewInt(int64(key.E)).Bytes()),
	}
}
