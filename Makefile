build:
	go build -o wfv ./cmd/...	

install:
	go install ./cmd/...

test:
	go test ./...

run:
	go run ./cmd/...

update:
	go list -m -u all \
	| awk -F" " '{ if ($$3 != "") print $$1 " " $$3; }' \
	| xargs -l bash -c 'VERSION=$(grep -Po "(?<=\[).+(?=\])" <<<$$1); go get $$0@$$VERSION'
	go mod tidy
