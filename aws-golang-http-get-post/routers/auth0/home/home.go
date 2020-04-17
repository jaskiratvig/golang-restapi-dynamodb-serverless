package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Handler function Using AWS Lambda Proxy Request
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var HTML = `
		<html>
			<body>
				<a href="/dev/login"> 
					Login 
				</a>
			</body>
		</html>
	`

	return events.APIGatewayProxyResponse{Body: HTML, Headers: map[string]string{
		"Content-Type": "text/html",
	}, StatusCode: 200}, nil
}

func main() {
	lambda.Start(Handler)
}
