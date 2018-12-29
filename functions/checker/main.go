package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type ProcessedStatus struct {
	Unprocessed []string
}

func getProcessingStatus(seqNos []string) ([]byte, error) {
	var unprocessed ProcessedStatus
	unprocessed.Unprocessed = append(unprocessed.Unprocessed, "2")
	return json.Marshal(&unprocessed)
}

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Printf("Processing request data for request %s.\n", request.RequestContext.RequestID)
	fmt.Printf("Body: %-v\n", request.Body)

	fmt.Println("Headers:")
	for key, value := range request.Headers {
		fmt.Printf("    %s: %s\n", key, value)
	}

	var seqNos []string
	err := json.Unmarshal([]byte(request.Body), &seqNos)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 400}, nil
	}

	response, err := getProcessingStatus(seqNos)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 400}, nil
	}

	return events.APIGatewayProxyResponse{Body: string(response), StatusCode: 200}, nil
}

func main() {
	lambda.Start(handleRequest)
}
