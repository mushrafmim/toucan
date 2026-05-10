package config

import (
	"os"
	"strings"
	"time"

	"toucan/internal/database"
	"toucan/internal/identity"
	"toucan/internal/storage"
)

type Config struct {
	HTTPAddr string
	Database database.Config
	Identity identity.Config
	Storage  storage.Config
}

func Load() Config {
	return Config{
		HTTPAddr: getEnv("TOUCAN_HTTP_ADDR", ":8080"),
		Database: database.Config{
			DSN: getEnv("TOUCAN_POSTGRES_DSN", "postgres://postgres:postgres@localhost:5432/toucan?sslmode=disable"),
		},
		Identity: identity.Config{
			Enabled:      getEnvBool("TOUCAN_AUTH_ENABLED", false),
			Issuer:       os.Getenv("TOUCAN_AUTH_ISSUER"),
			Audience:     os.Getenv("TOUCAN_AUTH_AUDIENCE"),
			JWKSURL:      os.Getenv("TOUCAN_AUTH_JWKS_URL"),
			JWKSCacheTTL: getEnvDuration("TOUCAN_AUTH_JWKS_CACHE_TTL", 15*time.Minute),
			DevPrincipal: identity.Principal{
				Subject: getEnv("TOUCAN_DEV_USER_SUBJECT", "dev-user"),
				Email:   getEnv("TOUCAN_DEV_USER_EMAIL", "dev@toucan.local"),
				Roles:   getEnvList("TOUCAN_DEV_USER_ROLES", []string{"admin"}),
				Scopes:  getEnvList("TOUCAN_DEV_USER_SCOPES", []string{"dev"}),
			},
		},
		Storage: storage.Config{
			Driver:    strings.ToLower(getEnv("TOUCAN_BLOB_DRIVER", storage.BlobDriverLocal)),
			LocalPath: getEnv("TOUCAN_LOCAL_STORAGE_PATH", "./uploads"),
			S3Bucket:  os.Getenv("TOUCAN_S3_BUCKET"),
			S3Region:  getEnv("TOUCAN_S3_REGION", "us-east-1"),
		},
	}
}

func getEnv(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	value := strings.TrimSpace(strings.ToLower(os.Getenv(key)))
	switch value {
	case "1", "true", "yes", "on":
		return true
	case "0", "false", "no", "off":
		return false
	default:
		return fallback
	}
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	duration, err := time.ParseDuration(value)
	if err != nil {
		return fallback
	}
	return duration
}

func getEnvList(key string, fallback []string) []string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	parts := strings.Split(value, ",")
	items := make([]string, 0, len(parts))
	for _, part := range parts {
		item := strings.TrimSpace(part)
		if item != "" {
			items = append(items, item)
		}
	}
	if len(items) == 0 {
		return fallback
	}
	return items
}
