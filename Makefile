GIT_SHA := $(shell git rev-parse --short HEAD 2>/dev/null)
GIT_TAG := $(shell git describe --abbrev=0 HEAD 2>/dev/null)
LD_FLAGS := '-s -w \
	-X main.versionTag=$(GIT_TAG)-$(GIT_SHA) \
	-X main.bashCompletion=$(shell base64 -w0 bash-completion)'

build:
	go build -ldflags=$(LD_FLAGS) -trimpath ./cmd/...

install:
	go install ./cmd/...

test:
	go test ./...

run:
	go run ./cmd/...

update:
	go list -m -u all | while read line; do \
		{ test -n "$(echo $line | cut -d' ' -f3)" && echo $line; }; done \
	| sed 1d | while read line; do \
		PKG=$(echo "$line"|cut -d' ' -f1); \
		VERSION=$(echo "$line"|cut -d' ' -f3 ); \
		go get $PKG@$(echo $VERSION|sed -e 's/\[//' -e 's/\]//'); done
	go mod tidy
