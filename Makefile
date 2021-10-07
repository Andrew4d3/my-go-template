PROJECT_NAME = template-go-api # Change this to real project name

host-build:
	go build -v -o $(PROJECT_NAME) ./cmd/server

host-run:
	go run ./cmd/server