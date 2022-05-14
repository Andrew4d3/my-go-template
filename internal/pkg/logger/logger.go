package logger

import (
	"encoding/json"
	"fmt"
	"net/url"
	"template-go-api/configs"
	"template-go-api/internal/pkg/customerror"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type hostData struct {
	ContainerID string `json:"containerId"`
	Name        string `json:"name"`
	IP          string `json:"ip"`
}

type loggerData struct {
	RequestID     string   `json:"requestId"`
	TransactionID string   `json:"transactionId"`
	SessionID     string   `json:"sessionId"`
	ChannelID     string   `json:"channelId"`
	ConsumerName  string   `json:"consumerName"`
	Environment   string   `json:"environment"`
	Host          hostData `json:"host"`
}

// TraceData defines the input trace data
type TraceData struct {
	TransactionID string
	SessionID     string
	ChannelID     string
	ConsumerName  string
}

// Logger defines the logger common interface
type Logger interface {
	Info(msg string, payload ...map[string]interface{})
	Debug(msg string, payload ...map[string]interface{})
	Warn(msg string, payload ...map[string]interface{})
	Error(err error)
}

type zapI interface {
	Info(msg string, fields ...zap.Field)
	Debug(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Sync() error
}

// ZapLogger defines a zap logger common instance
type ZapLogger struct {
	data       loggerData
	logHandler zapI
}

var getLoggerConfig = configs.GetLoggerConfig

// Info logs an info message
func (z *ZapLogger) Info(msg string, payload ...map[string]interface{}) {
	z.logHandler.Info(msg, zap.Any("payload", checkPayload(payload)))
	z.logHandler.Sync()
}

// Debug logs a debug message
func (z *ZapLogger) Debug(msg string, payload ...map[string]interface{}) {
	z.logHandler.Debug(msg, zap.Any("payload", checkPayload(payload)))
	z.logHandler.Sync()
}

// Warn logs a warn message
func (z *ZapLogger) Warn(msg string, payload ...map[string]interface{}) {
	z.logHandler.Warn(msg, zap.Any("payload", checkPayload(payload)))
	z.logHandler.Sync()
}

// Error logs an error message
func (z *ZapLogger) Error(err error) {
	logPayload := zap.Any("payload", map[string]interface{}{})
	errMsg := err.Error()

	if customError, isCustom := (err).(customerror.CustomError); isCustom {
		logPayload = zap.Any("payload", customError.GetData())
	}

	if echoError, isEcho := (err).(*echo.HTTPError); isEcho {
		logPayload = zap.Any("payload", map[string]interface{}{
			"statusCode": echoError.Code,
			"data":       echoError.Internal,
		})
	}

	z.logHandler.Error(errMsg, logPayload)
	z.logHandler.Sync()
}

type lumberjackSink struct {
	*lumberjack.Logger
}

func (lumberjackSink) Sync() error {
	return nil
}

func checkPayload(payload []map[string]interface{}) map[string]interface{} {
	if len(payload) == 1 {
		return payload[0]
	}

	return map[string]interface{}{}
}

func getZapLevel(level string) zapcore.Level {
	if level == "debug" {
		return zapcore.DebugLevel
	} else if level == "warn" {
		return zapcore.WarnLevel
	} else if level == "error" {
		return zapcore.ErrorLevel
	} else {
		return zapcore.InfoLevel
	}
}

func newZapLogger(ld loggerData, loggerConfig configs.LoggerConfig) (*zap.Logger, error) {
	var initialFields map[string]interface{}
	jsonInitialFields, _ := json.Marshal(&ld)

	json.Unmarshal(jsonInitialFields, &initialFields)

	zapConfig := zap.Config{
		Level:       zap.NewAtomicLevelAt(getZapLevel(loggerConfig.Level)),
		Encoding:    "json",
		OutputPaths: []string{"stdout", fmt.Sprint("lumberjack:", loggerConfig.Filepath)},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:  "message",
			LevelKey:    "level",
			TimeKey:     "date",
			EncodeLevel: zapcore.LowercaseLevelEncoder,
			EncodeTime:  zapcore.ISO8601TimeEncoder,
			LineEnding:  zapcore.DefaultLineEnding,
		},
		InitialFields: initialFields,
	}

	return zapConfig.Build()
}

var zapRegisterSink = zap.RegisterSink

// SetupFileStore sets the directory where logs will be stored
func SetupFileStore() error {
	logConfig := getLoggerConfig()

	lumbLogger := lumberjack.Logger{
		Filename:   logConfig.Filepath,
		MaxSize:    1024,
		MaxBackups: 30,
		MaxAge:     5,
		Compress:   true,
	}

	return zapRegisterSink("lumberjack", func(*url.URL) (zap.Sink, error) {
		return lumberjackSink{
			Logger: &lumbLogger,
		}, nil
	})
}

var _newZapLogger = newZapLogger

// NewLogger creates a new logger instance
func NewLogger(traceData TraceData) (Logger, error) {
	loggerConfig := getLoggerConfig()
	ld := loggerData{
		RequestID:     uuid.NewString(),
		TransactionID: traceData.TransactionID,
		SessionID:     traceData.SessionID,
		ChannelID:     traceData.ChannelID,
		ConsumerName:  traceData.ConsumerName,
		Environment:   loggerConfig.Env,
		Host: hostData{
			ContainerID: loggerConfig.ContainerID,
			Name:        loggerConfig.HostName,
			IP:          loggerConfig.HostIP,
		},
	}

	zapLogger, err := _newZapLogger(ld, loggerConfig)

	return &ZapLogger{
		data:       ld,
		logHandler: zapLogger,
	}, err
}
