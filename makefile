.PHONY: all build justrun run test fmt

all: ;

build:
	go build -o bin/app cmd/app/main.go

justrun:
	bin/app $(ARG)

run: build justrun

test:
	go test ./...

fmt: 
	go fmt ./...