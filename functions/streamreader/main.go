package main

import (
	"context"
	"fmt"
	"math/big"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

//https://play.golang.org/p/955Nm7AYuWu

var div = new(big.Int)
var bigZero = new(big.Int)
var snsSvc *sns.SNS
var topicArn *string

func init() {
	div.SetInt64(5)
	bigZero.SetInt64(0)
	session := session.Must(session.NewSession())
	snsSvc = sns.New(session)
	topicArn = aws.String(os.Getenv("TOPIC"))
}

func processIt(sequenceNo string) bool {
	i := new(big.Int)
	i.SetString(sequenceNo, 10)
	m1 := new(big.Int).Mod(i, div)
	return m1.Cmp(bigZero) != 0
}

func sendIt(data string) error {
	fmt.Printf("send %s\n", data)

	pubIn := sns.PublishInput{
		Message:  aws.String(data),
		TopicArn: topicArn,
	}

	pubOut, err := snsSvc.Publish(&pubIn)
	if err != nil {
		fmt.Println(pubOut)
	}
	return err
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
		if err := sendIt(string(kinesisRecord.SequenceNumber)); err != nil {
			fmt.Println("Error publishing:", err)
		}
	}

	return nil
}

func main() {
	lambda.Start(handler)
}
