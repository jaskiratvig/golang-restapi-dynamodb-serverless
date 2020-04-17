package main

import (
	"aws-golang-http-get-post/dynamoDB"
	"aws-golang-http-get-post/models"
	"aws-golang-http-get-post/ses"
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
		TableName: aws.String("Artists"),
		Key: map[string]*dynamodb.AttributeValue{
			"Name": {
				S: aws.String(request.PathParameters["Name"]),
			},
		},
	})
	if err != nil {
		message := fmt.Sprintf(err.Error())
		return events.APIGatewayProxyResponse{Body: message, StatusCode: 404}, err
	}

	item := models.Artist{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		message := fmt.Sprintf("Failed to unmarshal Record, %v", err)
		return events.APIGatewayProxyResponse{Body: message, StatusCode: 404}, err
	}

	if item.Name == "" {
		message := fmt.Sprintf("Could not find '" + item.Name)
		return events.APIGatewayProxyResponse{Body: message, StatusCode: 404}, err
	}

	//Generate message that want to be sent as body
	message := fmt.Sprintf("Found artist: ArtistID: %+v Name: %+v Subcategory: %+v Songs: %+v Domestic: %+v ", item.ArtistID, item.Name, item.Subcategory, item.Songs, item.Domestic)

	HTMLBody := "<h1>Success</h1><p> " + message + "</p>"

	return ses.SendEmail(HTMLBody, message)
}

func main() {
	lambda.Start(Handler)
}
