package authz

import (
	"context"
	"encoding/json"
	"time"

	jwt "github.com/cristalhq/jwt/v3"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const authorizationType = "bearer"

// Base auth details
type Base struct {
	Algorithm jwt.Algorithm
	// scopes will be populated from the token that
	// comes with the request
	scopes []string
	// ScopeMapping should hold full rpc method, given by
	// grpc.UnaryServerInfo.FullMethod to allowed scopes
	ScopeMapping map[string][]string
	// ValidateScopeFunc is a user defined func for validating a token
	// has the correct scopes. This will take in the decoded token json and
	// unmarshal into any struct the user wants. See the test files for examples.
	ValidateScopeFunc         func(tokenClaims []byte, scopes []string) error
	Audience                  string
	DisableAudienceValidation bool
}

func unauthenticatedError(msg string) error {
	return status.Errorf(codes.Unauthenticated, "access token is invalid: %s", msg)
}

func permissionDeniedError(msg string) error {
	return status.Errorf(codes.PermissionDenied, "no permission to access this RPC %s", msg)
}

func (b *Base) doProtected(ctx context.Context) (string, error) {
	token, err := grpc_auth.AuthFromMD(ctx, authorizationType)
	if err != nil {
		return token, err
	}
	fullMethodName, _ := grpc.Method(ctx)
	var protected bool
	b.scopes, protected = b.ScopeMapping[fullMethodName]
	if !protected {
		return token, nil
	}
	return token, nil
}

func (b *Base) doVerify(ctx context.Context, token string, verifier jwt.Verifier) ([]byte, error) {
	newToken, err := jwt.ParseAndVerifyString(token, verifier)
	if err != nil {
		return nil, unauthenticatedError(err.Error())
	}

	var newClaims jwt.StandardClaims
	err = json.Unmarshal(newToken.RawClaims(), &newClaims)
	if err != nil {
		return nil, unauthenticatedError(err.Error())
	}

	// Perform standard JWT validations
	if !newClaims.IsValidAt(time.Now()) {
		return nil, unauthenticatedError("access token is invalid: not valid")
	}
	if !newClaims.IsValidExpiresAt(time.Now()) {
		return nil, unauthenticatedError("access token is invalid: expired")
	}

	// Perform audience claim validation
	if !b.DisableAudienceValidation {
		if !newClaims.IsForAudience(b.Audience) {
			return nil, unauthenticatedError("not for audience")
		}
	}

	return newToken.RawClaims(), nil
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
