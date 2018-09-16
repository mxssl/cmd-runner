BINARY_NAME=cmd-runner

.PHONY: all build clean lint test dep build-linux build-darwin

all: build

cross-compilation: build-linux build-darwin

build:
	go build -o ${BINARY_NAME} -v

clean:
	rm -f ${BINARY_NAME}

lint:
	golangci-lint run -v
	
test:
	go test -v ./...

dep:
	dep ensure

# Cross compilation
build-linux:
	env GOOS=linux GOARCH=amd64 go build -o ${BINARY_NAME}-linux-amd64 -v
build-darwin:
	env GOOS=darwin GOARCH=amd64 go build -o ${BINARY_NAME}-darwin-amd64 -v
