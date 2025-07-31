package iam

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/joho/godotenv"
)

var ctx = context.Background()
var client = GetClient()

func GetClient() *iam.Client {
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error loading .env file, %v", err)
	}

	cfg, err := config.LoadDefaultConfig(ctx)

	if err != nil {
		log.Fatalf("unable to load AWS SDK config, %v", err)
	}

	return iam.NewFromConfig(cfg)
}
