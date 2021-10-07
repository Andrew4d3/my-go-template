package configs

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
)

type AppConfig struct {
	PORT           string
	ENVIRONMENT    string
	CONTAINER_ID   string
	HOSTNAME       string
	HOST_IP        string
	PRODUCT_ID     string
	COUNTRY        string
	JWT_SECRET     string
	SERVICE_NAME   string
	MONGO_URI      string
	MONGO_DATABASE string
}

var appConfig *AppConfig

func getConfigFromEnv(envName string, defaultValue string) string {
	var envar string

	if envar = os.Getenv(envName); envar == "" {
		envar = defaultValue
	}

	return envar
}

func getIp() (string, error) {
	req, err := http.Get("http://ip-api.com/json/")
	if err != nil {
		return "", err
	}
	defer req.Body.Close()

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return "", err
	}

	var ip map[string]string
	json.Unmarshal(body, &ip)

	return ip["query"], nil
}

func GetAppConfig() AppConfig {
	if appConfig == nil {
		appConfig = &AppConfig{}

		appConfig.ENVIRONMENT = getConfigFromEnv("ENVIRONMENT", "dev")
		appConfig.PORT = getConfigFromEnv("PORT", "3000")

		var hostname = "no-hostname"
		if name, err := os.Hostname(); err == nil {
			hostname = name
		}

		appConfig.CONTAINER_ID = getConfigFromEnv("CONTAINER_ID", hostname)
		appConfig.HOSTNAME = getConfigFromEnv("HOSTNAME", hostname)

		var iP = "127.0.0.1"
		if _ip, err := getIp(); err == nil {
			iP = _ip
		}

		appConfig.HOST_IP = getConfigFromEnv("HOST_IP", iP)
		appConfig.PRODUCT_ID = "SEGUROS_CL" // Change this to the corresponding product Id
		appConfig.COUNTRY = "CL"
		appConfig.JWT_SECRET = getConfigFromEnv("JWT_SECRET", "")
		appConfig.SERVICE_NAME = os.Args[0]
		appConfig.MONGO_URI = getConfigFromEnv("MONGO_URI", "")
		appConfig.MONGO_DATABASE = getConfigFromEnv("MONGO_DATABASE", "")
	}

	return *appConfig
}
