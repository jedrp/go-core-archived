package apicore

import (
	"context"
	"runtime/debug"

	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/jedrp/go-core/pllog"
	"google.golang.org/grpc"
)

//UnaryServerPanicInterceptor hanle unexpected panic
func UnaryServerPanicInterceptor(logger pllog.PlLogger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				if logger != nil {
					pllog.CreateLogEntryFromContext(ctx, logger).Error(r, string(debug.Stack()))
				}
			}
		}()

		return handler(ctx, req)
	}
}

func getRecoveryHandlerFuncContextHandler(logger pllog.PlLogger) grpc_recovery.RecoveryHandlerFuncContext {
	return grpc_recovery.RecoveryHandlerFuncContext(
		func(ctx context.Context, p interface{}) error {
			if logger != nil {
				pllog.CreateLogEntryFromContext(ctx, logger).Error(p, string(debug.Stack()))
			}
			return nil
		},
	)
}
