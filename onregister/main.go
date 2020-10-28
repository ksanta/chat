package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/ksanta/chat/model"
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
	// Extract player response from the request
	registerMessage := model.Message{}
	err := json.Unmarshal([]byte(request.Body), &registerMessage)
	if err != nil {
		return newErrorResponse("error unmarshalling JSON body", err)
	}

	user := &model.User{
		ConnectionId: request.RequestContext.ConnectionID,
		Room:         registerMessage.RoomName,
	}
	fmt.Printf("Registering user %+v\n", *user)

	itemMap, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		return newErrorResponse("Failed to marshal storage item map", err)
	}

	input := &dynamodb.PutItemInput{
		TableName: tableName,
		Item:      itemMap,
	}

	_, err = svc.PutItem(input)
	if err != nil {
		return newErrorResponse("Failed to put item into storage", err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "Connected",
	}, nil
}

func newErrorResponse(msg string, err error) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: 500,
		Body:       fmt.Sprintf("%s: %v", msg, err),
	}, err
}

func main() {
	lambda.Start(handler)
}
