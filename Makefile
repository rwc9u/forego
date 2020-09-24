BIN = forego
SRC = $(shell find . -name '*.go' -not -path './vendor/*')
VERSION:=$(shell git describe --tags)
LDFLAGS:=-X main.Version=$(VERSION)
# https://stackoverflow.com/a/58185179
LDFLAGS_EXTRA=-linkmode external -w -extldflags "-static"

.PHONY: all build clean lint release test dist dist-clean

all: build

build: $(BIN)

clean:
	go mod tidy
	rm -f $(BIN)

lint: $(SRC)
	go fmt

dist-clean:
	go mod tidy
	rm -rf dist release/
	rm -f forego

dist: dist-clean
	mkdir -p dist/alpine-linux/amd64 && GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS) ${LDFLAGS_EXTRA}" -a -tags netgo -installsuffix netgo -o dist/alpine-linux/amd64/forego
	mkdir -p dist/alpine-linux/arm64 && GOOS=linux GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -a -tags netgo -installsuffix netgo -o dist/alpine-linux/arm64/forego
	mkdir -p dist/alpine-linux/armhf && GOOS=linux GOARCH=arm GOARM=6 go build -ldflags "$(LDFLAGS)" -a -tags netgo -installsuffix netgo -o dist/alpine-linux/armhf/forego
	mkdir -p dist/linux/amd64 && GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS) ${LDFLAGS_EXTRA}" -o dist/linux/amd64/forego
	mkdir -p dist/linux/arm64 && GOOS=linux GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o dist/linux/arm64/forego
	mkdir -p dist/linux/i386  && GOOS=linux GOARCH=386 go build -ldflags "$(LDFLAGS)" -o dist/linux/i386/forego
	mkdir -p dist/linux/armel  && GOOS=linux GOARCH=arm GOARM=5 go build -ldflags "$(LDFLAGS)" -o dist/linux/armel/forego
	mkdir -p dist/linux/armhf  && GOOS=linux GOARCH=arm GOARM=6 go build -ldflags "$(LDFLAGS)" -o dist/linux/armhf/forego
	mkdir -p dist/darwin/amd64 && GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/darwin/amd64/forego
	mkdir -p dist/darwin/i386  && GOOS=darwin GOARCH=386 go build -ldflags "$(LDFLAGS)" -o dist/darwin/i386/forego

release: dist
	mkdir -p release
	tar -cvzf release/forego-alpine-linux-amd64-$(VERSION).tar.gz -C dist/alpine-linux/amd64 forego
	tar -cvzf release/forego-alpine-linux-arm64-$(VERSION).tar.gz -C dist/alpine-linux/arm64 forego
	tar -cvzf release/forego-alpine-linux-armhf-$(VERSION).tar.gz -C dist/alpine-linux/armhf forego
	tar -cvzf release/forego-linux-amd64-$(VERSION).tar.gz -C dist/linux/amd64 forego
	tar -cvzf release/forego-linux-arm64-$(VERSION).tar.gz -C dist/linux/arm64 forego
	tar -cvzf release/forego-linux-i386-$(VERSION).tar.gz -C dist/linux/i386 forego
	tar -cvzf release/forego-linux-armel-$(VERSION).tar.gz -C dist/linux/armel forego
	tar -cvzf release/forego-linux-armhf-$(VERSION).tar.gz -C dist/linux/armhf forego
	tar -cvzf release/forego-darwin-amd64-$(VERSION).tar.gz -C dist/darwin/amd64 forego
	tar -cvzf release/forego-darwin-i386-$(VERSION).tar.gz -C dist/darwin/i386 forego

test: lint build
	go test -v -race -cover ./...

$(BIN): $(SRC)
ifeq ("$(GOOS)-$(GOARCH)","linux-amd")
	go build -ldflags "${LDFLAGS} ${LDFLAGS_EXTRA}" -o $@
else
	go build -ldflags "${LDFLAGS}" -o $@
endif
