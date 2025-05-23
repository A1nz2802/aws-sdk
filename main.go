package main

import (
	"aws-sdk/dynamodb"
	"log"
)

// https://www.youtube.com/watch?v=DIQVJqiSUkE
func main() {

	/* err := dynamodb.Init()

	if err != nil {
		log.Fatalf("Initialization Failed: %v", err)
		return
	} */

	err := dynamodb.CreateOrderStatusDateGSI()

	if err != nil {
		log.Fatalf("failed to get price: %v", err)
	}

}
