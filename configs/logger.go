package configs

type LogLevel string

const (
	Debug LogLevel = "debug"
	Info  LogLevel = "info"
	Warn  LogLevel = "warn"
	Error LogLevel = "error"
)

type LoggerConfig struct {
	Env         string
	Level       string
	Filepath    string
	ServiceName string
	ProductId   string
	Commerce    string
	HostName    string
	HostIp      string
	ContainerId string
}

func GetLoggerConfig() LoggerConfig {
	loggerConfig := LoggerConfig{}

	appConfig := GetAppConfig()
	loggerConfig.Env = appConfig.ENVIRONMENT
	loggerConfig.Level = getConfigFromEnv("LOG_LEVEL", string(Info))
	loggerConfig.Filepath = getConfigFromEnv("LOG_FILE_PATH", "./logs/logfile.log")
	loggerConfig.ServiceName = appConfig.SERVICE_NAME
	loggerConfig.ProductId = appConfig.PRODUCT_ID
	loggerConfig.Commerce = "SEGUROS"
	loggerConfig.HostName = appConfig.HOSTNAME
	loggerConfig.HostIp = appConfig.HOST_IP
	loggerConfig.ContainerId = appConfig.CONTAINER_ID

	return loggerConfig
}
