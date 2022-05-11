package main

import (
	"fmt"
	"os"
	"template-go-api/configs"
	"template-go-api/internal/app/middlewares"
	"template-go-api/internal/app/routes"
	"template-go-api/internal/pkg/customerror"
	"template-go-api/internal/pkg/drivers/mongodb"
	"template-go-api/internal/pkg/logger"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func getPort() string {
	var PORT string
	if PORT = os.Getenv("PORT"); PORT == "" {
		PORT = "3000"
	}

	return PORT
}

func panicError(err error) {
	if err != nil {
		fmt.Printf("Unexpected error: %s\nExiting program...", err)
		panic(err)
	}
}

func connectToDB(mainLogger logger.Logger) {
	mainLogger.Info("Connecting to Mongo DB")
	appConfigs := configs.GetAppConfig()

	if err := mongodb.ConnectDB(appConfigs.MongoURI, appConfigs.MongoDatabase); err != nil {
		mainLogger.Error(customerror.NewFromErrror(err))
		// Panic here if you are required to have a DB connection for running your app
		return
	}
	mainLogger.Info("Mongo connection established")
}

func main() {
	godotenv.Load()

	err := logger.SetupFileStore()
	panicError(err)

	mainLogger, err := logger.NewLogger(logger.TraceData{})
	panicError(err)

	// Uncomment this line when you get your DB connection ready
	// connectToDB(mainLogger)

	echoServer := echo.New()
	middlewares.SetMiddlewares(echoServer)
	routes.SetRoutes(echoServer)

	port := configs.GetAppConfig().Port
	mainLogger.Info("Starting server in port: " + port)
	if err := echoServer.Start(":" + port); err != nil {
		mainLogger.Error(customerror.NewFromErrror(err))
	}

}
