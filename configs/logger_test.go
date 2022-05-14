package configs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetLoggerConfig(t *testing.T) {
	ogGetAppConfig := getAppConfig
	defer func() {
		getAppConfig = ogGetAppConfig
	}()

	t.Run("Should return a correct LoggerConfig object", func(t *testing.T) {
		getAppConfig = func() AppConfig {
			return AppConfig{
				Environment: "test",
				ServiceName: "test-service",
				HostName:    "test-host",
				HostIP:      "1.2.3.4",
				ContainerID: "12345",
			}
		}

		expectedObject := LoggerConfig{
			Env:         "test",
			Level:       "info",
			Filepath:    "./logs/logfile.log",
			ServiceName: "test-service",
			HostName:    "test-host",
			HostIP:      "1.2.3.4",
			ContainerID: "12345",
		}

		loggerConfig := GetLoggerConfig()
		areEqual := assert.ObjectsAreEqual(expectedObject, loggerConfig)
		assert.True(t, areEqual)
	})
}
