package auth

import (
	"context"

	"github.com/coreos/go-oidc"
)

type TokenVerifier interface {
	Verify(ctx context.Context, rawIDToken string) (*oidc.IDToken, error)
}
