SHELL=/bin/bash

default: build
clean:
	rm -rf build

deps:
	go get -u ./...

test:
	go test ./...

build: 
	mkdir -p build
	go build -o build/todo

.PHONY: build
