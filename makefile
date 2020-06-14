.PHONY: all build justrun run test

all: ;

build:
	go build -o bin/app cmd/app/main.go

justrun:
	bin/app

run: build justrun

test:
	go test ./...