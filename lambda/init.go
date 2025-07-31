package lambda

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/joho/godotenv"
)

var ctx = context.Background()
var client = GetClient()

func GetClient() *lambda.Client {
	err := godotenv.Load()

	if err != nil {
		// or panic
		log.Fatalf("Error loading .env file, %v", err)
	}

	cfg, err := config.LoadDefaultConfig(ctx)

	if err != nil {
		log.Fatalf("unable to load AWS SDK config, %v", err)
	}

	return lambda.NewFromConfig(cfg)
}
