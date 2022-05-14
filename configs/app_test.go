package configs

import (
	"errors"
	"io"
	"net/http"
	"os"
	"template-go-api/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

var bodyMock = new(mocks.IoReadCloser)
var httpGetSuccess = func(url string) (resp *http.Response, err error) {
	return &http.Response{
		Body: bodyMock,
	}, nil
}

var ioultilReadAllSucces = func(r io.Reader) ([]byte, error) {
	return []byte("{\"query\": \"201.189.26.204\"}"), nil
}

func Test_getConfigFromEnv(t *testing.T) {
	t.Run("Should get corresponding envar", func(t *testing.T) {
		os.Setenv("FOO", "BAR")

		envar := getConfigFromEnv("FOO", "NA")
		assert.Equal(t, "BAR", envar)
	})

	t.Run("Should get default value if envar does not exist", func(t *testing.T) {
		envar := getConfigFromEnv("BAR", "defaultValue")
		assert.Equal(t, "defaultValue", envar)
	})
}

func Test_getIp(t *testing.T) {
	originalHTTP := httpGet
	originalReadAll := ioultilReadAll

	defer func() {
		httpGet = originalHTTP
		ioultilReadAll = originalReadAll
	}()

	t.Run("Should throw error if service for getting public IP fails", func(t *testing.T) {
		httpGet = func(url string) (resp *http.Response, err error) {
			return nil, errors.New("Boom GetIP")
		}

		_, err := getIP()
		assert.Errorf(t, err, "Boom GetIP")
	})

	t.Run("Should throw error if there is a problem reading the response body", func(t *testing.T) {
		httpGet = httpGetSuccess
		ioultilReadAll = func(r io.Reader) ([]byte, error) {
			return nil, errors.New("Boom ReadBody")
		}
		bodyMock.On("Close").Return(nil)

		_, err := getIP()
		assert.Errorf(t, err, "Boom ReadBody")
	})

	t.Run("Should return the IP if the service response is OK", func(t *testing.T) {
		httpGet = httpGetSuccess
		ioultilReadAll = ioultilReadAllSucces
		bodyMock.On("Close").Return(nil)

		result, err := getIP()
		assert.NoError(t, err)
		assert.Equal(t, "201.189.26.204", result)
	})
}

func Test_GetAppConfig(t *testing.T) {
	originalHostNameFn := osHostname
	originalGetIP := _getIP

	defer func() {
		osHostname = originalHostNameFn
		_getIP = originalGetIP
	}()

	_getIPSuccess := func() (string, error) {
		return "1.2.3.4", nil
	}

	osHostnameSuccess := func() (name string, err error) {
		return "test-host", nil
	}

	setup := func() {
		appConfig = nil
		os.Setenv("HOSTNAME", "") // Workaround for pipeline
		osArgs = []string{"test-server"}
	}

	t.Run("Should return no-hostname as HOSTNAME if there is an error getting the hostname", func(t *testing.T) {
		setup()
		_getIP = _getIPSuccess
		osHostname = func() (name string, err error) {
			return "", errors.New("Boom Hostname")
		}

		appConfig := GetAppConfig()
		assert.Equal(t, "no-hostname", appConfig.HostName)
	})

	t.Run("Should return 127.0.0.1 as HOST_IP if there is an error getting the public IP", func(t *testing.T) {
		setup()
		osHostname = osHostnameSuccess
		_getIP = func() (string, error) {
			return "", errors.New("Boom IP")
		}

		appConfig := GetAppConfig()
		assert.Equal(t, "127.0.0.1", appConfig.HostIP)
	})

	t.Run("Should return the corresponding struct with all the correct values", func(t *testing.T) {
		setup()
		osHostname = osHostnameSuccess
		_getIP = _getIPSuccess
		os.Setenv("ENVIRONMENT", "test")
		os.Setenv("PORT", "3001")
		os.Setenv("CONTAINER_ID", "test-container")
		os.Setenv("HOSTNAME", "test-host")
		os.Setenv("HOST_IP", "10.1.2.3")
		os.Setenv("JWT_SECRET", "test-secret")
		os.Setenv("MONGO_URI", "mongo://test:2121")
		os.Setenv("MONGO_DATABASE", "test-db")

		expectedAppConfig := AppConfig{
			Port:          "3001",
			Environment:   "test",
			ContainerID:   "test-container",
			HostName:      "test-host",
			HostIP:        "10.1.2.3",
			JWTSecret:     "test-secret",
			ServiceName:   "test-server",
			MongoURI:      "mongo://test:2121",
			MongoDatabase: "test-db",
		}

		appConfig := GetAppConfig()

		areEqual := assert.ObjectsAreEqual(appConfig, expectedAppConfig)
		assert.True(t, areEqual)
	})
}
