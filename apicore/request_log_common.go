package apicore

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/HoaHuynhSoft/go-core/pllog"
	uuid "github.com/satori/go.uuid"
)

func HandlePanicMiddleware(handler http.Handler, logger pllog.PlLogger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := NewLogContext(r)
		defer func() {
			if r := recover(); r != nil {
				if logger != nil {
					pllog.CreateLogEntryFromContext(ctx, logger).Error(r)
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

func NewLogContext(r *http.Request) context.Context {
	ctx := r.Context()
	reqID := r.Header.Get(pllog.RequestIDHeaderKey)
	if reqID != "" {
		context.WithValue(ctx, pllog.RequestID, reqID)
	} else {
		context.WithValue(ctx, pllog.RequestID, uuid.NewV4().String())
	}

	corID := r.Header.Get(pllog.CorrelationIDHeaderKey)
	if corID != "" {
		context.WithValue(ctx, pllog.CorrelationID, corID)
	}
	return ctx
}
