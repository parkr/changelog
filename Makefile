REV:=$(shell git rev-parse HEAD)

all: build test run

testdeps:
	go get github.com/stretchr/testify/assert

dist:
	mkdir -p dist

build: dist
	go build
	go build -o dist/changelogger ./changelogger

test: testdeps
	go test -v ./...

run: build
	dist/changelogger -h || true
	dist/changelogger -out=History.markdown

docker-build:
	docker build -t parkr/changelog:$(REV) .
