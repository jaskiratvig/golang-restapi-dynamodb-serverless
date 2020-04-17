package main

import (
	"aws-golang-http-get-post/auth"
	"aws-golang-http-get-post/dynamoDB"
	"aws-golang-http-get-post/models"
	"context"
	"fmt"
	"net/http"

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

	return events.APIGatewayProxyResponse{Headers: map[string]string{
		"location": "https://bhvn5rgkmd.execute-api.us-east-1.amazonaws.com/dev/loggedIn",
	}, StatusCode: http.StatusTemporaryRedirect}, nil
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

func authenticate(request events.APIGatewayProxyRequest) (*oidc.IDToken, error) {
	authenticator, err := auth.NewAuthenticator()
	if err != nil {
		return nil, err
	}

	token, err := authenticator.Config.Exchange(context.TODO(), request.QueryStringParameters["code"])
	if err != nil {
		err = fmt.Errorf("No token found: %+v", err)
		return nil, err
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		err = fmt.Errorf("No id_token field in oauth2 token")
		return nil, err
	}

	oidcConfig := &oidc.Config{
		ClientID: "kf9yX2qaBa7J5tV1PtL5SpcdZ2GXHEo9",
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

	av, err := dynamodbattribute.MarshalMap(bodyRequest)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("SessionData"),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	lambda.Start(Handler)
}
