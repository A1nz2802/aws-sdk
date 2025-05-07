package main

import (
	"aws-sdk/dynamodb"
	"log"
)

// https://www.youtube.com/watch?v=DIQVJqiSUkE
func main() {

	err := dynamodb.Init()

	if err != nil {
		log.Fatalf("Initialization Failed: %v", err)
		return
	}

	orders, err := dynamodb.GetOrders("a1nzdev")

	if err != nil {
		log.Fatalf("failed to get orders: %v", err)
	}

	for _, order := range orders {
		log.Printf("Order: %v - %v\n", order.Username, order.IdOrder)
	}
}
