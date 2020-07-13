package main

import (
	"fmt"

	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type profile struct {
	id    int
	user  string
	email string
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	body := request.Body
	var token string
	json.Unmarshal([]byte(body), &token)
	sess, err := session.NewSession()
	if err != nil {
		fmt.Println("error en sesion")
	}
	svc := dynamodb.New(sess)
	input := &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"id": {
				N: aws.String("1"),
			},
			"nombre": {
				S: aws.String("jose"),
			},
			"email": {
				S: aws.String("romaj1805@gmail.com"),
			},
		},
		TableName: aws.String("Usuarios"),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf("Error insertando elemento"),
			StatusCode: 500,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf("Guardado"),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
