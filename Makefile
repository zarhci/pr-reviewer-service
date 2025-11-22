build:
	go build -o app ./cmd

run:
	go run cmd/main.go

clean:
	rm -f app

docker-up:
	docker-compose up --build