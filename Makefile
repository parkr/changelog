REV:=$(shell git rev-parse HEAD)

all: build test bench run

testdeps:
	go get github.com/stretchr/testify/assert

dist:
	mkdir -p dist

build: dist
	go build
	go build -o dist/changelogger ./changelogger

test: testdeps
	go test -v ./...

bench: testdeps
	go test -bench=. -v ./...

run: build
	dist/changelogger -h || true
	dist/changelogger -file=testdata/History.markdown -out=dist/History-changelogger.markdown
	diff testdata/History.markdown dist/History-changelogger.markdown

docker-build:
	docker build -t parkr/changelog:$(REV) .
