package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	streamName        = kingpin.Flag("stream", "Stream to check").Required().String()
	reconcileEndpoint = kingpin.Flag("endpoint", "Reconcile API endpoint").Required().String()
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

func processRecords(records []*kinesis.Record, endpoint *string) {
	var seqNos []string

	//Extract sequence numbers to check
	for _, rec := range records {
		seqNos = append(seqNos, *rec.SequenceNumber)
	}

	fmt.Println("checking", seqNos)

	//Form payload
	payloadBytes, err := json.Marshal(&seqNos)
	if err != nil {
		fmt.Printf("WARNING: errors marshaling payload: %s", err.Error())
		return
	}

	//Invoke endpoint
	resp, err := http.Post(*endpoint, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		fmt.Println("Error invoking reconcile API", err.Error())
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading body", err.Error())
		return
	}

	fmt.Println(string(body))
}

func readStreamRecords(kinesisSvc *kinesis.Kinesis, shardID, streamName string, wg *sync.WaitGroup) {
	fmt.Println("Processing records for shard id", shardID)

	input := &kinesis.GetShardIteratorInput{
		ShardId:           aws.String(shardID),
		ShardIteratorType: aws.String("TRIM_HORIZON"),
		StreamName:        aws.String(streamName),
	}

	var shardIterator *string
	recordSetSize := int64(20) //Process 20 records at a time - current value is arbitrary

	output, err := kinesisSvc.GetShardIterator(input)
	if err != nil {
		fmt.Printf("Error obtaining initial stream iterator: %s", err.Error())
		goto End
	}

	shardIterator = output.ShardIterator
	fmt.Println(shardIterator)

	for {
		getRecordsOutput, err := kinesisSvc.GetRecords(&kinesis.GetRecordsInput{
			ShardIterator: shardIterator,
			Limit:         &recordSetSize,
		})
		if err != nil {
			fmt.Println("Error reading records:", err.Error())
			break
		}

		switch len(getRecordsOutput.Records) {
		case 0:
			fmt.Println("No records to process...")
			break
		default:
			processRecords(getRecordsOutput.Records, reconcileEndpoint)
		}

		shardIterator = getRecordsOutput.NextShardIterator
		if shardIterator == nil {
			fmt.Println("Nil shard iterator returned")
			break
		}

		//We don't want to call GetRecords more than  5 times/sec - we'll call it even less
		//frequently here to reduce the output verbosity in dev mode
		time.Sleep(6 * time.Second)

	}

End:
	wg.Done()
}

func main() {
	kingpin.Parse()
	fmt.Println(*streamName)

	session := session.Must(session.NewSession())
	kinesisSvc := kinesis.New(session)

	out, err := getShardIds(kinesisSvc, streamName)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Printf("%-v\n", out)

	var wg sync.WaitGroup
	wg.Add(len(out))

	for _, id := range out {
		go readStreamRecords(kinesisSvc, id, *streamName, &wg)
	}

	wg.Wait()

}
