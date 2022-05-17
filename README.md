# My Go Template for building REST APIs

[[Project Description here]]

## Template Description

> Delete this section after initializing your project

### **How to use this template?**

Clone this repo and then delete the `.git` folder so that you can start a new git repository using: `git init`. Alternatively, you can use the "Use this template" feature from Github.

### **Stack**

-  Echo as web server (with complements)
-  Zap as logger library
-  Testify as testing library
-  JWT Integration through middleware

### **Logger**

This template includes a logger integration that you can use through the Echo context:

```go
import (
	"template-go-api/internal/app/middlewares"
)

// ...

func Controller(c echo.Context) error {
	cc:= c.(middlewares.CustomContext)
	logger := cc.CustomLogger()
	logger.Info("Hello")
}
```

Which will print the following log entry:

```
{"level":"info","date":"2022-05-12T08:37:25.716-0400","message":"Hello","channelId":"","consumerName":"","environment":"dev","host":{"containerId":"andrew-XPS-15-9560","ip":"186.00.00.00","name":"andrew-XPS-15-9560"},"requestId":"8368b925-0e3a-4964-81ef-6f2b7788f160","sessionId":"","transactionId":"","payload":{}}
```

Additionaly, If you're working with microservices, you can propagate correlation IDs in your log entries by using the following headers:

-  `x-transaction-id`: To identify an end-to-end call chain.
-  `x-session-id`: To identify multiple related call chains.
-  `x-channel-id`: To identify the channel that was used to start the call chain. (web, mobile app, etc).
-  `x-consumer`: To identify the caller service.

### **Custom Errors**

Optionally, you can use the CustomError package found at `internal/pkg/customerror`, if you want to provide more detailed errors in your logs:

```go
customError := customerror.New("Boom", customerror.InputData{
	Data: map[string]interface{}{ // Any interface here
		"foo": "bar",
	},
})
```

### **Authorization**

You can use the authorization middleware for controling access to your protected routes. Example:

```go
import (
	"template-go-api/internal/app/middlewares"
)


e.GET("/protected", controllers.HealthCheck, middlewares.AuthorizationMiddleware)
```

If a valid JWT is not provided as Bearer token in the authorization header, an HTTP 401 response like this will be returned:

```json
{
   "message": "Authorization header is not present"
}
```

## Requirements

-  Go 1.16 (\*)
-  Make
-  Docker (If you want to use Docker for development)
-  [golangci-lint](https://github.com/golangci/golangci-lint) for linter

(\*) For managing multiple Go versions, I recommend using [Go Version Manager](https://github.com/moovweb/gvm) (gvm)

## Instructions for running via Host

1. Create a `.env` file at the project's root:

```env
ENVIRONMENT=dev
PORT=3000
JWT_SECRET=<YOUR_SECRET_HERE>
# If you are using mongo, otherwise locate your DB secrets here
MONGO_URI=<YOUR_MONGO_URI>
MONGO_DATABASE=<YOUR_MONGO_DATABASE>
```

**Note**: In production you should inject your env variables directly into the procces

2. To start the dev server, run the make command:

```sh
make host-run
```

3. Now, the project should start at port 3000 or any other port indicated by the env variables.

## Instructions for running via Docker

1. Create the same env file described in the first step of the previous section.

2. Run the docker image using the make command:

```sh
make docker-dev-build
```

3. If the image was built successfully, run the following make command to start a container using the corresponding image:

```sh
make docker-dev-run
```

4. Now, the project should start at port 3000 or any other port indicated by the env variables.`

## Unit tests and Linter

For unit tests, just run:

```sh
make test
```

For linter(\*), run:

```sh
make lint
```

(\*) Remember to have golangci-lint installed

## Debugging using the VSCode editor

If you're using VSCode as main editor, you can leverage to use its poweful debugger. The only thing you need to do is to locate the following `launch.json` file into the `.vscode` folder at the project's root.

```json
{
   // Use IntelliSense to learn about possible attributes.
   // Hover to view descriptions of existing attributes.
   // For more information, visit: https://go.microsoft.com/fwlink/?linkid=830387
   "version": "0.2.0",
   "configurations": [
      {
         "name": "Launch Package",
         "type": "go",
         "request": "launch",
         "mode": "auto",
         "program": "cmd/server"
      }
   ]
}
```

The configurations above will tell VSCode to start a debugger from the main go file.

## Debugging using Docker

It's also possible to start a Debug server using docker. First run the following make command to build the corresponding Docker image:

```sh
make docker-debug-build
```

Now to start the debugg server, just run:

```sh
make docker-debug-run
```

A debug server should start through port 4000. You can connect to this server using any Debugger client. If you're using VSCode as Debugger client, you can use the following configuration in your `.vscode/launch.json` file:

```json
{
   // Use IntelliSense to learn about possible attributes.
   // Hover to view descriptions of existing attributes.
   // For more information, visit: https://go.microsoft.com/fwlink/?linkid=830387
   "version": "0.2.0",
   "configurations": [
      {
         "name": "Connect to server",
         "type": "go",
         "request": "attach",
         "mode": "remote",
         "remotePath": "/app",
         "port": 4000,
         "host": "127.0.0.1"
      }
   ]
}
```

## For Production

This template comes with a minimal Dockerfile located at `docker/deploy/Dockerfile` that you can use to deploy your application into production or any other environment you migh have. Feel free to extend it to include any other additional software you migh need to correctly run your Go application.

If you want to build the binary file using a CI/CD tool, you should use the following command:

```sh
go build -v -o server ./cmd/server
```

Additionaly, you should include in your pipeline the corresponding commands for testing:

```sh
go test -v ./...
```

And linting:

```sh
golint -set_exit_status  ./...
```
