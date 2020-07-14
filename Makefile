.PHONY: build

build:
	go get fmt
	go get github.com/aws/aws-lambda-go/events
	go get github.com/aws/aws-lambda-go/lambda
	go get github.com/aws/aws-sdk-go/aws
	go get github.com/aws/aws-sdk-go/aws/session
	go get github.com/aws/aws-sdk-go/service/dynamodb
	go get github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute
	go get github.com/zmb3/spotify
	go get net/http
	go get os
	go get strings
	go get encoding/json
	go get time
	go get strconv
	go get io/ioutil
	go get github.com/dgrijalva/jwt-go
	GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags '-s -w' -a -installsuffix cgo -o ./go_services/historialFotos/main ./go_services/historialFotos
	GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags '-s -w' -a -installsuffix cgo -o ./go_services/login/main ./go_services/login
	GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags '-s -w' -a -installsuffix cgo -o ./go_services/login_token/main ./go_services/login_token
	GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags '-s -w' -a -installsuffix cgo -o ./go_services/obtenerPerfil/main ./go_services/obtenerPerfil
	GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags '-s -w' -a -installsuffix cgo -o ./go_services/recomendacionFoto/main ./go_services/recomendacionFoto
	GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags '-s -w' -a -installsuffix cgo -o ./go_services/recomendacionInicial/main ./go_services/recomendacionInicial
	npm --prefix ./npm_services/* install ./npm_services/*
