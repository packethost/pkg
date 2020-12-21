package authz

import (
	"context"
	"fmt"
	"testing"

	jwt "github.com/cristalhq/jwt/v3"
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
	md := metadata.Pairs("authorization", fmt.Sprintf("%s %v", scheme, token))
	nCtx := metautils.NiceMD(md).ToOutgoing(ctx)
	return nCtx
}
