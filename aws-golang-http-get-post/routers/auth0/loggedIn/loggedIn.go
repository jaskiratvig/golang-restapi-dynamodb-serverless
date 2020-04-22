package main

import (
	"aws-golang-http-get-post/dynamoDB"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Handler function Using AWS Lambda Proxy Request
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	item, err := dynamoDB.GetSessionData()
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 404}, err
	}

	return events.APIGatewayProxyResponse{Body: "Hello " + fmt.Sprintf("%v", item.Profile["name"]), StatusCode: 200}, nil
}

func main() {
	lambda.Start(Handler)
}
