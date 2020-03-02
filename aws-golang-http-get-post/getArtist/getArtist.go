package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/ses"
)

type Artist struct {
	ArtistID    string
	Name        string
	Songs       []string
	Subcategory string
	Domestic    bool
}

const (
	Sender    = "jaskiratvig@gmail.com"
	Recipient = "jaskiratvig@gmail.com"
	Subject   = "Success"
	TextBody  = "This email was sent with Amazon SES using the AWS SDK for Go."
	CharSet   = "UTF-8"
)

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

	item := Artist{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		message := fmt.Sprintf("Failed to unmarshal Record, %v", err)
		return events.APIGatewayProxyResponse{Body: message, StatusCode: 404}, err
	}

	if item.Name == "" {
		message := fmt.Sprintf("Could not find '" + item.Name)
		return events.APIGatewayProxyResponse{Body: message, StatusCode: 404}, err
	}

	//Generate message that want to be sent as body
	message := fmt.Sprintf("Found artist: ArtistID: %+v Name: %+v Subcategory: %+v Songs: %+v Domestic: %+v ", item.ArtistID, item.Name, item.Subcategory, item.Songs, item.Domestic)

	//SES Integration
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)
	svcSes := ses.New(sess)

	HTMLBody := "<h1>Success</h1><p> " + message + "</p>"

	// Assemble the email.
	inputSes := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{
				aws.String(Recipient),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(HTMLBody),
				},
				Text: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(TextBody),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(CharSet),
				Data:    aws.String(Subject),
			},
		},
		Source: aws.String(Sender),
		// Uncomment to use a configuration set
		//ConfigurationSetName: aws.String(ConfigurationSet),
	}

	// Attempt to send the email.
	resultSes, err := svcSes.SendEmail(inputSes)

	// Display error messages if they occur.
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			return events.APIGatewayProxyResponse{Body: aerr.Error(), StatusCode: 404}, nil
		}
	}

	fmt.Println("Email Sent to address: " + Recipient)
	fmt.Println(resultSes)

	//Returning response with AWS Lambda Proxy Response
	return events.APIGatewayProxyResponse{Body: message, StatusCode: 200}, nil
}

func createDynamoDBClient() *dynamodb.DynamoDB {
	sess := session.Must(
		session.NewSessionWithOptions(
			session.Options{
				SharedConfigState: session.SharedConfigEnable,
			}))
	return dynamodb.New(sess)
}

func main() {
	lambda.Start(Handler)
}
