.PHONY: build test clean run lint fmt tidy

BINARY=dab-cloudcost
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS=-ldflags "-X github.com/amayabdaniel/dab-cloudcost/internal/cmd.version=$(VERSION)"

build:
	go build $(LDFLAGS) -o bin/$(BINARY) ./cmd/dab-cloudcost

run: build
	./bin/$(BINARY)

test:
	go test -v -race -cover ./...

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

lint:
	golangci-lint run

fmt:
	go fmt ./...

tidy:
	go mod tidy

clean:
	rm -rf bin/ coverage.out coverage.html

install:
	go install $(LDFLAGS) ./cmd/dab-cloudcost
