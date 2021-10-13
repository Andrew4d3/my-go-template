FROM hub.fif.tech/base/golang-dev:1.16-latest

WORKDIR /app

COPY . .

RUN go mod download

CMD ["go", "run", "./cmd/server"]
