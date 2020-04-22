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
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/coreos/go-oidc"
	"github.com/dgrijalva/jwt-go"
)

// Handler function Using AWS Lambda Proxy Request
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	sess := session.New()
	svcSSM := ssm.New(sess)

	loggedInURL, err := svcSSM.GetParameter(
		&ssm.GetParameterInput{
			Name: aws.String("/dev/LoggedInURL"),
		},
	)

	state, err := getStateDynamoDB()
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 400}, err
	}

	if state != request.QueryStringParameters["state"] {
		err = fmt.Errorf("Invalid state parameter")
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 400}, err
	}

	idToken, err := authenticate(request, state)
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
		"location": aws.StringValue(loggedInURL.Parameter.Value),
	}, StatusCode: http.StatusTemporaryRedirect}, nil
}

func getStateDynamoDB() (string, error) {
	item, err := dynamoDB.GetSessionData()
	if err != nil {
		return "", err
	}

	if item.State == "" {
		return "", fmt.Errorf("Could not find '" + item.State)
	}

	return item.State, nil
}

func authenticate(request events.APIGatewayProxyRequest, state string) (*oidc.IDToken, error) {
	item, err := dynamoDB.GetSessionData()
	if err != nil {
		return nil, err
	}

	if err = jsonAuthenticate(item, state); err != nil {
		return nil, err
	}

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

	sess := session.New()
	svcSSM := ssm.New(sess)

	clientID, err := svcSSM.GetParameter(
		&ssm.GetParameterInput{
			Name: aws.String("/dev/ClientID"),
		},
	)

	oidcConfig := &oidc.Config{
		ClientID: aws.StringValue(clientID.Parameter.Value),
	}

	idToken, err := authenticator.Provider.Verifier(oidcConfig).Verify(context.TODO(), rawIDToken)
	if err != nil {
		err = fmt.Errorf("Failed to verify ID Token: " + err.Error())
		return nil, err
	}

	return idToken, nil
}

func jsonAuthenticate(item models.SessionData, state string) error {
	if item.Token != "" {
		token, err := jwt.Parse(item.Token, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				merror := fmt.Errorf("Error parsing JWT token")
				return nil, merror
			}
			return []byte(state), nil
		})
		if err != nil {
			return err
		}

		if !token.Valid {
			return fmt.Errorf("Not authorized")
		}
	} else {
		return fmt.Errorf("Not authorized")
	}

	return nil
}

func saveProfileDynamoDB(profile map[string]interface{}) error {

	svc := dynamoDB.CreateDynamoDBClient()

	item, err := dynamoDB.GetSessionData()
	if err != nil {
		return err
	}

	// BodyRequest will be used to take the json response from client and build it
	bodyRequest := models.SessionData{
		ClientID: item.ClientID,
		State:    item.State,
		Profile:  profile,
		Token:    item.Token,
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
