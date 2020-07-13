package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"strconv"
	"time"
)

type profile struct {
	id    int
	user  string
	email string
}

type jsonString struct {
	token string `json:"token"`
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	body := request.Body
	now := time.Now()
	sec := now.Unix()
	id := strconv.FormatInt(sec, 10)
	token := jsonString{}
	json.Unmarshal([]byte(body), &token)
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
				S: aws.String(token.token),
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

	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf("este es el token: " + token.token),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
