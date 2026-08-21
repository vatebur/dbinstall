.PHONY: build test test-race fmt vet check

build:
	go build -o bin/dbinstall ./cmd/dbinstall

test:
	go test ./...

test-race:
	go test -race ./...

fmt:
	gofmt -w $$(find . -name '*.go' -not -path './vendor/*')

vet:
	go vet ./...

check: test-race vet build
