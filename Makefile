all-example: build-example run-example

build-example:
	@go build -o example ./cmd

run-example:
	@./example
