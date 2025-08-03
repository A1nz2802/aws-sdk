package main

import (
	"aws-sdk/lambda"
	"log"
)

func main() {

	/* err := sqs.CreateExampleQueues()
	if err != nil {
		log.Printf("failed to create example queues: %v", err)
		return
	}

	err = sqs.AttachDeadLetterQueue("my-std-queue", "my-dlq-queue", false)
	if err != nil {
		log.Printf("failed to attach dlq to queue: %v", err)
		return
	} */

	/* err := sqs.TestEmptyQueue("my-std-queue")
	if err != nil {
		log.Printf("failed to test empety queue: %v", err)
		return
	} */

	/* err := sqs.SimulateCommunication("my-std-queue", false)
	if err != nil {
		log.Printf("failed to simulate communication: %v", err)
		return
	*/

	err := lambda.CreateExample1()

	if err != nil {
		log.Println("❌ error creating lambda example 01")
	}

}
