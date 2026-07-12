package auth

import "github.com/coreos/go-oidc"

type ServiceInterface interface {
	ExchangeCode(code string) (*TokenResponse, error)
	VerifyToken(token string) (*oidc.IDToken, error)
}
