package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Middleware interface {
	Use(middleware ...echo.MiddlewareFunc)
}

func SetMiddlewares(e Middleware) {
	e.Use(middleware.Secure())
	e.Use(bindCustomContext)
	e.Use(traceMiddeware)
	e.Use(errorMiddleware)
	e.Use(logMiddleware)
}
