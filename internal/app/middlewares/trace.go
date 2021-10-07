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
			TransactionId: headers.Get("x-transaction-id"),
			SessionId:     headers.Get("x-session-id"),
			ChannelId:     headers.Get("x-channel-id"),
			ConsumerName:  headers.Get("x-consumer"),
		}

		cc.setTraceData(traceData)

		return next(c)
	}
}
