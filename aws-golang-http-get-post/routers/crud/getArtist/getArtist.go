package main

import (
	"aws-golang-http-get-post/dynamoDB"
	"aws-golang-http-get-post/ses"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Handler function Using AWS Lambda Proxy Request
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	item, err := dynamoDB.GetArtist(request)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 404}, err
	}

	if item.Name == "" {
		message := fmt.Sprintf("Could not find '" + item.Name)
		return events.APIGatewayProxyResponse{Body: message, StatusCode: 404}, err
	}

	message := fmt.Sprintf("Found artist: ArtistID: %+v Name: %+v Subcategory: %+v Songs: %+v Domestic: %+v ", item.ArtistID, item.Name, item.Subcategory, item.Songs, item.Domestic)
	HTMLBody := "<h1>Success</h1><p> " + message + "</p>"

	err = ses.SendEmail(HTMLBody)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 404}, err
	}

	return events.APIGatewayProxyResponse{Body: message, StatusCode: 200}, nil
}

func main() {
	lambda.Start(Handler)
}
