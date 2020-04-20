package ses

import (
	"aws-golang-http-get-post/constants"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/aws/aws-sdk-go/service/ssm"
)

//SendEmail is the function that sends an email using the HTMLBody input to the recepient defined in AWS Parameter Store
func SendEmail(HTMLBody string, response string) (events.APIGatewayProxyResponse, error) {
	//SES Integration
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)
	svcSes := ses.New(sess)
	svc := ssm.New(sess)

	recipient, err := svc.GetParameter(
		&ssm.GetParameterInput{
			Name: aws.String("/dev/Recipient"),
		},
	)

	// Assemble the email.
	inputSes := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{
				aws.String(aws.StringValue(recipient.Parameter.Value)),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(constants.CharSet),
					Data:    aws.String(HTMLBody),
				},
				Text: &ses.Content{
					Charset: aws.String(constants.CharSet),
					Data:    aws.String(constants.TextBody),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(constants.CharSet),
				Data:    aws.String(constants.Subject),
			},
		},
		Source: aws.String(constants.Sender),
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

	fmt.Println("Email Sent to address: " + aws.StringValue(recipient.Parameter.Value))
	fmt.Println(resultSes)

	//Returning response with AWS Lambda Proxy Response
	return events.APIGatewayProxyResponse{Body: string(response), StatusCode: 200}, nil
}
