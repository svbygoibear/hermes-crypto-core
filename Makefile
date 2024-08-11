.PHONY: build

build:
	$ENV:GOOS="linux" go build -o bin/bootstrap main.go 