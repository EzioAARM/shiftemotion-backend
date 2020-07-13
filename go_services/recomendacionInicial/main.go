package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	test := request.QueryStringParameters["email"]
	fmt.Println(request.QueryStringParameters["email"])
	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf(test),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
