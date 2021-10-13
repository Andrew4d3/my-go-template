package middlewares

import (
	"reflect"
	"runtime"
	"template-go-api/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var foo interface{}

func bar() {}

func Test_SetMiddlewares(t *testing.T) {
	mockedMiddlewareServer := new(mocks.Middleware)
	mockedMiddlewareServer.On("Use", mock.Anything)
	SetMiddlewares(mockedMiddlewareServer)

	assertFn := func(t *testing.T, expected interface{}, actual interface{}) {
		funcName1 := runtime.FuncForPC(reflect.ValueOf(expected).Pointer()).Name()
		funcName2 := runtime.FuncForPC(reflect.ValueOf(actual).Pointer()).Name()

		assert.Equal(t, funcName1, funcName2)
	}

	t.Run("Should register all the required middlewares in correct order", func(t *testing.T) {
		calls := mockedMiddlewareServer.Calls
		assertFn(t, bindCustomContext, calls[1].Arguments[0])
		assertFn(t, traceMiddeware, calls[2].Arguments[0])
		assertFn(t, errorMiddleware, calls[3].Arguments[0])
		assertFn(t, logMiddleware, calls[4].Arguments[0])
	})

}
