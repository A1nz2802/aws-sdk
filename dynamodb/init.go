package dynamodb

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/joho/godotenv"
)

func Init() error {
	_, err := CreateTable()

	if err != nil {
		log.Printf("Error creating table: %v", err)
		return err
	}

	err = populate()

	if err != nil {
		log.Printf("Error populating data: %v", err)
		return err
	}

	return nil
}

func populate() error {
	for _, user := range Users {
		err := AddUser(user)

		if err != nil {
			log.Printf("Error inserting user: %v", err)
			return err
		}
	}

	for _, address := range Addresses {
		err := AddAddress(address)

		if err != nil {
			log.Printf("Error inserting address: %v", err)
			return err
		}
	}

	for _, order := range Orders {
		err := AddOrder(order)

		if err != nil {
			log.Printf("Error inserting order: %v", err)
			return err
		}
	}

	return nil
}

func GetClient() *dynamodb.Client {
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	return dynamodb.NewFromConfig(cfg)
}
