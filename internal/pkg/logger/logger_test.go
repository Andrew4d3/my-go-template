package logger

import (
	"net/url"
	"template-go-api/configs"
	"template-go-api/mocks"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func getTestLogger() (*ZapLogger, *mocks.ZapMock) {
	zapMocked := new(mocks.ZapMock)
	testLogger := &ZapLogger{
		logHandler: zapMocked,
	}

	return testLogger, zapMocked
}

func setExpectations(zapPayload zapcore.Field, zapMocked *mocks.ZapMock) {
	zapMocked.On("Info", "test info", zapPayload)
	zapMocked.On("Debug", "test debug", zapPayload)
	zapMocked.On("Warn", "test warn", zapPayload)
	zapMocked.On("Sync").Return(nil)
}

func TestLoggerMethods(t *testing.T) {

	t.Run("Should log each level without payload", func(t *testing.T) {
		testLogger, zapMocked := getTestLogger()
		zapPayload := zap.Any("payload", map[string]interface{}{})

		setExpectations(zapPayload, zapMocked)

		testLogger.Info("test info")
		testLogger.Debug("test debug")
		testLogger.Warn("test warn")

		zapMocked.AssertExpectations(t)
	})

	t.Run("Should log each level with payload", func(t *testing.T) {
		testLogger, zapMocked := getTestLogger()
		zapPayload := zap.Any("payload", map[string]interface{}{"msg": "test"})

		setExpectations(zapPayload, zapMocked)

		payload := map[string]interface{}{"msg": "test"}
		testLogger.Info("test info", payload)
		testLogger.Debug("test debug", payload)
		testLogger.Warn("test warn", payload)

		zapMocked.AssertExpectations(t)
	})

	t.Run("Should log custom errors correctly", func(t *testing.T) {
		testLogger, zapMocked := getTestLogger()
		mockedCustomErr := new(mocks.CustomError)
		customErrorData := map[string]interface{}{
			"statusCode": 400,
		}
		zapPayload := zap.Any("payload", customErrorData)

		mockedCustomErr.On("Error").Return("test error")
		mockedCustomErr.On("GetData").Return(customErrorData)
		zapMocked.On("Error", "test error", zapPayload)
		zapMocked.On("Sync").Return(nil)

		testLogger.Error(mockedCustomErr)
	})

	t.Run("Should log echo errors correctly", func(t *testing.T) {
		testLogger, zapMocked := getTestLogger()
		echoError := echo.NewHTTPError(404, "echo error")
		echoErrorData := map[string]interface{}{
			"statusCode": echoError.Code,
			"data":       echoError.Internal,
		}
		zapPayload := zap.Any("payload", echoErrorData)

		zapMocked.On("Error", "code=404, message=echo error", zapPayload)
		zapMocked.On("Sync").Return(nil)

		testLogger.Error(echoError)
	})
}

func Test_NewLogger(t *testing.T) {
	ogGetLoggerConfig := getLoggerConfig
	ogNewZapLogger := _newZapLogger

	defer func() {
		getLoggerConfig = ogGetLoggerConfig
		_newZapLogger = ogNewZapLogger
	}()

	t.Run("Should return a ZapLogger instance", func(t *testing.T) {
		getLoggerConfig = func() configs.LoggerConfig {
			return configs.LoggerConfig{}
		}

		_newZapLogger = func(ld loggerData, loggerConfig configs.LoggerConfig) (*zap.Logger, error) {
			return &zap.Logger{}, nil
		}

		logger, err := NewLogger(TraceData{})
		assert.NoError(t, err)
		assert.NotNil(t, logger)
	})
}

func Test_getZapLevel(t *testing.T) {
	t.Run("Should return the correct zap level", func(t *testing.T) {
		assert.Equal(t, zapcore.InfoLevel, getZapLevel("info"))
		assert.Equal(t, zapcore.DebugLevel, getZapLevel("debug"))
		assert.Equal(t, zapcore.WarnLevel, getZapLevel("warn"))
		assert.Equal(t, zapcore.ErrorLevel, getZapLevel("error"))
	})
}

func Test_SetupFileStore(t *testing.T) {
	ogGetLoggerConfig := getLoggerConfig
	ogZapRegisterSink := zapRegisterSink

	defer func() {
		getLoggerConfig = ogGetLoggerConfig
		zapRegisterSink = ogZapRegisterSink
	}()

	t.Run("Should create the file store for the logs", func(t *testing.T) {
		zapRegisterSink = func(_ string, fn func(*url.URL) (zap.Sink, error)) error {
			ls, err := fn(nil)
			assert.NotNil(t, ls)
			assert.NoError(t, err)
			return nil
		}

		getLoggerConfig = func() configs.LoggerConfig {
			return configs.LoggerConfig{}
		}

		err := SetupFileStore()
		assert.NoError(t, err)
	})
}
