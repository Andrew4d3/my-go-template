FROM golang:1.16

WORKDIR /app

COPY . .

RUN go mod download && \
    go install github.com/go-delve/delve/cmd/dlv@latest

CMD go build -v -o server ./cmd/server && \
    dlv --listen=:4000 --headless=true --api-version=2 --accept-multiclient exec ./server