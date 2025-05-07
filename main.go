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

	user, err := dynamodb.GetUserByPk("a1nzdev")

	if err != nil {
		log.Fatalf("failed to get user: %v", err)
		return
	}

	log.Println(user)

	/*
		 	orders, err := dynamodb.GetOrdersByUserAndStatus("a1nzdev", "SHIPPED")

			if err != nil {
				log.Fatalf("Error getting orders: %v", err)
				return
			}

			for _, order := range orders {
				fmt.Printf("ID: %s, Status: %s, Created at: %s\n", order.IdOrder, order.Status, order.CreatedAt)
			}
	*/
}
