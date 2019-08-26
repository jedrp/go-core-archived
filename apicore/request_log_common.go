package apicore

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/HoaHuynhSoft/go-core/pllog"
	uuid "github.com/satori/go.uuid"
)

func HandlePanicMiddleware(handler http.Handler, logger pllog.PlLogger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := NewLogContext(r)
		defer func() {
			if rErr := recover(); rErr != nil {
				ctx := r.Context()
				fmt.Println("rid", ctx.Value("RequestId"))
				if logger != nil {
					pllog.CreateLogEntryFromContext(r.Context(), logger).Error(rErr)
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
	fmt.Println(reqID)
	if reqID != "" {
		ctx = context.WithValue(ctx, "RequestId", reqID)
	} else {
		ctx = context.WithValue(ctx, "RequestId", uuid.NewV4().String())
	}

	corID := r.Header.Get(pllog.CorrelationIDHeaderKey)
	if corID != "" {
		ctx = context.WithValue(ctx, pllog.CorrelationID, corID)
	}
	return ctx
}
