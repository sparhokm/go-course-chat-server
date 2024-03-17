package interceptor

import (
	"context"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	descAccess "github.com/sparhokm/go-course-ms-auth/pkg/access_v1"
)

type accessCheck struct {
	client descAccess.AccessV1Client
}

func NewAccessCheck(client descAccess.AccessV1Client) *accessCheck {
	return &accessCheck{client}
}

func (a *accessCheck) UnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("metadata is not provided")
	}

	authHeader, ok := md["authorization"]
	if !ok || len(authHeader) == 0 {
		return nil, errors.New("authorization header is not provided")
	}

	ctx = metadata.NewOutgoingContext(context.Background(), md)
	k, err := a.client.Check(ctx, &descAccess.CheckIn{
		EndpointAddress: info.FullMethod,
	})
	if err != nil {
		return nil, errors.New("access denied")
	}
	_ = k
	return handler(ctx, req)
}
