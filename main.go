package main

import (
	"aws-sdk/dynamodb"
	"log"
)

// https://www.youtube.com/watch?v=DIQVJqiSUkE

func main() {

	user := dynamodb.UserDto{
		Username: "alexdebrie",
		Fullname: "Alex DeBrie",
		Email:    "alexdebrie1@gmail.com",
	}

	_, err := dynamodb.CreateTable()

	if err != nil {
		log.Printf("Error while trying to create a table")
		return
	}

	err = dynamodb.AddUser(user)

	if err != nil {
		log.Printf("Error inserting user: %v", err)
		return
	}

	address := dynamodb.AddressDto{
		Username: "alexdebrie",
		Label:    "Work",
		Address: dynamodb.Address{
			Street:      "123 Main St",
			PostalCode:  90211,
			State:       "CA",
			CountryCode: "USAA",
		},
	}

	err = dynamodb.AddAddress(address)

	if err != nil {
		log.Printf("Error al añadir dirección: %v", err)
		return
	}

	order := dynamodb.OrderDto{
		Username: "alexdebrie",
		OrderId:  "5e7272b7",
		Status:   "PLACED",
		Address: dynamodb.Address{
			Street:      "123 Main St",
			PostalCode:  90211,
			State:       "CA",
			CountryCode: "USAA",
		},
		Items: []dynamodb.OrderItemDto{{
			IdItem:      "ab970628",
			ProductName: "Mackbook Pro",
			Price:       1399.99,
			Status:      "FILLED",
		}, {
			IdItem:      "2eae1dee",
			ProductName: "Amazon Echo",
			Price:       69.99,
			Status:      "FILLED",
		}},
	}

	err = dynamodb.AddOrder(order)

	if err != nil {
		log.Printf("Error al añadir orden: %v", err)
		return
	}

}
