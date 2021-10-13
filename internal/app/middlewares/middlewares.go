package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Middleware interface for DI
type Middleware interface {
	Use(middleware ...echo.MiddlewareFunc)
}

// SetMiddlewares registers all the required middlewares
func SetMiddlewares(e Middleware) {
	e.Use(middleware.Secure())
	e.Use(bindCustomContext)
	e.Use(traceMiddeware)
	e.Use(errorMiddleware)
	e.Use(logMiddleware)
}
