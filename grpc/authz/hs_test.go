package authz

import (
	"context"
	"encoding/json"
	"testing"

	jwt "github.com/cristalhq/jwt/v3"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_testing "github.com/grpc-ecosystem/go-grpc-middleware/testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	hsKey = []byte("Hello123$(ASM@_ASJ@@#)WR)SADJ@#T(Q#$")
)

type AuthzHSTestSuite struct {
	*grpc_testing.InterceptorTestSuite
}

func createTokenHS(scopes []string, algo jwt.Algorithm, key []byte, audience string) (*jwt.Token, error) {
	signer, err := jwt.NewSignerHS(algo, key)
	if err != nil {
		return nil, err
	}
	builder := jwt.NewBuilder(signer)

	claims := &userClaims{
		StandardClaims: jwt.StandardClaims{
			Audience: []string{audience},
			ID:       "random-unique-string",
		},
		Scopes: scopes,
	}
	return builder.Build(claims)
}

func TestHSTestSuite(t *testing.T) {
	verifyScopeFn := func(tokenClaims []byte, scopes []string) error {
		type CustomClaims struct {
			jwt.StandardClaims
			Scopes []string `json:"scopes"`
		}
		var newClaims CustomClaims
		err := json.Unmarshal(tokenClaims, &newClaims)
		if err != nil {
			return unauthenticatedError(err.Error())
		}

		if !contains(newClaims.Scopes, "read") {
			return permissionDeniedError("no matching scope found")
		}

		return nil
	}

	a := &Config{
		Algorithm: jwt.HS256,
		Audience:  "admin",
		ScopeMapping: map[string][]string{
			"/mwitkow.testproto.TestService/Ping":      {"read"},
			"/mwitkow.testproto.TestService/PingEmpty": {"read"},
		},
		ValidateScopeFunc: verifyScopeFn,
		HSKey:             hsKey,
	}

	s := &AuthzHSTestSuite{
		InterceptorTestSuite: &grpc_testing.InterceptorTestSuite{
			TestService: &assertingPingService{&grpc_testing.TestPingService{T: t}, t},
			ServerOpts: []grpc.ServerOption{
				grpc.StreamInterceptor(grpc_auth.StreamServerInterceptor(a.AuthFunc)),
				grpc.UnaryInterceptor(grpc_auth.UnaryServerInterceptor(a.AuthFunc)),
			},
		},
	}
	suite.Run(t, s)
}

func (s *AuthzHSTestSuite) TestUnary_AuthHS_Passes() {
	sc := []string{"read", "write"}
	tk, _ := createTokenHS(sc, jwt.HS256, hsKey, "admin")
	_, err := s.Client.Ping(ctxWithToken(s.SimpleCtx(), "bearer", tk.String()), goodPing)
	if err != nil {
		s.Suite.Fail(err.Error())
	}
}

func (s *AuthzHSTestSuite) TestUnary_AuthHS_VerifyError() {
	_, err := s.Client.Ping(ctxWithToken(s.SimpleCtx(), "bearer", "bad_token"), goodPing)
	assert.Error(s.T(), err, "there must be an error")
	assert.Equal(s.T(), codes.Unauthenticated, status.Code(err), "token format is not valid")
}

func (s *AuthzHSTestSuite) TestUnary_AuthHS_BadAudience_Error() {
	sc := []string{"read", "write"}
	tk, _ := createTokenHS(sc, jwt.HS256, hsKey, "user")
	_, err := s.Client.Ping(ctxWithToken(s.SimpleCtx(), "bearer", tk.String()), goodPing)
	assert.Error(s.T(), err, "there must be an error")
	assert.Equal(s.T(), codes.Unauthenticated, status.Code(err), "not for audience")
}

func (s *AuthzHSTestSuite) TestUnary_AuthHS_EndpointNotProtectedByAuth() {
	sc := []string{"read", "write"}
	tk, _ := createTokenHS(sc, jwt.HS256, hsKey, "admin")
	_, err := s.Client.PingEmpty(ctxWithToken(s.SimpleCtx(), "bearer", tk.String()), emptyPing)
	if err != nil {
		s.Suite.Fail(err.Error())
	}
}
func (s *AuthzHSTestSuite) TestUnary_AuthHS_NoTokenError() {
	_, err := s.Client.PingEmpty(context.Background(), emptyPing)
	assert.Error(s.T(), err, "there must be an error")
	assert.Equal(s.T(), codes.Unauthenticated, status.Code(err), "token format is not valid")
}
