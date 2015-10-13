all: build test run

build:
	go build

test:
	go test ./...

run:
	go run changelogger/changelogger.go
