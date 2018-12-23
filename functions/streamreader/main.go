package main

import (
	"context"
	"fmt"
	"math/big"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

//https://play.golang.org/p/955Nm7AYuWu

var div = new(big.Int)
var bigZero = new(big.Int)

func init() {
	div.SetInt64(5)
	bigZero.SetInt64(0)
}

func processIt(sequenceNo string) bool {
	i := new(big.Int)
	i.SetString(sequenceNo, 10)
	m1 := new(big.Int).Mod(i, div)
	return m1.Cmp(bigZero) != 0
}

func handler(ctx context.Context, kinesisEvent events.KinesisEvent) error {
	for _, record := range kinesisEvent.Records {
		kinesisRecord := record.Kinesis
		//dataBytes := kinesisRecord.Data
		//dataText := string(dataBytes)

		//fmt.Printf("%s Data = %s \n", record.EventName, dataText)
		if processIt(kinesisRecord.SequenceNumber) == false {
			fmt.Printf("Skip processing of %+v\n", kinesisRecord)
			return nil
		}

		fmt.Printf("Process record %+v\n", kinesisRecord)
	}

	return nil
}

func main() {
	lambda.Start(handler)
}
