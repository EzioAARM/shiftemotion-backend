package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/go-querystring/query"
	"io/ioutil"
	"net/http"
	"strings"
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

type aToken struct {
	Access string `json:"access_token"`
}

type reqInput struct {
	Type  string `url:"grant_type"`
	Token string `url:"refresh_token"`
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	token := request.QueryStringParameters["token"]
	clisec := "Zjk0MGQ2MTk4OGE2NDg0ZmJkY2M5OGE1OTZkNDc5ZWM6OGZiMzA1ZjA3NzIzNGZhMjhmNjI5YThlYjFmMTI4MmQ="
	//aqui empieza la peticion para cambiar el refresh por el access
	var access aToken
	data := reqInput{"refresh_token", token}
	opt, _ := query.Values(data)
	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(opt.Encode()))
	if err != nil {
		fmt.Println("Error reading request. ", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Basic "+clisec)

	client := &http.Client{Timeout: time.Second * 10}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error reading response. ", err)
	}
	defer resp.Body.Close()
	body3, err2 := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading body. ", err2)
	}
	errJ3 := json.Unmarshal([]byte(body3), &access)
	if errJ3 != nil {
		fmt.Println("Error Parseando: ", errJ3)
	}
	//aqui termina
	var res recomendacion
	req2, err := http.NewRequest("GET", "https://api.spotify.com/v1/recommendations?limit=10&market=ES&seed_tracks=7ytR5pFWmSjzHJIeQkgog4%2C0VjIjW4GlUZAMYd2vXMi3b%2C7ju97lgwC2rKQ6wwsf9no9%2C1rgnBhdG2JDFTbYkYRZAku", nil)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf("Error en la Peticion a spotify: " + err.Error()),
			StatusCode: 500,
		}, nil
	}
	req2.Header.Set("Authorization", "Bearer "+access.Access)

	client2 := &http.Client{Timeout: time.Second * 10}

	resp2, err := client2.Do(req2)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf("Error en la Peticion: " + err.Error()),
			StatusCode: 500,
		}, nil
	}
	defer resp2.Body.Close()

	body2, err := ioutil.ReadAll(resp2.Body)
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
	jsonString := `{"tracks":[{"name":"` + res.Tracks[0].Name + `", "artist":"` + res.Tracks[0].Artist[0].Name + `", "id":"` + res.Tracks[0].ID + `"}`
	for i := 1; i < len(res.Tracks); i++ {
		jsonString += `,{"name":"` + res.Tracks[i].Name + `", "artist":"` + res.Tracks[i].Artist[0].Name + `", "id":"` + res.Tracks[i].ID + `"}`
	}
	jsonString += "]}"
	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf(jsonString),
		StatusCode: 200,
	}, nil
}

func main() {

	lambda.Start(handler)
}
