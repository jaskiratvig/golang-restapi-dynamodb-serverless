package main

import (
	"aws-golang-http-get-post/auth"
	"aws-golang-http-get-post/dynamoDB"
	"aws-golang-http-get-post/models"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/ssm"
	jwt "github.com/dgrijalva/jwt-go"
)

// Handler function Using AWS Lambda Proxy Request
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	state, err := generateRandomState()
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 400}, err
	}

	validToken, err := generateJWT(state)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 400}, err
	}

	errorMessage := saveStateTokenDynamoDB(state, validToken)
	if errorMessage != "" {
		err = fmt.Errorf(errorMessage)
		return events.APIGatewayProxyResponse{Body: errorMessage, StatusCode: 400}, err
	}

	authenticator, err := auth.NewAuthenticator()
	if err != nil {
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

func saveStateTokenDynamoDB(state string, token string) string {

	sess := session.New()
	svcSSM := ssm.New(sess)

	clientID, err := svcSSM.GetParameter(
		&ssm.GetParameterInput{
			Name: aws.String("/dev/ClientID"),
		},
	)

	session := models.SessionData{
		ClientID: aws.StringValue(clientID.Parameter.Value),
		State:    state,
		Profile:  nil,
		Token:    token,
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

func generateJWT(state string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["exp"] = time.Now().Add(time.Minute).Unix()

	tokenString, err := token.SignedString([]byte(state))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func main() {
	lambda.Start(Handler)
}
