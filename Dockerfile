FROM golang:latest

WORKDIR /app

COPY .env go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main ./cmd/p2pbot

CMD ["./main"]
