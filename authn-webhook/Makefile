.DEFAULT: all

.PHONY: all
all: build

.PHONY: build
build:
	echo "building authn-webhook binary"
	mkdir -p bin/amd64
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/amd64 .

.PHONY: run
run:
	./bin/amd64/authn-webhook

.PHONY: clean
clean:
	rm -rf ./bin
