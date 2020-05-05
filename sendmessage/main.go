package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"os"
)

var (
	mySession = session.Must(session.NewSession())
	svc       = dynamodb.New(mySession)
)

type message struct {
	// The action to take
	Message string `json:"message"`
	// The contents of the message
	Data string `json:"data"`
}

func handler(event events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Input param for a table scan
	scanInput := &dynamodb.ScanInput{
		TableName:            aws.String(os.Getenv("TABLE_NAME")),
		ProjectionExpression: aws.String("connectionId"),
	}

	// Scan the table for all the connection ids
	connectionData, err := svc.Scan(scanInput)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Error sending message",
		}, err
	}

	apiMgmtService := apigatewaymanagementapi.New(mySession, &aws.Config{
		Endpoint: aws.String(fmt.Sprintf("%s/%s", event.RequestContext.DomainName, event.RequestContext.Stage)),
	})

	var receivedMessage message
	err = json.Unmarshal([]byte(event.Body), &receivedMessage)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Error unmarshalling JSON body",
		}, err
	}

	for _, item := range connectionData.Items {
		connectionId := item["connectionId"].S
		postToConnectionRequest := &apigatewaymanagementapi.PostToConnectionInput{
			ConnectionId: connectionId,
			Data:         []byte(receivedMessage.Data),
		}
		_, err := apiMgmtService.PostToConnection(postToConnectionRequest)
		if err != nil {
			fmt.Println("Found stale connection, deleting", connectionId)
			deleteItem(connectionId)
		}
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "Data sent",
	}, nil
}

func deleteItem(connectionId *string) {
	deleteRequest := &dynamodb.DeleteItemInput{
		TableName: aws.String(os.Getenv("TABLE_NAME")),
		Key: map[string]*dynamodb.AttributeValue{
			"connectionId": {
				S: connectionId,
			},
		},
	}
	_, err := svc.DeleteItem(deleteRequest)
	if err != nil {
		fmt.Println("Error deleting connection:", err)
	}
}

func main() {
	lambda.Start(handler)
}
