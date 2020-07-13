package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"os"
	//"strings"
)

type profile struct {
	user  string
	email string
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	email := request.QueryStringParameters["email"]
	//user := strings.Split(email, "@")
	sess, err := session.NewSession()

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
		}, nil
	}
	svc := dynamodb.New(sess)

	tempProfile := profile{
		user:  "roma",
		email: "romaj1805@gmail.com",
	}
	av, err2 := dynamodbattribute.MarshalMap(tempProfile)
	if err2 != nil {
		fmt.Println("Got error marshalling new movie item:")
		fmt.Println(err.Error())
		os.Exit(1)
	}
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("Usuarios"),
	}
	_, err3 = svc.PutItem(input)
	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf("Guardado"),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
