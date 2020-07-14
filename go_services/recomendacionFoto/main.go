package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"net/http"
)

type recomendacion struct {
	ID    string `json:"id"`
	Track string `json:"cancion"`
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id := request.QueryStringParameters["id"]
	sess, err := session.NewSession()
	if err != nil {
		fmt.Println("error en sesion")
	}
	svc := dynamodb.New(sess)

	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				N: aws.String(id),
			},
		},
		TableName: aws.String("Usuarios"),
	}

	result, err2 := svc.GetItem(input)
	if err2 != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusNotFound,
		}, nil
	}
	item := recomendacion{}
	dynamodbattribute.UnmarshalMap(result.Item, &item)
	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf("Esto es una prueb con GO"),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
