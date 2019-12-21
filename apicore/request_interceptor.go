package apicore

import (
	"context"
	"fmt"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/jedrp/go-core/pllog"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

//UnaryServerRequestContextInterceptor set correlation id to context
func UnaryServerRequestContextInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {

		ctx, err = setUpRequestInfoToContext(ctx)
		if err != nil {
			return
		}
		return handler(ctx, req)
	}
}

func setUpRequestInfoToContext(baseCtx context.Context) (context.Context, error) {
	ctx := context.WithValue(baseCtx, pllog.RequestID, uuid.NewV4().String())

	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		corID := md.Get(pllog.CorrelationIDHeaderKey)
		if len(corID) > 0 && corID[0] != "" {
			ctx = context.WithValue(ctx, pllog.CorrelationID, corID[0])
		}
		return ctx, nil
	}
	return nil, fmt.Errorf("Unable to obtain metadata")
}

// StreamServerRequestInterceptor adding request and correlation id to context
func StreamServerRequestInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		newCtx, err := setUpRequestInfoToContext(stream.Context())
		if err != nil {
			return err
		}
		wrapped := grpc_middleware.WrapServerStream(stream)
		wrapped.WrappedContext = newCtx
		return handler(srv, wrapped)
	}
}
