package pllog

import "context"

const (
	RequestIDHeaderKey     = "Request-Id"
	CorrelationIDHeaderKey = "Correlation-Id"
	RequestID              = "RequestId"
	CorrelationID          = "CorrelationId"
)

func CreateLogEntryFromContext(ctx context.Context, log PlLogger) PlLogentry {
	correlationID, ok := ctx.Value(CorrelationID).(string)
	if ok {
		return log.WithFields(map[string]interface{}{
			CorrelationID: correlationID,
			RequestID:     ctx.Value(RequestID).(string),
		})
	}
	return log.WithFields(map[string]interface{}{
		RequestID: ctx.Value(RequestID).(string),
	})
}
