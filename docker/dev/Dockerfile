FROM golang:1.16

WORKDIR /app

COPY . .

RUN go mod download

CMD ["go", "run", "./cmd/server"]
