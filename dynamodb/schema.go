package dynamodb

import (
	"context"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// USER = USER#<username> + PROFILE#<username>
// USER_ADDRESS = NA + NA
// ORDER = USER#<username> + ORDER#<orderid>
// ORDER_ITEM = ITEM#<itemid> + ORDER#<orderid>

type Address struct {
	Street      string
	PostalCode  uint
	State       string
	CountryCode string
}

type User struct {
	PK string `dynamodbav:"PK"`
	SK string `dynamodbav:"SK"`

	Username  string               `dynamodbav:"username"`
	Fullname  string               `dynamodbav:"fullname"`
	Email     string               `dynamodbav:"email"`
	CreatedAt string               `dynamodbav:"created_at"`
	Addresses map[string][]Address `dynamodbav:"addresses"`
}

type Order struct {
	PK string
	SK string

	IdOrder   string
	Username  string
	Address   Address
	Status    string
	CreatedAt string
}

type Item struct {
	PK string
	SK string

	IdItem      string
	IdOrder     string
	ProductName string
	Price       float64
	Status      string
}

var ctx = context.Background()
var client = GetClient()
var tableName = "ecommerce"

func CreateTable() (*types.TableDescription, error) {
	table, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{{
			AttributeName: aws.String("PK"),
			AttributeType: types.ScalarAttributeTypeS,
		}, {
			AttributeName: aws.String("SK"),
			AttributeType: types.ScalarAttributeTypeS,
		}},
		KeySchema: []types.KeySchemaElement{{
			AttributeName: aws.String("PK"),
			KeyType:       types.KeyTypeHash,
		}, {
			AttributeName: aws.String("SK"),
			KeyType:       types.KeyTypeRange,
		}},
		TableName:   aws.String(tableName),
		BillingMode: types.BillingModePayPerRequest,
	})

	if err != nil {
		log.Printf("Couldn't create table %v. Here's why: %v\n", tableName, err)
	}

	waiter := dynamodb.NewTableExistsWaiter(client)

	err = waiter.Wait(ctx, &dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	}, 5*time.Minute)

	if err != nil {
		log.Printf("Wait for table exists failed. Here's why: %v\n", err)
	}

	log.Printf("Creating table :)")
	return table.TableDescription, err
}
