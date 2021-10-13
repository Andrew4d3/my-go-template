package middlewares

import (
	"strings"
	"template-go-api/configs"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

var getJWTSecret = configs.GetAppConfig

// AuthorizationMiddleware is a middleware that verifies if JWT is valid. Otherwise it returns error
func AuthorizationMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		jwtSecret := getJWTSecret().JWTSecret
		authHeader := c.Request().Header.Get("authorization")

		if authHeader == "" {
			return echo.NewHTTPError(401, "Authorization header is not present")
		}

		tokenString := ""
		if headerParts := strings.Split(authHeader, " "); len(headerParts) >= 2 {
			tokenString = headerParts[1]
		}

		_, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})

		if err == nil {
			return next(c)
		}

		return echo.NewHTTPError(401, err.Error())
	}
}
