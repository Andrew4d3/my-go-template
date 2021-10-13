FROM hub.fif.tech/base/golang-dev:1.16-latest

WORKDIR /app

COPY . .

RUN go mod download && \
    go get github.com/go-delve/delve/cmd/dlv

CMD go build -v -o server ./cmd/server && \
    dlv --listen=:4000 --headless=true --api-version=2 --accept-multiclient exec ./server