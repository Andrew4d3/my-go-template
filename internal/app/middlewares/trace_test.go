package middlewares

import (
	"template-go-api/internal/pkg/logger"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testHeaders = map[string]string{
	"x-transaction-id": "transaction-123",
	"x-session-id":     "session-123",
	"x-channel-id":     "test",
	"x-consumer":       "test-consumer",
}

func Test_traceMiddeware(t *testing.T) {
	cc := newCC()

	for key, value := range testHeaders {
		cc.Request().Header.Set(key, value)
	}

	t.Run("Should set the corresponding trace data", func(t *testing.T) {
		err := traceMiddeware(stubController)(cc)
		assert.NoError(t, err)
		traceData := cc.TraceData()
		assert.Equal(t, logger.TraceData{
			TransactionID: "transaction-123",
			SessionID:     "session-123",
			ChannelID:     "test",
			ConsumerName:  "test-consumer",
		}, traceData)
	})
}
