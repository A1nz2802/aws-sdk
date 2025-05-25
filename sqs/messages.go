package sqs

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

func SendMessage(queueURL string, messageBody string) error {
	input := &sqs.SendMessageInput{
		MessageBody: aws.String(messageBody),
		QueueUrl:    aws.String(queueURL),
	}

	// Para colas FIFO, necesitarías MessageGroupId y opcionalmente MessageDeduplicationId
	// if isFifoQueue {
	//    input.MessageGroupId = aws.String("my-message-group-id")
	//    input.MessageDeduplicationId = aws.String("my-deduplication-id") // O dejar que ContentBasedDeduplication lo genere
	// }

	_, err := client.SendMessage(ctx, input)

	if err != nil {
		return fmt.Errorf("error sending messages: %w", err)
	}

	fmt.Printf("message send successfully to queue: %s\n", messageBody)

	return nil
}

func ReceiveMessage(queueURL string, maxNumberOfMessages int32, waitTimeSeconds int32) ([]types.Message, error) {
	input := &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(queueURL),
		MaxNumberOfMessages: maxNumberOfMessages,
		WaitTimeSeconds:     waitTimeSeconds,
		// También puedes solicitar atributos de mensaje como SenderId, SentTimestamp, etc.
		// MessageAttributeNames: []string{string(types.MessageAttributeNameAll)},
	}

	result, err := client.ReceiveMessage(ctx, input)

	if err != nil {
		return nil, fmt.Errorf("error ocurred when received messages: %w", err)
	}

	if len(result.Messages) == 0 {
		fmt.Println("no messages were received.")
	}

	return result.Messages, nil
}

func Communication(queueName string) {

	fmt.Printf("Getting URL queue '%s'...\n", queueName)
	queueURL, err := getQueueURL(ctx, sqsClient, queueName)

	if err != nil {
		log.Printf("Error getting queueURL: %v", err)
	}
}
