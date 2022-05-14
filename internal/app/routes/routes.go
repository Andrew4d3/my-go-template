package routes

import (
	"template-go-api/internal/app/controllers"
	"template-go-api/internal/app/middlewares"

	"github.com/labstack/echo/v4"
)

// WebServer defines an interface for defining HTTP routes (endpoints)
type WebServer interface {
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// SetRoutes register all the required routes
func SetRoutes(e WebServer) {
	e.GET("/health", controllers.HealthCheck)
	// Example: Delete this once it's clear how it works
	e.GET("/protected", controllers.HealthCheck, middlewares.AuthorizationMiddleware)
}
