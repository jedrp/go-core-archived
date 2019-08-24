package apicore

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/HoaHuynhSoft/go-core/pllog"
	uuid "github.com/satori/go.uuid"
)

const (
	RequestIDHeaderKey     = "Request-Id"
	CorrelationIDHeaderKey = "Correlation-Id"
	RequestID              = "RequestId"
	CorrelationID          = "CorrelationId"
)

func HandleExceptionMiddleware(handler http.Handler, logger pllog.PlLogger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := NewLogContext(r)
		defer func() {
			if r := recover(); r != nil {
				if logger != nil {
					CreateLogEntryFromContext(ctx, logger).Error(r)
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(500)
				response, _ := json.Marshal(map[string]string{"message": "Internal server error"})
				w.Write(response)
			}
		}()
		handler.ServeHTTP(w, r.WithContext(ctx))
	})
}

func CreateLogEntryFromContext(ctx context.Context, log pllog.PlLogger) pllog.PlLogentry {
	return log.WithFields(map[string]interface{}{
		CorrelationID: ctx.Value(CorrelationID),
		RequestID:     ctx.Value(RequestID),
	})
}

func NewLogContext(r *http.Request) context.Context {
	ctx := r.Context()
	reqID := r.Header.Get(RequestIDHeaderKey)
	if reqID != "" {
		context.WithValue(ctx, RequestID, reqID)
	} else {
		context.WithValue(ctx, RequestID, uuid.NewV4().String())
	}

	corID := r.Header.Get(CorrelationIDHeaderKey)
	if corID != "" {
		context.WithValue(ctx, CorrelationID, corID)
	}
	return ctx
}
