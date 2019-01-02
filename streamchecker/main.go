package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	streamName = kingpin.Flag("stream", "Stream to check").Required().String()
)

func getShardIds(kinesisSvc *kinesis.Kinesis, streamName *string) ([]string, error) {
	response, err := kinesisSvc.ListShards(&kinesis.ListShardsInput{
		StreamName: streamName,
	})

	if err != nil {
		return nil, err
	}

	//TODO - implement handling of logic if more shards are present than can be returned
	// in a single call to list shards
	if response.NextToken != nil {
		fmt.Println("WARNING: MORE SHARDS EXISTS THAN A SINGLE CALL TO ListShards CAN RETURN")
	}

	var shardIds []string
	for _, shard := range response.Shards {
		shardIds = append(shardIds, *shard.ShardId)
	}

	return shardIds, err
}

func main() {
	kingpin.Parse()
	fmt.Println(*streamName)

	session := session.Must(session.NewSession())
	kinesisSvc := kinesis.New(session)

	//GetRecords - needs a shard iterator
	//Use GetShardIterator for the initial shard
	//  -- Needs a shard id to get the iterator
	//Use ListShards to get the shard names for the stream

	out, err := getShardIds(kinesisSvc, streamName)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Printf("%-v\n", out)
}
