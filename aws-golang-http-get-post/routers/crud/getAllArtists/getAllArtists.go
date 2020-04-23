package main

import (
	"aws-golang-http-get-post/dynamoDB"
	"aws-golang-http-get-post/models"
	"aws-golang-http-get-post/ses"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

// Handler function Using AWS Lambda Proxy Request
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	svc := dynamoDB.CreateDynamoDBClient()

	proj := expression.NamesList(expression.Name("Name"), expression.Name("Subcategory"), expression.Name("Songs"), expression.Name("Domestic"))
	expr, err := expression.NewBuilder().WithProjection(proj).Build()
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 404}, err
	}

	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String("Artists"),
	}
	result, err := svc.Scan(params)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 404}, err
	}

	message := fmt.Sprintf("")

	for _, i := range result.Items {
		item := models.Artist{}

		err = dynamodbattribute.UnmarshalMap(i, &item)
		if err != nil {
			return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 404}, err
		}

		message = message + fmt.Sprintf("Name: %+v Subcategory: %+v Songs: %+v Domestic: %+v ", item.Name, item.Subcategory, item.Songs, item.Domestic)
	}

	HTMLBody := "<h1>Success</h1><p>Here is a list of all artists in the database: " + message + "</p>"

	err = ses.SendEmail(HTMLBody)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 404}, err
	}

	return events.APIGatewayProxyResponse{Body: message, StatusCode: 200}, nil
}

func main() {
	lambda.Start(Handler)
}
