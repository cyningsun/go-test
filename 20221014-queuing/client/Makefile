BINARY_NAME=hello-world-cli

build:
	GOARCH=amd64 GOOS=darwin go build -o ${BINARY_NAME}-darwin main.go

run:
	./${BINARY_NAME}-darwin

brun: build run

clean:
	go clean
	rm ${BINARY_NAME}-darwin