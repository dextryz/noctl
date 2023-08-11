SHELL := /usr/bin/env bash

build:
	go build -o ./ix ./*.go

fmt:
	go mod tidy -compat=1.17
	gofmt -l -s -w .

test:
	go test ./...
