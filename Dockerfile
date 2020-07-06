FROM golang:latest
WORKDIR /app/changelog
COPY . .
RUN set -ex && go test ./... && go install ./...
