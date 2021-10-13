package routes

import (
	"reflect"
	"runtime"
	"template-go-api/internal/app/controllers"
	"template-go-api/mocks"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_SetRoutes(t *testing.T) {
	mockedWebServer := new(mocks.WebServer)
	mockedWebServer.On("GET", mock.Anything, mock.Anything).Return(&echo.Route{})
	// Delete the following line once the /protected sample endpoint is deleted
	mockedWebServer.On("GET", mock.Anything, mock.Anything, mock.Anything).Return(&echo.Route{})
	SetRoutes(mockedWebServer)

	assertFn := func(t *testing.T, expected interface{}, actual interface{}) {
		funcName1 := runtime.FuncForPC(reflect.ValueOf(expected).Pointer()).Name()
		funcName2 := runtime.FuncForPC(reflect.ValueOf(actual).Pointer()).Name()

		assert.Equal(t, funcName1, funcName2)
	}

	t.Run("Should register all the required controllers and their corresponding path", func(t *testing.T) {
		calls := mockedWebServer.Calls

		assert.Equal(t, "/health", calls[0].Arguments[0].(string))
		assertFn(t, controllers.HealthCheck, calls[0].Arguments[1])
	})
}
