package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bradfitz/gomemcache/memcache"
)

var memcachedEndpoint = fmt.Sprintf("%s:%s", os.Getenv("MEMCACHED_ENDPOINT"), os.Getenv("MEMCACHED_PORT"))

var client = memcache.New(memcachedEndpoint)

type ProcessedStatus struct {
	Unprocessed []string
}

func getProcessingStatus(seqNos []string) ([]byte, error) {
	var unprocessed ProcessedStatus

	for _, sn := range seqNos {
		fmt.Println("Checking", sn)
		_, err := client.Get(sn)
		if err == nil {
			continue
		}

		switch err {
		case memcache.ErrCacheMiss:
			unprocessed.Unprocessed = append(unprocessed.Unprocessed, sn)
			break
		default:
			return nil, err
		}
	}

	return json.Marshal(&unprocessed)
}

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Printf("Processing request data for request %s.\n", request.RequestContext.RequestID)
	fmt.Printf("Body: %-v\n", request.Body)

	fmt.Println("Headers:")
	for key, value := range request.Headers {
		fmt.Printf("    %s: %s\n", key, value)
	}

	fmt.Println("Unmarshall request")
	var seqNos []string
	err := json.Unmarshal([]byte(request.Body), &seqNos)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 400}, nil
	}

	fmt.Println("Get processing status")
	response, err := getProcessingStatus(seqNos)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 400}, nil
	}

	return events.APIGatewayProxyResponse{Body: string(response), StatusCode: 200}, nil
}

func main() {
	lambda.Start(handleRequest)
}
