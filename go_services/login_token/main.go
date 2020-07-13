package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	//"github.com/aws/aws-sdk-go/aws/session"
	//"github.com/aws/aws-sdk-go/service/dynamodb"
	//"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"os"
	//"strconv"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	//token := request.token
	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf("Prueba desde AWS Lambda con GO"),
		StatusCode: 200,
	}, nil
}

func main() {
	fmt.Println(os.Getenv("AWS_SECRET_ACCESS_KEY"))
	lambda.Start(handler)
}
