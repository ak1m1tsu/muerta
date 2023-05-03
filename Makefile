build:
	go build -o ./bin/muerta ./cmd/muerta/

run: build
	./bin/muerta

schema:
	go run ./cmd/postman/main.go

test:
	go test -v ./... -count=1
