FROM golang:1.22.2

COPY . .

RUN go mod download

RUN go build -o bot ./cmd/main.go

CMD ["./bot"]