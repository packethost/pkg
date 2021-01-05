package authz

import (
	"crypto/rsa"
	"encoding/json"
	"testing"

	jwt "github.com/cristalhq/jwt/v3"
	jwt_helper "github.com/dgrijalva/jwt-go"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_testing "github.com/grpc-ecosystem/go-grpc-middleware/testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	pubKey = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAnzyis1ZjfNB0bBgKFMSv
vkTtwlvBsaJq7S5wA+kzeVOVpVWwkWdVha4s38XM/pa/yr47av7+z3VTmvDRyAHc
aT92whREFpLv9cj5lTeJSibyr/Mrm/YtjCZVWgaOYIhwrXwKLqPr/11inWsAkfIy
tvHWTxZYEcXLgAXFuUuaS3uF9gEiNQwzGTU1v0FqkqTBr4B8nW3HCN47XUu0t8Y0
e+lf4s4OxQawWD79J9/5d3Ry0vbV3Am1FtGJiJvOwRsIfVChDpYStTcHTCMqtvWb
V6L11BWkpzGXSW4Hv43qa+GSYOD2QU68Mb59oSk2OB+BtOLpJofmbGEGgvmwyCI9
MwIDAQAB
-----END PUBLIC KEY-----`
	privKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEogIBAAKCAQEAnzyis1ZjfNB0bBgKFMSvvkTtwlvBsaJq7S5wA+kzeVOVpVWw
kWdVha4s38XM/pa/yr47av7+z3VTmvDRyAHcaT92whREFpLv9cj5lTeJSibyr/Mr
m/YtjCZVWgaOYIhwrXwKLqPr/11inWsAkfIytvHWTxZYEcXLgAXFuUuaS3uF9gEi
NQwzGTU1v0FqkqTBr4B8nW3HCN47XUu0t8Y0e+lf4s4OxQawWD79J9/5d3Ry0vbV
3Am1FtGJiJvOwRsIfVChDpYStTcHTCMqtvWbV6L11BWkpzGXSW4Hv43qa+GSYOD2
QU68Mb59oSk2OB+BtOLpJofmbGEGgvmwyCI9MwIDAQABAoIBACiARq2wkltjtcjs
kFvZ7w1JAORHbEufEO1Eu27zOIlqbgyAcAl7q+/1bip4Z/x1IVES84/yTaM8p0go
amMhvgry/mS8vNi1BN2SAZEnb/7xSxbflb70bX9RHLJqKnp5GZe2jexw+wyXlwaM
+bclUCrh9e1ltH7IvUrRrQnFJfh+is1fRon9Co9Li0GwoN0x0byrrngU8Ak3Y6D9
D8GjQA4Elm94ST3izJv8iCOLSDBmzsPsXfcCUZfmTfZ5DbUDMbMxRnSo3nQeoKGC
0Lj9FkWcfmLcpGlSXTO+Ww1L7EGq+PT3NtRae1FZPwjddQ1/4V905kyQFLamAA5Y
lSpE2wkCgYEAy1OPLQcZt4NQnQzPz2SBJqQN2P5u3vXl+zNVKP8w4eBv0vWuJJF+
hkGNnSxXQrTkvDOIUddSKOzHHgSg4nY6K02ecyT0PPm/UZvtRpWrnBjcEVtHEJNp
bU9pLD5iZ0J9sbzPU/LxPmuAP2Bs8JmTn6aFRspFrP7W0s1Nmk2jsm0CgYEAyH0X
+jpoqxj4efZfkUrg5GbSEhf+dZglf0tTOA5bVg8IYwtmNk/pniLG/zI7c+GlTc9B
BwfMr59EzBq/eFMI7+LgXaVUsM/sS4Ry+yeK6SJx/otIMWtDfqxsLD8CPMCRvecC
2Pip4uSgrl0MOebl9XKp57GoaUWRWRHqwV4Y6h8CgYAZhI4mh4qZtnhKjY4TKDjx
QYufXSdLAi9v3FxmvchDwOgn4L+PRVdMwDNms2bsL0m5uPn104EzM6w1vzz1zwKz
5pTpPI0OjgWN13Tq8+PKvm/4Ga2MjgOgPWQkslulO/oMcXbPwWC3hcRdr9tcQtn9
Imf9n2spL/6EDFId+Hp/7QKBgAqlWdiXsWckdE1Fn91/NGHsc8syKvjjk1onDcw0
NvVi5vcba9oGdElJX3e9mxqUKMrw7msJJv1MX8LWyMQC5L6YNYHDfbPF1q5L4i8j
8mRex97UVokJQRRA452V2vCO6S5ETgpnad36de3MUxHgCOX3qL382Qx9/THVmbma
3YfRAoGAUxL/Eu5yvMK8SAt/dJK6FedngcM3JEFNplmtLYVLWhkIlNRGDwkg3I5K
y18Ae9n7dHVueyslrb6weq7dTkYDi3iOYRW8HRkIQh06wEdbxt0shTzAJvvCQfrB
jg/3747WSsf/zBTcHihTRBdAv6OmdhV4/dD5YBfLAkLrd+mX7iE=
-----END RSA PRIVATE KEY-----`
)

type AuthzRSTestSuite struct {
	*grpc_testing.InterceptorTestSuite
	token string
}

func createTokenRS(scopes []string, algo jwt.Algorithm, privateKey *rsa.PrivateKey, audience string) (*jwt.Token, error) {
	signer, err := jwt.NewSignerRS(algo, privateKey)
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

func TestRSSuccessful(t *testing.T) {
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

	pKey, _ := jwt_helper.ParseRSAPublicKeyFromPEM([]byte(pubKey))
	a := &Config{
		Algorithm: jwt.RS256,
		ScopeMapping: map[string][]string{
			"/mwitkow.testproto.TestService/Ping": {"read"},
		},
		Audience:          "admin",
		ValidateScopeFunc: verifyScopeFn,
		RSAPublicKey:      pKey,
	}
	privateKey, _ := jwt_helper.ParseRSAPrivateKeyFromPEM([]byte(privKey))
	sc := []string{"read", "write"}
	tk, err := createTokenRS(sc, jwt.RS256, privateKey, "admin")
	if err != nil {
		t.Log(err)
	}

	s := &AuthzRSTestSuite{
		InterceptorTestSuite: &grpc_testing.InterceptorTestSuite{
			TestService: &assertingPingService{&grpc_testing.TestPingService{T: t}, t},
			ServerOpts: []grpc.ServerOption{
				grpc.StreamInterceptor(grpc_auth.StreamServerInterceptor(a.AuthFunc)),
				grpc.UnaryInterceptor(grpc_auth.UnaryServerInterceptor(a.AuthFunc)),
			},
		},
		token: tk.String(),
	}
	suite.Run(t, s)
}

func (s *AuthzRSTestSuite) TestUnary_PassesAuthRS() {
	_, err := s.Client.Ping(ctxWithToken(s.SimpleCtx(), "bearer", s.token), goodPing)
	if err != nil {
		s.Suite.Fail(err.Error())
	}
}

func (s *AuthzRSTestSuite) TestUnary_AuthRS_VerifyError() {
	_, err := s.Client.Ping(ctxWithToken(s.SimpleCtx(), "bearer", "bad_token"), goodPing)
	assert.Error(s.T(), err, "there must be an error")
	assert.Equal(s.T(), codes.Unauthenticated, status.Code(err), "token format is not valid")
}
