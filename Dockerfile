FROM golang:1.22.2

COPY . .

RUN go mod download

RUN go build -o ./cmd/bot/bot ./cmd/bot/main.go

CMD ["./cmd/bot/bot"]
