BIN = forego
SRC = $(shell find . -name '*.go' -not -path './vendor/*')
VERSION:=$(shell git describe --tags)
LDFLAGS:=-X main.Version=$(VERSION)
# https://stackoverflow.com/a/58185179
LDFLAGS_EXTRA=-linkmode external -w -extldflags "-static"

.PHONY: all build clean lint release test

all: build

build: $(BIN)

clean:
	go mod tidy
	rm -f $(BIN)

lint: $(SRC)
	go fmt

release:
	bin/release

test: lint build
	go test -v -race -cover ./...

$(BIN): $(SRC)
	go build -ldflags "${LDFLAGS} ${LDFLAGS_EXTRA}" -o $@
