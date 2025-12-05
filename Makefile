.PHONY: all fmt test tidy vet

all: fmt vet test

fmt:
	gofmt -w ./

test:
	go test ./...

vet:
	go vet ./...

tidy:
	go mod tidy
