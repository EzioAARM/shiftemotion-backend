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
	ID      string `json:"id"`
	Email   string `json:"user_id"`
	Emocion string `json:"emotion"`
	Foto    string `json:"picture_code"`
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	email := request.QueryStringParameters["email"]
	sess, err := session.NewSession()
	if err != nil {
		fmt.Println("error en sesion")
	}
	svc := dynamodb.New(sess)

	filt := expression.Name("user_id").Equal(expression.Value(email))
	proj := expression.NamesList(expression.Name("emotion"), expression.Name("picture_code"), expression.Name("user_id"), expression.Name("id"))
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
	itemInit := Item{}
	dynamodbattribute.UnmarshalMap(result.Items[0], &itemInit)
	resString := `{data:[{"foto":"` + itemInit.Foto + `", "emocion":"` + itemInit.Emocion + `", "id":"` + itemInit.ID + `"}`
	for i := 1; i < len(result.Items); i++ {
		item := Item{}
		dynamodbattribute.UnmarshalMap(result.Items[i], &item)
		resString += `, {"foto":"` + item.Foto + `", "emocion":"` + item.Emocion + `", "id":"` + item.ID + `"}`
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
