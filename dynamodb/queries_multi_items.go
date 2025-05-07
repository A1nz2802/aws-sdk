package dynamodb

import (
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// get orders by a particular user
func GetOrders(username string) ([]Order, error) {
	pk := "USER#" + username

	output, err := client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(tableName),
		KeyConditionExpression: aws.String("PK = :u AND begins_with(SK, :s)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":u": &types.AttributeValueMemberS{Value: pk},
			":s": &types.AttributeValueMemberS{Value: "ORDER#"},
		},
	})

	if err != nil {
		log.Printf("failed to get user orders: %v", err)
		return nil, err
	}

	var orders []Order
	err = attributevalue.UnmarshalListOfMaps(output.Items, &orders)

	if err != nil {
		log.Printf("failed to unmarshal orders: %v", err)
	}

	return orders, nil
}
