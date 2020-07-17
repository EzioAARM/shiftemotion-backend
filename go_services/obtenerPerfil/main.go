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

type profile struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
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
		TableName: aws.String("Usuarios"),
	}

	result, err2 := svc.GetItem(input)
	if err2 != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusNotFound,
		}, nil
	}
	item := profile{}
	dynamodbattribute.UnmarshalMap(result.Item, &item)
	cadena := `{"name":"` + item.Name + `", "email":"` + item.Email + `", "id":"` + item.ID + `", "status":"200"}`
	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf(cadena),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
