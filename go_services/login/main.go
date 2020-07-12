package main

import (
	//"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/zmb3/spotify"
	"net/http"
	"os"
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
	id := os.Getenv("SPOTIFY_ID")
	secret := os.Getenv("SPOTIFY_SECRET")
	auth.SetAuthInfo(id, secret)
	url := auth.AuthURL(state)
	return events.APIGatewayProxyResponse{
		Body:       "{Data:" + url + ", StatusCode:" + string(http.StatusOK) + "}",
		StatusCode: http.StatusOK,
	}, nil
}

func main() {
	lambda.Start(handler)
}
