package middlewares

import (
	"template-go-api/internal/pkg/logger"

	"github.com/labstack/echo/v4"
)

func traceMiddeware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		headers := c.Request().Header
		cc, _ := c.(*customContext)

		traceData := logger.TraceData{
			TransactionID: headers.Get("x-transaction-id"),
			SessionID:     headers.Get("x-session-id"),
			ChannelID:     headers.Get("x-channel-id"),
			ConsumerName:  headers.Get("x-consumer"),
		}

		cc.setTraceData(traceData)

		return next(c)
	}
}
