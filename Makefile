volumes = -v $(shell pwd)/cmd:/app/cmd -v $(shell pwd)/configs:/app/configs -v $(shell pwd)/internal:/app/internal -v $(shell pwd)/mocks:/app/mocks -v $(shell pwd)/logs:/app/logs
# Change this to real project name
PROJECT_NAME = template-go-api

host-build:
	go build -v -o server ./cmd/server

host-run:
	go run ./cmd/server

test:
	go test -v ./...

lint:
	golangci-lint run

docker-dev-build:
	docker build -f docker/dev/Dockerfile -t $(PROJECT_NAME):dev .

docker-dev-run:
	docker run --rm -it $(volumes) -p 3000:3000 --env-file ./.env $(PROJECT_NAME):dev

docker-dev-sh:
	docker run --rm -it $(volumes) -p 3000:3000 --env-file ./.env $(PROJECT_NAME):dev sh

docker-debug-build:
	docker build -f docker/debug/Dockerfile -t $(PROJECT_NAME):debug .

docker-debug-run:
	docker run --rm -it $(volumes) -p 3000:3000 -p 4000:4000 --security-opt=seccomp:unconfined --env-file ./.env $(PROJECT_NAME):debug