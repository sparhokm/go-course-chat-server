package interceptor

import (
	"context"
	"errors"

	descAccess "github.com/sparhokm/go-course-ms-auth/pkg/access_v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type accessStreamCheck struct {
	client descAccess.AccessV1Client
}

func NewAccessStreamCheck(client descAccess.AccessV1Client) *accessStreamCheck {
	return &accessStreamCheck{client}
}

type key int

var userKey key = 1

func (a *accessStreamCheck) StreamServerInterceptor(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	md, ok := metadata.FromIncomingContext(ss.Context())
	if !ok {
		return errors.New("metadata is not provided")
	}

	authHeader, ok := md["authorization"]
	if !ok || len(authHeader) == 0 {
		return errors.New("authorization header is not provided")
	}

	ctx := metadata.NewOutgoingContext(context.Background(), md)
	k, err := a.client.Check(ctx, &descAccess.CheckIn{
		EndpointAddress: info.FullMethod,
	})
	if err != nil {
		return errors.New("access denied")
	}

	return handler(srv, &serverStream{
		ss,
		context.WithValue(ss.Context(), userKey, k.GetId()),
	})
}

func GetUserId(ctx context.Context) (int64, bool) {
	u, ok := ctx.Value(userKey).(int64)
	return u, ok
}

type serverStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (s serverStream) Context() context.Context {
	return s.ctx
}
