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

type recomendacion struct {
	ID      string `json:"id"`
	Cancion string `json:"song_name"`
	Artista string `json:"song_artist"`
	Foto    string `json:"s3_code"`
	User    string `json:"user"`
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	foto := request.QueryStringParameters["foto"]
	sess, err := session.NewSession()
	if err != nil {
		fmt.Println("error en sesion")
	}
	svc := dynamodb.New(sess)
	filt := expression.Name("s3_code").Equal(expression.Value(foto))
	proj := expression.NamesList(expression.Name("song_name"), expression.Name("s3_code"), expression.Name("user"), expression.Name("song_artist"), expression.Name("id"))
	expr, err := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()
	tableName := "Recomendaciones"
	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(tableName),
	}
	result, err := svc.Scan(params)
	if err != nil {
		fmt.Println("error en resultado de scan: " + err.Error())
	}
	item := recomendacion{}
	dynamodbattribute.UnmarshalMap(result.Items[0], &item)
	if err != nil {
		fmt.Println("error en parseo de resultad: " + err.Error())
	}
	resString := `{"foto":"` + item.Foto + `", "canciones":[{"nombre":"` + item.Cancion + `", "artista":"` + item.Artista + `"}`
	for i := 1; i < len(result.Items); i++ {
		item := recomendacion{}
		dynamodbattribute.UnmarshalMap(result.Items[i], &item)
		resString += `, {"nombre":"` + item.Cancion + `", "artista":"` + item.Artista + `"}`
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
