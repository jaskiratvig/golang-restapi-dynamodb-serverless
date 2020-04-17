package main

import (
	"aws-golang-http-get-post/auth"
	"aws-golang-http-get-post/dynamoDB"
	"aws-golang-http-get-post/models"
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/coreos/go-oidc"
)

// Handler function Using AWS Lambda Proxy Request
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	state, err := getStateDynamoDB()
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 400}, err
	}
	if state == "" {
		err = fmt.Errorf("Could not find state")
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 400}, err
	}

	if state != request.QueryStringParameters["state"] {
		err = fmt.Errorf("Invalid state parameter")
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 400}, err
	}

	idToken, err := authenticate(request)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 400}, err
	}

	var profile map[string]interface{}
	if err := idToken.Claims(&profile); err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 400}, err
	}

	err = saveProfileDynamoDB(profile)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 400}, err
	}

	return events.APIGatewayProxyResponse{Body: request.QueryStringParameters["state"], StatusCode: 200}, nil
}

func getStateDynamoDB() (string, error) {
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
		return "", err
	}

	item := models.SessionData{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		return "", err
	}

	if item.State == "" {
		return "", fmt.Errorf("Could not find '" + item.State)
	}

	return item.State, nil
}

func authenticate(request events.APIGatewayProxyRequest) (*IDToken, error) {
	authenticator, err := auth.NewAuthenticator()
	if err != nil {
		return nil, err
	}

	token, err := authenticator.Config.Exchange(context.TODO(), request.QueryStringParameters["code"])
	if err != nil {
		err = fmt.Errorf("No token found: ", err)
		return nil, err
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		err = fmt.Errorf("No id_token field in oauth2 token")
		return nil, err
	}

	oidcConfig := &oidc.Config{
		ClientID: os.Getenv("AUTH0_CLIENT_ID"),
	}

	idToken, err := authenticator.Provider.Verifier(oidcConfig).Verify(context.TODO(), rawIDToken)
	if err != nil {
		err = fmt.Errorf("Failed to verify ID Token: " + err.Error())
		return nil, err
	}

	return idToken, nil
}

func saveProfileDynamoDB(profile map[string]interface{}) error {

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
		return err
	}

	item := models.SessionData{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		return err
	}

	// BodyRequest will be used to take the json response from client and build it
	bodyRequest := models.SessionData{
		ClientID: item.ClientID,
		State:    item.State,
		Profile:  profile,
	}

	// Unmarshal the json, return 404 if error
	err = json.Unmarshal([]byte(request.Body), &bodyRequest)
	if err != nil {
		return err
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
		return err
	}

	return nil
}

func main() {
	lambda.Start(Handler)
}
