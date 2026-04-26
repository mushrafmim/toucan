package identity

import "time"

type Config struct {
	Enabled      bool
	Issuer       string
	Audience     string
	JWKSURL      string
	JWKSCacheTTL time.Duration
	DevPrincipal Principal
}
