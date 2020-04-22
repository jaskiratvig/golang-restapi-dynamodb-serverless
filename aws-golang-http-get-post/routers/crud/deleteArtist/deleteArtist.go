package main

import (
	"aws-golang-http-get-post/dynamoDB"
	"aws-golang-http-get-post/ses"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// Handler function Using AWS Lambda Proxy Request
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	svc := dynamoDB.CreateDynamoDBClient()

	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"Name": {
				S: aws.String(request.PathParameters["Name"]),
			},
		},
		TableName: aws.String("Artists"),
	}

	_, err := svc.DeleteItem(input)
	if err != nil {
		fmt.Println("Got error calling DeleteItem:")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	message := fmt.Sprintf("Deleted artist: Name: %+v ", request.PathParameters["Name"])
	HTMLBody := "<h1>Success</h1><p>Artist " + request.PathParameters["Name"] + " has been deleted from the database.</p>"

	err = ses.SendEmail(HTMLBody)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 404}, err
	}

	return events.APIGatewayProxyResponse{Body: message, StatusCode: 200}, nil
}

func main() {
	lambda.Start(Handler)
}
