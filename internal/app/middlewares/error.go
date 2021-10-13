package middlewares

import (
	"fmt"
	"template-go-api/internal/pkg/customerror"

	"github.com/labstack/echo/v4"
)

func transformError(err error) *echo.HTTPError {
	echoError, isEcho := err.(*echo.HTTPError)
	if isEcho {
		return echoError
	}

	customError, isCustom := err.(customerror.CustomError)
	if isCustom {
		return echo.NewHTTPError(customError.GetStatusCode(), customError.Error())
	}

	return echo.NewHTTPError(500, err.Error())
}

func logError(c echo.Context, err error) {
	cc, isCC := c.(CustomContext)
	if isCC && cc.CustomLogger() != nil {
		cc.CustomLogger().Error(err)
	} else {
		fmt.Println(err.Error())
	}
}

var _logError = logError

func errorMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)
		if err != nil {
			_logError(c, err)
			err = transformError(err)
		}

		return err
	}
}
