package sqs

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

func SendMessage(isFifoQueue bool, queueURL string, messageBody string) error {
	input := &sqs.SendMessageInput{
		MessageBody: aws.String(messageBody),
		QueueUrl:    aws.String(queueURL),
	}

	if isFifoQueue {
		input.MessageGroupId = aws.String("my-message-group-id")
		input.MessageDeduplicationId = aws.String("my-deduplication-id")
	}

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

func DeleteMessage(queueURL string, msg types.Message) error {
	fmt.Printf("🗑️  Attempting to delete message ID %s from queue...\n", aws.ToString(msg.MessageId))
	deleteInput := &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(queueURL),
		ReceiptHandle: msg.ReceiptHandle,
	}

	_, err := client.DeleteMessage(ctx, deleteInput)

	if err != nil {
		return err
	}

	fmt.Printf("✅ Message ID %s deleted successfully.\n", aws.ToString(msg.MessageId))

	return nil
}

func SimulateCommunication(queueName string, isFifoQueue bool) error {

	fmt.Printf("🔍 Getting URL for queue '%s'...\n", queueName)
	queueURL, err := GetQueueURL(queueName)
	if err != nil {
		return fmt.Errorf("❌ failed to get URL for queue %s: %w", queueName, err)
	}

	fmt.Printf("✅ Queue URL obtained: %s\n\n", queueURL)

	messages := []string{"First Message :)", "Second Message :(", "Test Message :')"}

	for i, message := range messages {
		err := SendMessage(isFifoQueue, queueURL, message)

		if err != nil {
			return fmt.Errorf("failed to send message %d: %w", i, err)
		}
	}

	time.Sleep(2 * time.Second)

	maxMessagesToReceive := int32(5)
	waitTimeSeconds := int32(20) // Use long polling

	receivedMessages, err := ReceiveMessage(queueURL, maxMessagesToReceive, waitTimeSeconds)

	if err != nil {
		return fmt.Errorf("failed to receive messages: %w", err)
	}

	if len(receivedMessages) == 0 {
		fmt.Println("📪 No messages received from the queue this time.")
		return nil
	}

	fmt.Printf("📫 Received %d messages:\n", len(receivedMessages))
	for i, msg := range receivedMessages {
		fmt.Printf("--- Message %d ---\n", i+1)
		fmt.Printf("🆔 Message ID: %s\n", aws.ToString(msg.MessageId))

		if msg.ReceiptHandle == nil {
			fmt.Printf("⚠️ message with ID %s has no ReceiptHandle. It cannot be deleted.", aws.ToString(msg.MessageId))
			continue
		}

		fmt.Printf("  📄 Body: %s\n", aws.ToString(msg.Body))

		// --- 3. Delete Received Message ---
		err := DeleteMessage(queueURL, msg)

		if err != nil {
			return fmt.Errorf("❌ error deleting message with ID %s: %w", aws.ToString(msg.MessageId), err)
		}
	}

	fmt.Printf("\n-------------------------------------------------------------\n")
	fmt.Printf("🏁 SQS communication simulation finished. 👋\n")

	return nil
}

func TestEmptyQueue(queueName string) error {
	queueURL, err := GetQueueURL(queueName)
	if err != nil {
		return fmt.Errorf("❌ failed to get URL for queue %s: %w", queueName, err)
	}

	maxMessagesToReceive := int32(5)
	waitTimeSeconds := int32(20) // Use long polling

	_, err = ReceiveMessage(queueURL, maxMessagesToReceive, waitTimeSeconds)
	if err != nil {
		return fmt.Errorf("❌ failed to received messages %s: %w", queueName, err)
	}

	fmt.Println("No messages received")

	return nil
}
