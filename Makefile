build:
	go build -o ./bin/whiteboard

run: build
	./bin/whiteboard

test:
	go test ./...