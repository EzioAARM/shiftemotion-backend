.PHONY: build

build:
	sam build
	GOARCH=amd64 GOOS=linux go build -o ./go_services/*
	npm --prefix ./npm_services/* install ./npm_services/*
