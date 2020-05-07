package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"os"
)

var (
	svc       *dynamodb.DynamoDB
	tableName *string
)

func init() {
	mySession := session.Must(session.NewSession())
	svc = dynamodb.New(mySession)
	tableName = aws.String(os.Getenv("TABLE_NAME"))
}

func handler(request events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	input := &dynamodb.DeleteItemInput{
		TableName: tableName,
		Key: map[string]*dynamodb.AttributeValue{
			"connectionId": {
				S: aws.String(request.RequestContext.ConnectionID),
			},
		},
	}

	_, err := svc.DeleteItem(input)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf("Failed to disconnect: %v", err),
		}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "Disconnected",
	}, nil
}

func main() {
	lambda.Start(handler)
}
