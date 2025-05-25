package main

import (
	"aws-sdk/sqs"
	"fmt"
	"log"
)

func main() {

	/* _, err := sqs.CreateQueue(false)

	if err != nil {
		log.Fatalf("Initialization Failed: %v", err)
		return
	} */

	sqs.ListQueues()

	err := sqs.SendMessage("https://sqs.us-east-1.amazonaws.com/266735829330/MyQueue", "WEIRD TEXT HERE :D")

	if err != nil {
		log.Printf("Error here :D: %v", err)
		return
	}

	messages, err := sqs.ReceiveMessage("https://sqs.us-east-1.amazonaws.com/266735829330/MyQueue", 10, 1)

	if err != nil {
		log.Printf("Error here :D: %v", err)
		return
	}

	for _, message := range messages {
		fmt.Println(message)
	}
}
