package dynamoDB

import (
	"aws-golang-http-get-post/models"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/ssm"
)

// CreateDynamoDBClient sets up a dynamoDB client
func CreateDynamoDBClient() *dynamodb.DynamoDB {
	sess := session.Must(
		session.NewSessionWithOptions(
			session.Options{
				SharedConfigState: session.SharedConfigEnable,
			}))
	return dynamodb.New(sess)
}

// GetSessionData retrieves session data from the dynamoDB table
func GetSessionData() (models.SessionData, error) {
	svc := CreateDynamoDBClient()

	sess := session.New()
	svcSSM := ssm.New(sess)

	clientID, err := svcSSM.GetParameter(
		&ssm.GetParameterInput{
			Name: aws.String("/dev/ClientID"),
		},
	)

	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("SessionData"),
		Key: map[string]*dynamodb.AttributeValue{
			"ClientID": {
				S: aws.String(aws.StringValue(clientID.Parameter.Value)),
			},
		},
	})
	if err != nil {
		return models.SessionData{}, err
	}

	item := models.SessionData{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		return models.SessionData{}, err
	}

	return item, nil
}

// GetArtist retrieves an artist from the dynamoDB table
func GetArtist(request events.APIGatewayProxyRequest) (models.Artist, error) {
	svc := CreateDynamoDBClient()

	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("Artists"),
		Key: map[string]*dynamodb.AttributeValue{
			"Name": {
				S: aws.String(request.PathParameters["Name"]),
			},
		},
	})
	if err != nil {
		return models.Artist{}, err
	}

	item := models.Artist{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		return models.Artist{}, err
	}

	return item, nil
}
