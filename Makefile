.PHONY: build test fmt vet check

build:
	go build -o bin/dbinstall ./cmd/dbinstall

test:
	go test ./...

fmt:
	gofmt -w $$(find . -name '*.go' -not -path './vendor/*')

vet:
	go vet ./...

check: test vet
