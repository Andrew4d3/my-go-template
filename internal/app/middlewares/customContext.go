package middlewares

import (
	"template-go-api/internal/pkg/logger"

	"github.com/labstack/echo/v4"
)

type customContext struct {
	echo.Context
}

type CustomContext interface {
	CustomLogger() logger.Logger
	TraceData() logger.TraceData
}

func (c *customContext) setCtxLogger(logger logger.Logger) {
	c.Set("logger", logger)
}

func (c *customContext) setTraceData(traceData logger.TraceData) {
	c.Set("traceData", traceData)
}

func (c customContext) CustomLogger() logger.Logger {
	logger, _ := c.Get("logger").(logger.Logger)
	return logger
}

func (c customContext) TraceData() logger.TraceData {
	traceData, _ := c.Get("traceData").(logger.TraceData)
	return traceData
}

func bindCustomContext(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cc := &customContext{c}
		return next(cc)
	}
}
