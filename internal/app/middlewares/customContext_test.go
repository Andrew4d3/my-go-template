package middlewares

import (
	"template-go-api/internal/pkg/logger"
	"template-go-api/mocks"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func Test_CustomContext(t *testing.T) {
	mockedLogger := new(mocks.Logger)
	mockedTraceData := new(logger.TraceData)
	c, _ := getTestContext()
	cc := &customContext{c}

	cc.setCtxLogger(mockedLogger)
	cc.setTraceData(*mockedTraceData)

	t.Run("Should get the corresponding custom logger", func(t *testing.T) {
		customLogger := cc.CustomLogger()
		assert.Equal(t, mockedLogger, customLogger)
	})

	t.Run("Should get the corresponding trace data object", func(t *testing.T) {
		traceData := cc.TraceData()
		assert.Equal(t, *mockedTraceData, traceData)
	})
}

func Test_bindCustomContext(t *testing.T) {
	c, _ := getTestContext()

	t.Run("Should bind the corresponding custom context", func(t *testing.T) {
		mockedController := func(c echo.Context) error {
			_, isCC := c.(*customContext)
			assert.True(t, isCC)
			return c.JSON(200, "OK")
		}

		err := bindCustomContext(mockedController)(c)
		assert.NoError(t, err)
	})

}
