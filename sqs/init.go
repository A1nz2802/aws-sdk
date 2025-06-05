package sqs

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/joho/godotenv"
)

var ctx = context.Background()
var client = GetClient()

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

// If is fifo queue, queue name should have a .fifo suffix.
func CreateQueue(isFifoQueue bool, queueName string) (string, error) {

	attrsMap := make(map[string]string)

	attrsMap["ReceiveMessageWaitTimeSeconds"] = "20"

	if isFifoQueue {
		attrsMap["FifoQueue"] = "true"
	}

	queue, err := client.CreateQueue(ctx, &sqs.CreateQueueInput{
		QueueName:  &queueName,
		Attributes: attrsMap,
	})

	if err != nil {
		return "", err
	}

	return *queue.QueueUrl, err
}

// CreateExampleQueues creates three example SQS queues:
//
// 1. A standard queue.
//
// 2. A FIFO queue.
//
// 3. A FIFO queue with a Dead Letter Queue (DLQ) configured.
func CreateExampleQueues() error {
	const (
		firstQueueName  = "my-std-queue"
		secondQueueName = "my-fifo-queue.fifo"
		dlqQueueName    = "my-dlq-queue"
	)

	fmt.Println("Attempting to create example SQS queues...")

	// Create a standard queue
	firstQueueURL, err := CreateQueue(false, firstQueueName)
	if err != nil {
		return fmt.Errorf("❌ failed to create standard queue '%s': %w", firstQueueName, err)
	}

	fmt.Printf("✅ Standard queue '%s' created successfully with URL: %s\n", firstQueueName, firstQueueURL)

	// Create a FIFO queue
	secondQueueURL, err := CreateQueue(true, secondQueueName)
	if err != nil {
		return fmt.Errorf("❌ failed to create FIFO queue '%s': %w", secondQueueName, err)
	}

	fmt.Printf("✅ FIFO queue '%s' created successfully with URL: %s\n", secondQueueName, secondQueueURL)

	// Create a queue with an attached DLQ
	err = AttachDeadLetterQueue(secondQueueName, dlqQueueName, false)
	if err != nil {
		return fmt.Errorf("failed to attach DLQ '%s' to queue '%s': %w", dlqQueueName, secondQueueName, err)
	}

	fmt.Printf("✅ DLQ '%s' successfully configured for queue '%s'.\n", dlqQueueName, secondQueueName)

	fmt.Println("\nAll example queues created and configured successfully!")

	return nil
}

// Creates a Dead Letter Queue (DLQ) and configures an existing main queue.
func AttachDeadLetterQueue(mainQueueName string, dlqQueueName string, isFifoQueue bool) error {

	fmt.Printf("Attempting to attach DLQ '%s' to main queue '%s'...\n", dlqQueueName, mainQueueName)

	// Create DLQ
	dlqQueueURL, err := CreateQueue(isFifoQueue, dlqQueueName)
	if err != nil {
		return fmt.Errorf("❌ failed to create DLQ '%s': %w", dlqQueueName, err)
	}

	fmt.Printf("✅ DLQ '%s' created successfully with URL: %s\n", dlqQueueName, dlqQueueURL)

	// Get the main queue URL and ARN
	mainQueueUrl, err := GetQueueURL(mainQueueName)
	if err != nil {
		return fmt.Errorf("❌ failed to get URL for main queue '%s': %w", mainQueueName, err)
	}

	queueArn, err := GetQueueArn(mainQueueUrl)
	if err != nil {
		return fmt.Errorf("❌ failed to get ARN for main queue '%s': %w", mainQueueName, err)
	}

	redrivePolicyMap := map[string]interface{}{
		"deadLetterTargetArn": queueArn,
		"maxReceiveCount":     10, // maybe this can configurable
	}

	redrivePolicyJSON, err := json.Marshal(redrivePolicyMap)
	if err != nil {
		return fmt.Errorf("❌ failed to marshal RedrivePolicy for DLQ '%s': %w", dlqQueueName, err)

	}

	// Set the RedrivePolicy attribute on the main queue
	_, err = client.SetQueueAttributes(ctx, &sqs.SetQueueAttributesInput{
		QueueUrl: &mainQueueUrl,
		Attributes: map[string]string{
			"RedrivePolicy": string(redrivePolicyJSON),
		},
	})

	if err != nil {
		return fmt.Errorf("❌ failed to set RedrivePolicy for main queue '%s': %w", mainQueueName, err)
	}

	fmt.Printf("✅ DLQ '%s' successfully attached to main queue '%s'.\n", dlqQueueName, mainQueueName)

	return nil
}
