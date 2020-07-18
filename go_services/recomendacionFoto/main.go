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
	ID     string    `json:"id"`
	Tracks []cancion `json:"cancion"`
	Foto   string    `json:"foto"`
}

type cancion struct {
	Name   string `json:"nombre"`
	Artist string `json:"artista"`
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
		TableName: aws.String("HistorialEscuchadas"),
	}

	result, err2 := svc.GetItem(input)
	if err2 != nil {
		return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf(`{"Status":"500"}`),
			StatusCode: http.StatusOK,
		}, nil
	}
	item := recomendacion{}
	dynamodbattribute.UnmarshalMap(result.Item, &item)
	resString := `{"id":"` + item.ID + `", "foto":"` + item.Foto + `", "canciones":[{"nombre":"` + item.Tracks[0].Name + `", "artista":"` + item.Tracks[0].Artist + `"}`
	for i := 1; i < len(item.Tracks); i++ {
		resString += `, {"nombre":"` + item.Tracks[i].Name + `", "artista":"` + item.Tracks[i].Artist + `"}`
	}
	resString += `], "status":"200"}`
	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf(resString),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
