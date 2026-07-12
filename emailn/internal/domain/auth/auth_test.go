package auth

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewConfig(t *testing.T) {
	assert := assert.New(t)

	os.Setenv("KEYCLOAK_BASE_URL", "http://localhost:8080")
	os.Setenv("KEYCLOAK_ISSUER_URL", "http://localhost:8080/realms/emailn")
	os.Setenv("KEYCLOAK_TOKEN_URI", "realms/emailn/protocol/openid-connect/token")
	os.Setenv("KEYCLOAK_AUTH_URI", "realms/emailn/protocol/openid-connect/auth")
	os.Setenv("KEYCLOAK_CALLBACK_URL", "http://localhost:8081/auth/callback")
	os.Setenv("KEYCLOAK_CLIENT_ID", "emailn")
	os.Setenv("KEYCLOAK_CLIENT_SECRET", "secret123")

	defer func() {
		os.Unsetenv("KEYCLOAK_BASE_URL")
		os.Unsetenv("KEYCLOAK_ISSUER_URL")
		os.Unsetenv("KEYCLOAK_TOKEN_URI")
		os.Unsetenv("KEYCLOAK_AUTH_URI")
		os.Unsetenv("KEYCLOAK_CALLBACK_URL")
		os.Unsetenv("KEYCLOAK_CLIENT_ID")
		os.Unsetenv("KEYCLOAK_CLIENT_SECRET")
	}()

	config := NewConfig()

	assert.Equal("http://localhost:8080", config.BaseURL)
	assert.Equal("http://localhost:8080/realms/emailn", config.IssuerURL)
	assert.Equal("realms/emailn/protocol/openid-connect/token", config.TokenURI)
	assert.Equal("realms/emailn/protocol/openid-connect/auth", config.AuthURI)
	assert.Equal("http://localhost:8081/auth/callback", config.CallbackURL)
	assert.Equal("emailn", config.ClientID)
	assert.Equal("secret123", config.ClientSecret)
}

func Test_NewConfig_EmptyEnv(t *testing.T) {
	assert := assert.New(t)

	os.Unsetenv("KEYCLOAK_BASE_URL")
	os.Unsetenv("KEYCLOAK_ISSUER_URL")
	os.Unsetenv("KEYCLOAK_TOKEN_URI")
	os.Unsetenv("KEYCLOAK_AUTH_URI")
	os.Unsetenv("KEYCLOAK_CALLBACK_URL")
	os.Unsetenv("KEYCLOAK_CLIENT_ID")
	os.Unsetenv("KEYCLOAK_CLIENT_SECRET")

	config := NewConfig()

	assert.Equal("", config.BaseURL)
	assert.Equal("", config.IssuerURL)
	assert.Equal("", config.TokenURI)
	assert.Equal("", config.AuthURI)
	assert.Equal("", config.CallbackURL)
	assert.Equal("", config.ClientID)
	assert.Equal("", config.ClientSecret)
}
