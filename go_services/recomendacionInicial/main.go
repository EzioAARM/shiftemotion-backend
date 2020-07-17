package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
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

type refresh struct {
	ID    string `json:"id"`
	Token string `json:"token"`
}

type test struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type arrayTest struct {
	Array []test `json:"items"`
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	email := request.QueryStringParameters["email"]

	sess, err := session.NewSession()
	if err != nil {
		fmt.Println("error en sesion")
	}
	svc := dynamodb.New(sess)

	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(email),
			},
		},
		TableName: aws.String("PasswordsTokens"),
	}

	result, err2 := svc.GetItem(input)
	if err2 != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusNotFound,
		}, nil
	}
	item := refresh{}
	dynamodbattribute.UnmarshalMap(result.Item, &item)

	clisec := "Zjk0MGQ2MTk4OGE2NDg0ZmJkY2M5OGE1OTZkNDc5ZWM6OGZiMzA1ZjA3NzIzNGZhMjhmNjI5YThlYjFmMTI4MmQ="
	//aqui empieza la peticion para cambiar el refresh por el access
	var access aToken
	data := reqInput{"refresh_token", item.Token}
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
	//aqui busca gustos personales
	req4, err := http.NewRequest("GET", "https://api.spotify.com/v1/me/top/tracks?time_range=medium_term&limit=10&offset=0", nil)
	if err != nil {
		fmt.Println(err.Error())
	}
	req4.Header.Set("Authorization", "Bearer "+access.Access)
	req4.Header.Set("Content-Type", "application/json")
	req4.Header.Set("Accept", "application/json")

	client4 := &http.Client{Timeout: time.Second * 10}

	resp4, err := client4.Do(req4)
	if err != nil {
		fmt.Println(resp4)
	}
	//fmt.Println(resp2)
	defer resp4.Body.Close()
	var prueba arrayTest
	body4, err := ioutil.ReadAll(resp4.Body)
	if err != nil {
		fmt.Println(err.Error())
	}
	errJ4 := json.Unmarshal([]byte(body4), &prueba)
	if errJ4 != nil {
		fmt.Println(errJ4.Error())
	}
	seeds := prueba.Array[0].ID + "%2C" + prueba.Array[1].ID + "%2C" + prueba.Array[2].ID + "%2C" + prueba.Array[3].ID
	//aqui terminan los gustos personales
	var res recomendacion
	req2, err := http.NewRequest("GET", "https://api.spotify.com/v1/recommendations?limit=10&seed_tracks="+seeds, nil)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf("Error en la Peticion a spotify: " + err.Error()),
			StatusCode: 200,
		}, nil
	}
	req2.Header.Set("Authorization", "Bearer "+access.Access)

	client2 := &http.Client{Timeout: time.Second * 10}

	resp2, err := client2.Do(req2)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf("Error en la Peticion: " + err.Error()),
			StatusCode: 200,
		}, nil
	}
	defer resp2.Body.Close()

	body2, err := ioutil.ReadAll(resp2.Body)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf("Error Leyendo respuesta de spotify: " + err.Error()),
			StatusCode: 200,
		}, nil
	}
	errJ2 := json.Unmarshal([]byte(body2), &res)
	if errJ2 != nil {
		return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf("Error Parseando: " + errJ2.Error()),
			StatusCode: 200,
		}, nil
	}
	jsonString := `{"tracks":[{"name":"` + res.Tracks[0].Name + `", "artist":"` + res.Tracks[0].Artist[0].Name + `", "id":"` + res.Tracks[0].ID + `"}`
	for i := 1; i < len(res.Tracks); i++ {
		jsonString += `,{"name":"` + res.Tracks[i].Name + `", "artist":"` + res.Tracks[i].Artist[0].Name + `", "id":"` + res.Tracks[i].ID + `"}`
	}
	jsonString += `], "status":"200"}`
	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf(jsonString),
		StatusCode: 200,
	}, nil
}

func main() {

	lambda.Start(handler)
}
