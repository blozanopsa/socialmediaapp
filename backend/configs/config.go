package configs

import (
	"os"
)

type Config struct {
	DBUser                 string
	DBPassword             string
	DBHost                 string
	DBPort                 string
	DBName                 string
	OAuthClientID          string
	OAuthClientSecret      string
	OAuthRedirectURL       string
	OAuthMicrosoftTenantID string
	FacebookAppID          string
	FacebookAppSecret      string
}

func LoadConfig() *Config {
	return &Config{
		DBUser:                 os.Getenv("DB_USER"),
		DBPassword:             os.Getenv("DB_PASSWORD"),
		DBHost:                 os.Getenv("DB_HOST"),
		DBPort:                 os.Getenv("DB_PORT"),
		DBName:                 os.Getenv("DB_NAME"),
		OAuthClientID:          os.Getenv("OAUTH_CLIENT_ID"),
		OAuthClientSecret:      os.Getenv("OAUTH_CLIENT_SECRET"),
		OAuthRedirectURL:       os.Getenv("OAUTH_REDIRECT_URL"),
		OAuthMicrosoftTenantID: os.Getenv("OAUTH_MICROSOFT_TENANT_ID"),
		FacebookAppID:          os.Getenv("FACEBOOK_APP_ID"),
		FacebookAppSecret:      os.Getenv("FACEBOOK_APP_SECRET"),
	}
}
