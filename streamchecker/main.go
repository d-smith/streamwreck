package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	streamName = kingpin.Flag("stream", "Stream to check").Required().String()
)

func main() {
	kingpin.Parse()
	fmt.Println(*streamName)

	session := session.Must(session.NewSession())
	kinesisSvc := kinesis.New(session)
	fmt.Printf("%-v\n", kinesisSvc)

	//GetRecords - needs a shard iterator
	//Use GetShardIterator for the initial shard
	//  -- Needs a shard id to get the iterator
	//Use ListShards to get the shard names for the stream
}
