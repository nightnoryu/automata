all: build

.PHONY: modules
modules:
	go mod tidy

.PHONY: build
build: modules
	go build -o bin/ ./cmd/lab1
