build:
	go build -o wfv ./cmd/...	

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
		echo go get $PKG@$(echo $VERSION|sed -e 's/\[//' -e 's/\]//'); done
	go mod tidy
