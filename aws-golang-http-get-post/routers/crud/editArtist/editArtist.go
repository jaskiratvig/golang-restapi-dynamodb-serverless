package main

import (
	"aws-golang-http-get-post/dynamoDB"
	"aws-golang-http-get-post/models"
	"aws-golang-http-get-post/ses"
	"encoding/json"
	"fmt"
	"os"

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

	// BodyRequest will be used to take the json response from client and build it
	bodyRequest := models.Artist{
		ArtistID:    item.ArtistID,
		Name:        item.Name,
		Songs:       item.Songs,
		Subcategory: item.Subcategory,
		Domestic:    item.Domestic,
	}

	// Unmarshal the json, return 404 if error
	err = json.Unmarshal([]byte(request.Body), &bodyRequest)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 404}, err
	}

	av, err := dynamodbattribute.MarshalMap(bodyRequest)
	if err != nil {
		fmt.Println("Got error marshalling new movie item:")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("Artists"),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		fmt.Println("Got error calling PutItem:")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// Marshal the response into json bytes, if error return 404
	response, err := json.Marshal(&bodyRequest)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 404}, err
	}

	HTMLBody := "<h1>Success</h1><p>An artist has been editted to have the following attributes: " + string(response) + "</p>"

	return ses.SendEmail(HTMLBody, string(response))
}

func main() {
	lambda.Start(Handler)
}