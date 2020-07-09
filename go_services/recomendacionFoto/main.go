package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// Item Create struct to hold info about new item
type Item struct {
	id      int
	name    string
	cancion string
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	svc := dynamodb.New(sess)
	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("HistorialFotosEmociones"),
		Key: map[string]*dynamodb.AttributeValue{
			"nombre": {
				N: aws.String("roma"),
			},
		},
	})
	if err != nil {
		fmt.Println(err.Error)
	}
	item := Item{}
	dynamodbattribute.UnmarshalMap(result.Item, &item)
	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf(item.cancion),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
