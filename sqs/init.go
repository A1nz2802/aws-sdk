package sqs

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/joho/godotenv"
)

/*
SQS TYPES (FIFO AND STANDARD)
SQS DEAD LETTER QUEUE
SQS DELAY QUEUE
SQS VISIBILITY TIMEOUT
LONG POLLING VS SHORT POLLING

*/

var ctx = context.Background()
var client = GetClient()
var queueName = "MySecondQueue"

func GetClient() *sqs.Client {
	err := godotenv.Load()

	if err != nil {
		// or panic
		log.Fatalf("Error loading .env file, %v", err)
	}

	sdkConfig, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		log.Fatalf("Couldn't load default configuration: %v", err)
	}

	return sqs.NewFromConfig(sdkConfig)
}

func CreateQueue(isFifoQueue bool) (string, error) {
	var queueUrl string
	queueAttributes := map[string]string{}

	if isFifoQueue {
		queueAttributes["FifoQueue"] = "true"
	}

	queue, err := client.CreateQueue(ctx, &sqs.CreateQueueInput{
		QueueName:  &queueName,
		Attributes: queueAttributes,
	})

	if err != nil {
		log.Printf("Couldn't create queue %v. Here's why: %v\n", queueName, err)
		return "", err
	}

	queueUrl = *queue.QueueUrl

	return queueUrl, err
}
