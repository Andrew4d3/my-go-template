package configs

// LogLevel is a type for log level
type LogLevel string

const (
	// Debug is a log level
	Debug LogLevel = "debug"
	// Info is a log level
	Info LogLevel = "info"
	// Warn is a log level
	Warn LogLevel = "warn"
	// Error is a log level
	Error LogLevel = "error"
)

// LoggerConfig defines the logger configuration
type LoggerConfig struct {
	Env         string
	Level       string
	Filepath    string
	ServiceName string
	ProductID   string
	Commerce    string
	HostName    string
	HostIP      string
	ContainerID string
}

var getAppConfig = GetAppConfig

// GetLoggerConfig gets the logger configuration
func GetLoggerConfig() LoggerConfig {
	loggerConfig := LoggerConfig{}

	appConfig := getAppConfig()
	loggerConfig.Env = appConfig.Environment
	loggerConfig.Level = getConfigFromEnv("LOG_LEVEL", string(Info))
	loggerConfig.Filepath = getConfigFromEnv("LOG_FILE_PATH", "./logs/logfile.log")
	loggerConfig.ServiceName = appConfig.ServiceName
	loggerConfig.ProductID = appConfig.ProductID
	loggerConfig.Commerce = "SEGUROS"
	loggerConfig.HostName = appConfig.HostName
	loggerConfig.HostIP = appConfig.HostIP
	loggerConfig.ContainerID = appConfig.ContainerID

	return loggerConfig
}
