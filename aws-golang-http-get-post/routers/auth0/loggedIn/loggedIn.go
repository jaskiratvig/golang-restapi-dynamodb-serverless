package main

import (
	"aws-golang-http-get-post/dynamoDB"
	"aws-golang-http-get-post/models"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// Handler function Using AWS Lambda Proxy Request
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	svc := dynamoDB.CreateDynamoDBClient()

	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("SessionData"),
		Key: map[string]*dynamodb.AttributeValue{
			"ClientID": {
				S: aws.String("kf9yX2qaBa7J5tV1PtL5SpcdZ2GXHEo9"),
			},
		},
	})
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 400}, err
	}

	item := models.SessionData{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 400}, err
	}

	return events.APIGatewayProxyResponse{Body: "Hello " + fmt.Sprintf("%v", item.Profile["name"]), StatusCode: 200}, nil
}

func main() {
	lambda.Start(Handler)
}
