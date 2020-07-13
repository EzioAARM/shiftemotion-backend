package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Println(request.QueryStringParameters["email"])
	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf("Prueba desde AWS Lambda con GO"),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
