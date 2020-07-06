.PHONY: build

build:
	go get fmt
	go get github.com/aws/aws-lambda-go/events
	go get github.com/aws/aws-lambda-go/lambda
	GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags '-s -w' -a -installsuffix cgo -o ./go_services/*
	npm --prefix ./npm_services/* install ./npm_services/*
