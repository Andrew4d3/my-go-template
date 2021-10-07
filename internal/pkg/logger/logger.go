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
	ContainerId string `json:"containerId"`
	Name        string `json:"name"`
	Ip          string `json:"ip"`
}

type loggerData struct {
	RequestId     string   `json:"requestId"`
	TransactionId string   `json:"transactionId"`
	SessionId     string   `json:"sessionId"`
	ProductId     string   `json:"productId"`
	ChannelId     string   `json:"channelId"`
	ConsumerName  string   `json:"consumerName"`
	Environment   string   `json:"environment"`
	Commerce      string   `json:"commerce"`
	Host          hostData `json:"host"`
}

type TraceData struct {
	TransactionId string
	SessionId     string
	ChannelId     string
	ConsumerName  string
}

type Logger interface {
	Info(msg string, payload ...map[string]interface{})
	Debug(msg string, payload ...map[string]interface{})
	Warn(msg string, payload ...map[string]interface{})
	Error(err error)
}

type ZapLogger struct {
	data       loggerData
	logHandler *zap.Logger
}

func (z *ZapLogger) Info(msg string, payload ...map[string]interface{}) {

	z.logHandler.Info(msg, zap.Any("payload", checkPayload(payload)))
	z.logHandler.Sync()
}

func (z *ZapLogger) Debug(msg string, payload ...map[string]interface{}) {

	z.logHandler.Debug(msg, zap.Any("payload", checkPayload(payload)))
	z.logHandler.Sync()
}

func (z *ZapLogger) Warn(msg string, payload ...map[string]interface{}) {

	z.logHandler.Warn(msg, zap.Any("payload", checkPayload(payload)))
	z.logHandler.Sync()
}

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

func SetupFileStore() error {
	logConfig := configs.GetLoggerConfig()

	lumbLogger := lumberjack.Logger{
		Filename:   logConfig.Filepath,
		MaxSize:    1024,
		MaxBackups: 30,
		MaxAge:     5,
		Compress:   true,
	}

	return zap.RegisterSink("lumberjack", func(*url.URL) (zap.Sink, error) {
		return lumberjackSink{
			Logger: &lumbLogger,
		}, nil
	})
}

func NewLogger(traceData TraceData) (Logger, error) {
	loggerConfig := configs.GetLoggerConfig()
	ld := loggerData{
		RequestId:     uuid.NewString(),
		TransactionId: traceData.TransactionId,
		SessionId:     traceData.SessionId,
		ProductId:     loggerConfig.ProductId,
		ChannelId:     traceData.ChannelId,
		ConsumerName:  traceData.ConsumerName,
		Environment:   loggerConfig.Env,
		Commerce:      loggerConfig.Commerce,
		Host: hostData{
			ContainerId: loggerConfig.ContainerId,
			Name:        loggerConfig.HostName,
			Ip:          loggerConfig.HostIp,
		},
	}

	zapLogger, err := newZapLogger(ld, loggerConfig)

	return &ZapLogger{
		data:       ld,
		logHandler: zapLogger,
	}, err
}
