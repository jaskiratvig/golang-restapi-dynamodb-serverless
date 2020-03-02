package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
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

	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"Name": {
				S: aws.String(request.PathParameters["Name"]),
			},
		},
		TableName: aws.String("Artists"),
	}

	_, err := svc.DeleteItem(input)
	if err != nil {
		fmt.Println("Got error calling DeleteItem:")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	//Generate message that want to be sent as body
	message := fmt.Sprintf("Deleted artist: Name: %+v ", request.PathParameters["Name"])

	//SES Integration
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)
	svcSes := ses.New(sess)

	HTMLBody := "<h1>Success</h1><p>Artist " + request.PathParameters["Name"] + " has been deleted from the database.</p>"

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
	result, err := svcSes.SendEmail(inputSes)

	// Display error messages if they occur.
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			return events.APIGatewayProxyResponse{Body: aerr.Error(), StatusCode: 404}, nil
		}
	}

	fmt.Println("Email Sent to address: " + Recipient)
	fmt.Println(result)

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
