package authz

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"time"

	jwt "github.com/cristalhq/jwt/v3"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const authorizationType = "bearer"

// Config auth details
type Config struct {
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
	// HSKey for use with HS algorithms
	HSKey []byte
	// RSAPublicKey for use with RS algorithms
	RSAPublicKey *rsa.PublicKey
}

// ConfigOption for setting optional values
type ConfigOption func(*Config)

// WithScopeMapping sets the ScopeMapping option
func WithScopeMapping(scopeMap map[string][]string) ConfigOption {
	return func(args *Config) { args.ScopeMapping = scopeMap }
}

// WithValidateScopeFunc sets the ValidateScopeFunc option
func WithValidateScopeFunc(scopeFunc func(tokenClaims []byte, scopes []string) error) ConfigOption {
	return func(args *Config) { args.ValidateScopeFunc = scopeFunc }
}

// WithAudience sets the audience
func WithAudience(aud string) ConfigOption {
	return func(args *Config) { args.Audience = aud }
}

// WithDisableAudienceValidation sets the WithDisableAudienceValidation option
func WithDisableAudienceValidation(disable bool) ConfigOption {
	return func(args *Config) { args.DisableAudienceValidation = disable }
}

// WithHSKey sets the HS key
func WithHSKey(hKey []byte) ConfigOption {
	return func(args *Config) { args.HSKey = hKey }
}

// WithRSAPubKey sets the RSA public key
func WithRSAPubKey(rsaPubKey *rsa.PublicKey) ConfigOption {
	return func(args *Config) { args.RSAPublicKey = rsaPubKey }
}

// NewConfig returns a new config with options
func NewConfig(algo jwt.Algorithm, opts ...ConfigOption) *Config {
	defaultConfig := &Config{
		Algorithm:         algo,
		ValidateScopeFunc: func(tokenClaims []byte, scopes []string) error { return nil },
	}
	for _, opt := range opts {
		opt(defaultConfig)
	}
	return defaultConfig
}

// AuthFunc authorization function
func (c *Config) AuthFunc(ctx context.Context) (context.Context, error) {
	token, protected, err := c.doProtected(ctx)
	if err != nil {
		return ctx, err
	}
	if !protected {
		return ctx, nil
	}
	var verifier jwt.Verifier
	switch c.Algorithm {
	case jwt.HS256, jwt.HS384, jwt.HS512:
		verifier, err = jwt.NewVerifierHS(c.Algorithm, c.HSKey)
		if err != nil {
			return ctx, status.Errorf(codes.FailedPrecondition, "verifier error: %v", err.Error())
		}
	case jwt.RS256, jwt.RS384, jwt.RS512:
		verifier, err = jwt.NewVerifierRS(c.Algorithm, c.RSAPublicKey)
		if err != nil {
			return ctx, status.Errorf(codes.FailedPrecondition, "verifier error: %v", err.Error())
		}
	default:
		return ctx, status.Errorf(codes.Unimplemented, "algorithm is not supported: %T", c.Algorithm)
	}

	rawToken, err := c.doVerify(ctx, token, verifier)
	if err != nil {
		return ctx, err
	}
	return ctx, c.ValidateScopeFunc(rawToken, c.scopes)
}

func unauthenticatedError(msg string) error {
	return status.Errorf(codes.Unauthenticated, "access token is invalid: %s", msg)
}

func permissionDeniedError(msg string) error {
	return status.Errorf(codes.PermissionDenied, "no permission to access this RPC %s", msg)
}

// doProtected checks if the method should be protected by auth or not. If so, then scopes
// for the method are saved to Config.scopes for later use.
func (c *Config) doProtected(ctx context.Context) (token string, protected bool, err error) {
	fullMethodName, _ := grpc.Method(ctx)
	c.scopes, protected = c.ScopeMapping[fullMethodName]
	if !protected {
		return token, false, nil
	}
	token, err = grpc_auth.AuthFromMD(ctx, authorizationType)
	if err != nil {
		return token, protected, err
	}
	return token, true, nil
}

// doVerify runs some standard JWT validations against a token
func (c *Config) doVerify(ctx context.Context, token string, verifier jwt.Verifier) ([]byte, error) {
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
	if !c.DisableAudienceValidation {
		if !newClaims.IsForAudience(c.Audience) {
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
