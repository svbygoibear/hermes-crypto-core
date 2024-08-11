.PHONY: build build-windows

build:
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o bin/bootstrap -tags lambda.norpc main.go
	zip -jrm bin/bootstrap.zip bin/bootstrap

build-windows:
	set GOOS=linux
	set GOARCH=arm64
	go build -tags lambda.norpc -o bin/bootstrap main.go
	tar -acf bin/bootstrap.zip bin/bootstrap