PROJECT_NAME = template-go-api # Change this to real project name

host-build:
	go build -i -v -o $(PROJECT_NAME) ./cmd/server