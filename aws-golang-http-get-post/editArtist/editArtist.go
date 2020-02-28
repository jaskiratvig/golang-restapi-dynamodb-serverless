package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// BodyRequest is our self-made struct to process JSON request from Client
type BodyRequest struct {
	ArtistID    string   `json:"ArtistID"`
	Name        string   `json:"Name"`
	Songs       []string `json:"Songs"`
	Subcategory string   `json:"Subcategory"`
	Domestic    bool     `json:"Domestic"`
}

// Handler function Using AWS Lambda Proxy Request
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	svc := createDynamoDBClient()

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

	item := BodyRequest{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		message := fmt.Sprintf("Failed to unmarshal Record, %v", err)
		return events.APIGatewayProxyResponse{Body: message, StatusCode: 404}, err
	}

	// BodyRequest will be used to take the json response from client and build it
	bodyRequest := BodyRequest{
		ArtistID:    item.ArtistID,
		Name:        item.Name,
		Songs:       item.Songs,
		Subcategory: item.Subcategory,
		Domestic:    item.Domestic,
	}

	// Unmarshal the json, return 404 if error
	err = json.Unmarshal([]byte(request.Body), &bodyRequest)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 404}, nil
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
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 404}, nil
	}

	//Returning response with AWS Lambda Proxy Response
	return events.APIGatewayProxyResponse{Body: string(response), StatusCode: 200}, nil
}

func createDynamoDBClient() *dynamodb.DynamoDB {
	sess := session.Must(
		session.NewSessionWithOptions(
			session.Options{
				SharedConfigState: session.SharedConfigEnable,
			}))
	return dynamodb.New(sess)
}

func (art BodyRequest) editArtist(artist BodyRequest) BodyRequest {
	art.ArtistID = artist.ArtistID
	art.Name = artist.Name
	art.Songs = artist.Songs
	art.Subcategory = artist.Subcategory
	art.Domestic = artist.Domestic
	return art
}

func main() {
	lambda.Start(Handler)
}
