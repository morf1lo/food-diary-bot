run:
	go run cmd/main.go

build:
	docker build -t food-diary-bot .

run-docker:
	docker run --name food-diary-bot --rm food-diary-bot