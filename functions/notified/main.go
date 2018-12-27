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

func handler(ctx context.Context, snsEvent events.SNSEvent) {
	fmt.Println(memcachedEndpoint)
	for _, record := range snsEvent.Records {
		snsRecord := record.SNS

		fmt.Printf("[%s %s] Message = %s \n", record.EventSource, snsRecord.Timestamp, snsRecord.Message)
		if err := client.Set(&memcache.Item{Key: snsRecord.Message, Value: []byte("processed")}); err != nil {
			fmt.Printf("Error setting key %s: %s\n", snsRecord.Message, err.Error())
		}

		fmt.Println("For grins, can we get what we just set?")
		item, err := client.Get(snsRecord.Message)
		switch err {
		case nil:
			fmt.Printf("Got item %-v", *item)
		default:
			fmt.Printf("Error on get: %s", err.Error())
		}

	}
}

func main() {
	lambda.Start(handler)
}
