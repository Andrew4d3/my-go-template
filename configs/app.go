package configs

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
)

// AppConfig defines application config
type AppConfig struct {
	Port          string
	Environment   string
	ContainerID   string
	HostName      string
	HostIP        string
	ProductID     string
	Country       string
	JWTSecret     string
	ServiceName   string
	MongoURI      string
	MongoDatabase string
}

var appConfig *AppConfig

func getConfigFromEnv(envName string, defaultValue string) string {
	var envar string

	if envar = os.Getenv(envName); envar == "" {
		envar = defaultValue
	}

	return envar
}

var httpGet = http.Get
var ioultilReadAll = ioutil.ReadAll

func getIP() (string, error) {
	res, err := httpGet("http://ip-api.com/json/")
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := ioultilReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var ip map[string]string
	json.Unmarshal(body, &ip)

	return ip["query"], nil
}

var osHostname = os.Hostname
var _getIP = getIP
var osArgs = os.Args

// GetAppConfig gets application config
func GetAppConfig() AppConfig {
	if appConfig == nil {
		appConfig = &AppConfig{}

		appConfig.Environment = getConfigFromEnv("ENVIRONMENT", "dev")
		appConfig.Port = getConfigFromEnv("PORT", "3000")

		var hostname = "no-hostname"
		if name, err := osHostname(); err == nil {
			hostname = name
		}

		appConfig.ContainerID = getConfigFromEnv("CONTAINER_ID", hostname)
		appConfig.HostName = getConfigFromEnv("HOSTNAME", hostname)

		var iP = "127.0.0.1"
		if _ip, err := _getIP(); err == nil {
			iP = _ip
		}

		appConfig.HostIP = getConfigFromEnv("HOST_IP", iP)
		appConfig.ProductID = "SEGUROS_CL" // Change this to the corresponding product Id
		appConfig.Country = "CL"
		appConfig.JWTSecret = getConfigFromEnv("JWT_SECRET", "")
		appConfig.ServiceName = osArgs[0]
		appConfig.MongoURI = getConfigFromEnv("MONGO_URI", "")
		appConfig.MongoDatabase = getConfigFromEnv("MONGO_DATABASE", "")
	}

	return *appConfig
}
