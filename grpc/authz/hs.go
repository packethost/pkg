package authz

import (
	"context"

	jwt "github.com/cristalhq/jwt/v3"
)

// AuthWithHS is authz using HS algorithms
type AuthWithHS struct {
	Base
	Key []byte
}

// AuthFunc authorization function satisfies the go-grpc-middleware/auth AuthFunc func signature
func (a *AuthWithHS) AuthFunc(ctx context.Context) (context.Context, error) {
	token, err := a.doProtected(ctx)
	if err != nil {
		return ctx, err
	}
	verifier, err := jwt.NewVerifierHS(a.Algorithm, a.Key)
	if err != nil {
		return ctx, unauthenticatedError(err.Error())
	}
	rawToken, err := a.doVerify(ctx, token, verifier)
	if err != nil {
		return ctx, err
	}
	return ctx, a.ValidateScopeFunc(rawToken, a.scopes)
}
