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

	item, err := dynamoDB.GetArtist(request)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 404}, err
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

	message := "An artist has been editted to have the following attributes: " + string(response)
	HTMLBody := "<h1>Success</h1><p>An artist has been editted to have the following attributes: " + string(response) + "</p>"

	err = ses.SendEmail(HTMLBody)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 404}, err
	}

	return events.APIGatewayProxyResponse{Body: message, StatusCode: 200}, nil
}

func main() {
	lambda.Start(Handler)
}
