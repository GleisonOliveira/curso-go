package oidc

import (
	"context"
	"fmt"

	"github.com/coreos/go-oidc"
)

func NewVerifier(issuerURL, clientID string) *oidc.IDTokenVerifier {
	provider, err := oidc.NewProvider(context.Background(), issuerURL)
	if err != nil {
		panic(fmt.Sprintf("failed to create OIDC provider: %s", err.Error()))
	}

	return provider.Verifier(&oidc.Config{ClientID: clientID})
}
