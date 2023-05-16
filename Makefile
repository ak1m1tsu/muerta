GO=go
GOCOVER=$(GO) tool cover
GOBUILD=$(GO) build
GOTEST=$(GO) test

.PHONY: build
build:
	$(GOBUILD) -o ./bin/muerta ./cmd/muerta/

run: build
	./bin/muerta

docker-up:
	docker compose up --build -d

docker-down:
	docker compose down --rmi local

swagger:
	swag fmt && swag init -d ./cmd/muerta/,./internal/api/ -o ./docs

.PHONY: test
test:
	$(GOTEST) -v ./... -count=1

.PHONY: test/cover
test/cover:
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCOVER) -func=coverage.out
	$(GOCOVER) -html=coverage.out