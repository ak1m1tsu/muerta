build:
	go build -o ./bin/muerta ./cmd/muerta/

run: build
	./bin/muerta

test:
	go test -v ./... -count=1
