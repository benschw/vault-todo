deps:
	go get -u ./...

test:
	go test ./...

build: 
	go build -o todo


