package authz

import (
	"context"
	"crypto/rsa"

	jwt "github.com/cristalhq/jwt/v3"
)

// AuthWithRS is authz using RS algorithms
type AuthWithRS struct {
	Base
	RSAPublicKey *rsa.PublicKey
}

// AuthFunc authorization function
func (a *AuthWithRS) AuthFunc(ctx context.Context) (context.Context, error) {
	token, err := a.doProtected(ctx)
	if err != nil {
		return ctx, err
	}
	verifier, err := jwt.NewVerifierRS(a.Algorithm, a.RSAPublicKey)
	if err != nil {
		return ctx, unauthenticatedError(err.Error())
	}
	rawToken, err := a.doVerify(ctx, token, verifier)
	if err != nil {
		return ctx, err
	}
	return ctx, a.ValidateScopeFunc(rawToken, a.scopes)
}
