package middlewares

import (
	"errors"
	"template-go-api/internal/pkg/logger"
	"template-go-api/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newCC() *customContext {
	c, _ := getTestContext()
	return &customContext{c}
}

func Test_logMiddleware(t *testing.T) {
	ogNewLogger := newLogger

	defer func() {
		newLogger = ogNewLogger
	}()

	t.Run("Should return an error if there is a problem getting a new logger instance", func(t *testing.T) {
		newLogger = func(_ logger.TraceData) (logger.Logger, error) {
			return nil, errors.New("Boom logger")
		}

		err := logMiddleware(stubController)(newCC())

		assert.Errorf(t, err, "Boom logger")
	})

	t.Run("Should set the logger to the custom context", func(t *testing.T) {
		mockedLogger := new(mocks.Logger)
		newLogger = func(_ logger.TraceData) (logger.Logger, error) {
			return mockedLogger, nil
		}

		cc := newCC()
		err := logMiddleware(stubController)(cc)

		assert.NoError(t, err)
		assert.Equal(t, mockedLogger, cc.CustomLogger())
	})
}
