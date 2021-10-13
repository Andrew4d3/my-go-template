package middlewares

import (
	"errors"
	"template-go-api/mocks"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func Test_transformError(t *testing.T) {
	t.Run("Should return the unaltered error if it's of echo type", func(t *testing.T) {
		echoError := echo.NewHTTPError(400, "boom!")
		resultError := transformError(echoError)

		assert.Equal(t, echoError, resultError)
	})

	t.Run("Should return the corresponding echo error if error is custom", func(t *testing.T) {
		mockedCustomError := new(mocks.CustomError)
		mockedCustomError.On("GetStatusCode").Return(403)
		mockedCustomError.On("Error").Return("Forbidden access")

		resultError := transformError(mockedCustomError)
		assert.Equal(t, resultError.Code, 403)
		assert.Equal(t, resultError.Message, "Forbidden access")
	})

	t.Run("Should return a default 500 echo error if error is unknown", func(t *testing.T) {
		normalError := errors.New("Normal error")
		resultError := transformError(normalError)
		assert.Equal(t, resultError.Code, 500)
		assert.Equal(t, resultError.Message, "Normal error")
	})
}

func Test_logError(t *testing.T) {
	mockedLogger := new(mocks.Logger)
	testError := errors.New("Boom error")
	mockedLogger.On("Error", testError)

	t.Run("Should log the error using the custom logger if exists", func(t *testing.T) {
		c, _ := getTestContext()
		cc := &customContext{c}
		cc.setCtxLogger(mockedLogger)

		logError(cc, testError)
		mockedLogger.AssertExpectations(t)
	})

	t.Run("Should log the error using the custom logger if exists", func(t *testing.T) {
		c, _ := getTestContext()
		cc := &customContext{c}
		mockedCustomErr := new(mocks.CustomError)
		mockedCustomErr.On("Error").Return("Boom")

		logError(cc, mockedCustomErr)
		mockedCustomErr.AssertExpectations(t)
	})
}

func Test_errorMiddleware(t *testing.T) {
	ogLogError := _logError
	defer func() {
		_logError = ogLogError
	}()

	t.Run("Should log and transform an error if controller returns an error", func(t *testing.T) {
		testError := errors.New("Boom error")
		mockedController := func(c echo.Context) error {
			return testError
		}

		_logError = func(_ echo.Context, err error) {
			assert.Equal(t, testError, err)
		}

		c, _ := getTestContext()

		err := errorMiddleware(mockedController)(c)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "500")
	})
}
