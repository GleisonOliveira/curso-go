package auth

import "os"

type Config struct {
	BaseURL      string
	IssuerURL    string
	TokenURI     string
	AuthURI      string
	CallbackURL  string
	ClientID     string
	ClientSecret string
}

type Claims struct {
	Email string `json:"email"`
}

func NewConfig() *Config {
	return &Config{
		BaseURL:      os.Getenv("KEYCLOAK_BASE_URL"),
		IssuerURL:    os.Getenv("KEYCLOAK_ISSUER_URL"),
		TokenURI:     os.Getenv("KEYCLOAK_TOKEN_URI"),
		AuthURI:      os.Getenv("KEYCLOAK_AUTH_URI"),
		CallbackURL:  os.Getenv("KEYCLOAK_CALLBACK_URL"),
		ClientID:     os.Getenv("KEYCLOAK_CLIENT_ID"),
		ClientSecret: os.Getenv("KEYCLOAK_CLIENT_SECRET"),
	}
}
