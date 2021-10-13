package middlewares

import (
	"template-go-api/internal/pkg/logger"

	"github.com/labstack/echo/v4"
)

var newLogger = logger.NewLogger

func logMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cc := c.(*customContext)
		mylogger, err := newLogger(cc.TraceData())

		if err != nil {
			return err
		}

		cc.setCtxLogger(mylogger)
		return next(cc)
	}
}
