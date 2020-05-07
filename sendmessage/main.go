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
	"github.com/ksanta/chat/model"
	"os"
	"sync"
)

var (
	svc            *dynamodb.DynamoDB
	apiMgmtService *apigatewaymanagementapi.ApiGatewayManagementApi
	tableName      *string
)

func init() {
	mySession := session.Must(session.NewSession())
	svc = dynamodb.New(mySession)
	tableName = aws.String(os.Getenv("TABLE_NAME"))
}

func handler(event events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	// One time setup using data from a request
	if apiMgmtService == nil {
		mySession := session.Must(session.NewSession())
		apiMgmtService = apigatewaymanagementapi.New(mySession, &aws.Config{
			Endpoint: aws.String(fmt.Sprintf("%s/%s", event.RequestContext.DomainName, event.RequestContext.Stage)),
		})
	}

	// Input param for a table scan
	scanInput := &dynamodb.ScanInput{
		TableName:            tableName,
		ProjectionExpression: aws.String("connectionId"),
	}

	// Scan the table for all the connection ids
	connectionsResult, err := svc.Scan(scanInput)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Error scanning for connections",
		}, err
	}

	var receivedMessage model.Message
	err = json.Unmarshal([]byte(event.Body), &receivedMessage)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Error unmarshalling JSON body",
		}, err
	}

	waitGroup := sync.WaitGroup{}
	waitGroup.Add(len(connectionsResult.Items))

	for _, item := range connectionsResult.Items {
		// Send messages to each connected person concurrently
		itemCopy := item
		go func() {
			connectionId := itemCopy["connectionId"].S
			postToConnectionRequest := &apigatewaymanagementapi.PostToConnectionInput{
				ConnectionId: connectionId,
				Data:         []byte(receivedMessage.Data),
			}
			_, err := apiMgmtService.PostToConnection(postToConnectionRequest)
			if err != nil {
				fmt.Println("Found stale connection, deleting", connectionId)
				deleteItem(connectionId)
			}
			waitGroup.Done()
		}()
	}
	waitGroup.Wait()

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
