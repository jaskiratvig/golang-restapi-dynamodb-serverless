package ses

import (
	"aws-golang-http-get-post/constants"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/aws/aws-sdk-go/service/ssm"
)

//SendEmail is the function that sends an email using the HTMLBody input to the recepient defined in AWS Parameter Store
func SendEmail(HTMLBody string) error {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)
	if err != nil {
		return err
	}

	svc := ssm.New(sess)

	recipient, err := svc.GetParameter(
		&ssm.GetParameterInput{
			Name: aws.String("/dev/Recipient"),
		},
	)

	inputSes := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{recipient.Parameter.Value},
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
	}

	svcSes := ses.New(sess)
	_, err = svcSes.SendEmail(inputSes)
	if err != nil {
		return err
	}

	return nil
}
