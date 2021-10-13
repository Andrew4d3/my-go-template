package controllers

import (
	"template-go-api/configs"

	"github.com/labstack/echo/v4"
)

var getAppConfig = configs.GetAppConfig

// HealthCheck GET /health
func HealthCheck(c echo.Context) error {
	return c.JSON(200, map[string]string{
		"message":     "API is healthy",
		"environment": getAppConfig().Environment,
	})
}
