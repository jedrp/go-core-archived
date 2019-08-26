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
		NewLogContext(r)
		defer func() {
			if rErr := recover(); rErr != nil {
				if logger != nil {
					pllog.CreateLogEntryFromContext(r.Context(), logger).Error(rErr)
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(500)
				response, _ := json.Marshal(map[string]string{"message": "Internal server error"})
				w.Write(response)
			}
		}()
		handler.ServeHTTP(w, r)
	})
}

func NewLogContext(r *http.Request) {
	ctx := r.Context()
	reqID := r.Header.Get(pllog.RequestIDHeaderKey)
	fmt.Println(reqID)
	if reqID != "" {
		context.WithValue(ctx, pllog.RequestID, reqID)
	} else {
		context.WithValue(ctx, pllog.RequestID, uuid.NewV4().String())
	}

	corID := r.Header.Get(pllog.CorrelationIDHeaderKey)
	if corID != "" {
		context.WithValue(ctx, pllog.CorrelationID, corID)
	}
}
