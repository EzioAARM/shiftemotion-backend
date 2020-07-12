package main

import (
	//"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/zmb3/spotify"
	"net/http"
)

const redirectURI = ""

type URL struct {
	Data       string
	statusCode int
}

var (
	auth  = spotify.NewAuthenticator(redirectURI, spotify.ScopeUserReadPrivate, spotify.ScopeUserLibraryRead)
	ch    = make(chan *spotify.Client)
	state = "abc123"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	auth.SetAuthInfo("d133bb8c721e476db214f1319fce2b11", "ecd0db16e5334a1f9d34a548c1624534")
	url := auth.AuthURL(state)
	return events.APIGatewayProxyResponse{
		Body:       "{Data:" + url + ", StatusCode:" + string(http.StatusOK) + "}",
		StatusCode: http.StatusOK,
	}, nil
}

func main() {
	lambda.Start(handler)
}
