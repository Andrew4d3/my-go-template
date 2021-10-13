package routes

import (
	"template-go-api/internal/app/controllers"
	"template-go-api/internal/app/middlewares"
	"template-go-api/internal/app/models/users"
	"template-go-api/internal/app/repository"
	"template-go-api/internal/pkg/drivers/mongodb"

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
	// Example: Delete this once it's clear how it works
	e.GET("/user/:name", func(c echo.Context) error {
		userCol, _ := mongodb.GetCollection("users")
		users := users.New(repository.New(userCol))

		user, err := users.FindByName(c.Param("name"))
		if err != nil {
			return echo.NewHTTPError(500, err.Error())
		}

		if user.ID == "" {
			return echo.NewHTTPError(404, "User not found")
		}

		return c.JSON(200, user)
	})
}
