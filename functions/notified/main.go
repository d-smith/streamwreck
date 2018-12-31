package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bradfitz/gomemcache/memcache"
)

var memcachedEndpoint = fmt.Sprintf("%s:%s", os.Getenv("MEMCACHED_ENDPOINT"), os.Getenv("MEMCACHED_PORT"))

var client = memcache.New(memcachedEndpoint)

//TODO - error handling at the lambda level - return error and retry? Log and continue?

func handler(ctx context.Context, snsEvent events.SNSEvent) {
	fmt.Println(memcachedEndpoint)
	for _, record := range snsEvent.Records {
		snsRecord := record.SNS

		fmt.Printf("[%s %s] Message = %s \n", record.EventSource, snsRecord.Timestamp, snsRecord.Message)
		itemToCache := memcache.Item{
			Key:        snsRecord.Message,
			Value:      []byte("processed"),
			Expiration: 36 * 60 * 60, //36 hours in seconds
		}
		if err := client.Set(&itemToCache); err != nil {
			fmt.Printf("Error setting key %s: %s\n", snsRecord.Message, err.Error())
		}
	}
}

func main() {
	lambda.Start(handler)
}
