package main

import (
	"aws-golang-http-get-post/auth"
	"aws-golang-http-get-post/dynamoDB"
	"aws-golang-http-get-post/models"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// Handler function Using AWS Lambda Proxy Request
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	state, err := generateRandomState()
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 400}, err
	}

	errorMessage := saveStateDynamoDB(state)
	if errorMessage != "" {
		err = fmt.Errorf(errorMessage)
		return events.APIGatewayProxyResponse{Body: errorMessage, StatusCode: 400}, err
	}

	authenticator, err := auth.NewAuthenticator()
	if err != nil {
		err = fmt.Errorf(errorMessage)
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 400}, err
	}

	return events.APIGatewayProxyResponse{Headers: map[string]string{
		"location": authenticator.Config.AuthCodeURL(state),
	}, StatusCode: http.StatusTemporaryRedirect}, nil
}

func generateRandomState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

func saveStateDynamoDB(state string) string {
	session := models.SessionData{
		ClientID: "kf9yX2qaBa7J5tV1PtL5SpcdZ2GXHEo9",
		State:    state,
		Profile:  nil,
	}

	svc := dynamoDB.CreateDynamoDBClient()

	av, err := dynamodbattribute.MarshalMap(session)
	if err != nil {
		return "Got error marshalling session data:" + err.Error()
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("SessionData"),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		return "Got error calling PutItem:" + err.Error()
	}

	return ""
}

func main() {
	lambda.Start(Handler)
}
