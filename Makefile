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
	go run ./cmd/client/main.go send -k 8df4135ecefc9a4d054e2c596cd3f56432e683431b27216fea917b01c8ef1fee

test:
	go test -v ./...