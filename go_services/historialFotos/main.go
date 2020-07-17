package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

// Item Create struct to hold info about new item
type Item struct {
	ID      int    `json:"id"`
	Email   string `json:"email"`
	Emocion string `json:"emocion"`
	Foto    string `json:"foto"`
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	email := request.QueryStringParameters["email"]
	sess, err := session.NewSession()
	if err != nil {
		fmt.Println("error en sesion")
	}
	svc := dynamodb.New(sess)

	filt := expression.Name("email").Equal(expression.Value(email))
	proj := expression.NamesList(expression.Name("emocion"), expression.Name("foto"))
	expr, err := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()
	tableName := "HistorialFotosEmociones"
	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(tableName),
	}
	result, err := svc.Scan(params)
	resString := `{data:[`
	for _, i := range result.Items {
		item := Item{}
		dynamodbattribute.UnmarshalMap(i, &item)
		resString += `{"foto":"`+item.Foto+`", "emocion":"`+item.Emocion+`"},`
	}

	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf(""),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
