.PHONY: build
build: lint
	go build -o ./bin/muerta ./cmd/muerta/

run: build
	./bin/muerta

lint:
	golangci-lint run ./...

containers-up:
	docker compose up --build -d

containers-down:
	docker compose down --rmi all

test:
	go test -v -race -coverprofile=c.out ./... \
	&& go tool cover -html=c.out \
	&& rm c.out

swagger:
	swag fmt && swag init -d ./cmd/muerta/,./internal/api/ -o ./internal/api/docs
