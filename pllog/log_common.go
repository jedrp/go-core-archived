package pllog

import "context"

const (
	RequestIDHeaderKey     = "Request-Id"
	CorrelationIDHeaderKey = "Correlation-Id"
	RequestID              = "RequestId"
	CorrelationID          = "CorrelationId"
)

func CreateLogEntryFromContext(ctx context.Context, log PlLogger) PlLogentry {
	return log.WithFields(map[string]interface{}{
		CorrelationID: ctx.Value(CorrelationID).(string),
		RequestID:     ctx.Value(RequestID).(string),
	})
}
