package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"io/ioutil"
	"net/http"
	"time"
)

type recomendacion struct {
	Tracks []tracks `json:"tracks"`
}

type tracks struct {
	Name   string   `json:"name"`
	Artist []artist `json:"artists"`
	ID     string   `json:"id"`
}

type artist struct {
	Name string `json:"name"`
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	token := request.QueryStringParameters["token"]
	var res recomendacion
	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/recommendations?limit=10&market=ES&seed_tracks=7ytR5pFWmSjzHJIeQkgog4%2C0VjIjW4GlUZAMYd2vXMi3b%2C7ju97lgwC2rKQ6wwsf9no9%2C1rgnBhdG2JDFTbYkYRZAku", nil)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf("Error en la Peticion a spotify: " + err.Error()),
			StatusCode: 500,
		}, nil
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{Timeout: time.Second * 10}

	resp, err := client.Do(req)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf("Error en la Peticion: " + err.Error()),
			StatusCode: 500,
		}, nil
	}
	defer resp.Body.Close()

	body2, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf("Error Leyendo respuesta de spotify: " + err.Error()),
			StatusCode: 500,
		}, nil
	}
	errJ2 := json.Unmarshal([]byte(body2), &res)
	if errJ2 != nil {
		return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf("Error Parseando: " + errJ2.Error()),
			StatusCode: 500,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf("%+v", res),
		StatusCode: 200,
	}, nil
}

func main() {

	lambda.Start(handler)
}
