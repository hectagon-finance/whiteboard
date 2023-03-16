build:
	go build -o ./bin/example

run: build
	./bin/example

test:
	go test -v ./...