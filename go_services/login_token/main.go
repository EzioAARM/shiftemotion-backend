package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type profile struct {
	ID    string `json:"id"`
	Name  string `json:"display_name"`
	Email string `json:"email"`
}

type jsonString struct {
	Token string `json:"token"`
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	body := request.Body
	now := time.Now()
	sec := now.Unix()
	id := strconv.FormatInt(sec, 10)
	var token jsonString
	var user profile
	errJ := json.Unmarshal([]byte(body), &token)
	if errJ != nil {
		return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf("Error parseando " + errJ.Error()),
			StatusCode: 500,
		}, nil
	}
	sess, err := session.NewSession()
	if err != nil {
		fmt.Println("error en sesion")
	}
	svc := dynamodb.New(sess)
	input := &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"id": {
				N: aws.String(id),
			},
			"token": {
				S: aws.String(token.Token),
			},
		},
		TableName: aws.String("PasswordsTokens"),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		/*return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf("Error insertando elemento " + err.Error()),
			StatusCode: 500,
		}, nil*/
	}

	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/me", nil)
	if err != nil {
		fmt.Println("Error reading request. ", err)
	}

	req.Header.Set("Authorization", "Bearer "+token.Token)

	client := &http.Client{Timeout: time.Second * 10}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error reading response. ", err)
	}
	defer resp.Body.Close()

	body2, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading body. ", err)
	}
	errJ2 := json.Unmarshal([]byte(body2), &user)
	if errJ2 != nil {
		fmt.Println("Error Parseando: ", err)
	}
	input2 := &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"id": {
				N: aws.String(user.ID),
			},
			"name": {
				S: aws.String(user.Name),
			},
			"email": {
				S: aws.String(user.Email),
			},
		},
		TableName: aws.String("Usuarios"),
	}

	_, err = svc.PutItem(input2)
	if err != nil {
		/*return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf("Error insertando elemento " + err.Error()),
			StatusCode: 500,
		}, nil*/
	}
	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf("este es el token: " + token.Token),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
