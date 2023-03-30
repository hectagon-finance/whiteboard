build:
	go build -o ./bin/example

run1:
	go run main.go 8080 genesis

run2:
	go run main.go 9000 8080

run3:
	go run main.go 9001 9000

run4:
	go run main.go 9002 9001

client1:
	go run ./cmd/client/main.go send -k afe79f8f118feb344332622997fff34a9ef9d6706e5766d2624ad8d65cf88126

client2:
	go run ./cmd/client/main.go send -k hahaah

test:
	go test -v ./...