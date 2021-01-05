package authz

import (
	"context"
	"fmt"
	"testing"

	jwt "github.com/cristalhq/jwt/v3"
	jwt_helper "github.com/dgrijalva/jwt-go"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	pb_testproto "github.com/grpc-ecosystem/go-grpc-middleware/testing/testproto"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
)

var (
	authedMarker = "some_context_marker"
	goodPing     = &pb_testproto.PingRequest{Value: "something", SleepTimeMs: 9999}
	emptyPing    = &pb_testproto.Empty{}
)

type userClaims struct {
	jwt.StandardClaims
	Scopes []string `json:"scopes"`
}

type assertingPingService struct {
	pb_testproto.TestServiceServer
	T *testing.T
}

func (s *assertingPingService) PingError(ctx context.Context, ping *pb_testproto.PingRequest) (*pb_testproto.Empty, error) {
	assertAuthMarkerExists(ctx, s.T)
	return s.TestServiceServer.PingError(ctx, ping)
}

func (s *assertingPingService) PingList(ping *pb_testproto.PingRequest, stream pb_testproto.TestService_PingListServer) error {
	assertAuthMarkerExists(stream.Context(), s.T)
	return s.TestServiceServer.PingList(ping, stream)
}

func assertAuthMarkerExists(ctx context.Context, t *testing.T) {
	assert.Equal(t, "marker_exists", ctx.Value(authedMarker).(string), "auth marker from buildDummyAuthFunction must be passed around")
}

func ctxWithToken(ctx context.Context, scheme string, token string) context.Context {
	md := metadata.Pairs("authorization", fmt.Sprintf("%s %s", scheme, token))
	return metautils.NiceMD(md).ToOutgoing(ctx)
}

func ctxWithTokenIncoming(ctx context.Context, scheme string, token string) context.Context {
	md := metadata.Pairs("authorization", fmt.Sprintf("%s %s", scheme, token))
	return metautils.NiceMD(md).ToIncoming(ctx)
}

func TestNewConfig(t *testing.T) {
	rsaPubKey, err := jwt_helper.ParseRSAPublicKeyFromPEM([]byte(pubKey))
	if err != nil {
		t.Fatal(err)
	}

	expectedConfig := &Config{
		Algorithm:    jwt.HS256,
		ScopeMapping: map[string][]string{"one": {"one"}},
		ValidateScopeFunc: func(tokenClaims []byte, scopes []string) error {
			return nil
		},
		Audience:                  "admin",
		DisableAudienceValidation: true,
		HSKey:                     hsKey,
		RSAPublicKey:              rsaPubKey,
	}

	config := NewConfig(
		jwt.HS256,
		WithScopeMapping(map[string][]string{"one": {"one"}}),
		WithValidateScopeFunc(func(tokenClaims []byte, scopes []string) error { return nil }),
		WithAudience("admin"),
		WithDisableAudienceValidation(true),
		WithHSKey(hsKey),
		WithRSAPubKey(rsaPubKey),
	)

	if diff := cmp.Diff(expectedConfig, config, cmpopts.IgnoreFields(Config{}, "ValidateScopeFunc"), cmpopts.IgnoreUnexported(Config{})); diff != "" {
		t.Fatalf(diff)
	}

	tk, err := createTokenHS([]string{}, jwt.HS256, hsKey, "admin")
	if err != nil {
		t.Fatal(err)
	}
	ctx := ctxWithTokenIncoming(context.Background(), "bearer", tk.String())

	_, err = config.AuthFunc(ctx)
	if err != nil {
		t.Fatal(err)
	}
}
