GO=go
GOCOVER=$(GO) tool cover
GOBUILD=$(GO) build
GOTEST=$(GO) test

.PHONY: build
build: format
	$(GOBUILD) -o ./bin/muerta ./cmd/muerta/

run: build
	./bin/muerta

format:
	gofumpt -w ./internal

docker-up:
	docker compose up --build -d

docker-down:
	docker compose down --rmi local

swagger:
	swag fmt && swag init -d ./cmd/muerta/,./internal/api/ -o ./internal/api/docs

.PHONY: test
test:
	$(GOTEST) -v -count=1 ./...

test100:
	$(GOTEST) -v -count=100 ./...

race:
	$(GOTEST) -v -race -count=1 ./...

.PHONY: cover
cover:
	$(GOTEST) -short -count=1 -race -coverprofile=coverage.out ./...
	$(GOCOVER) -html=coverage.out
	rm coverage.out
