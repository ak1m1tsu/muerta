build:
	go build -o ./bin/muerta ./cmd/muerta/

run: build
	./bin/muerta

docker-up:
	docker compose up --build -d

docker-down:
	docker compose down --rmi local

swagger:
	swag fmt && swag init -d ./cmd/muerta/,./internal/api/ -o ./docs

postman:
	go run ./cmd/postman/main.go

test:
	go test -v ./... -count=1
