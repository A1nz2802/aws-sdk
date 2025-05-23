package dynamodb

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// LSI example
func GetOrdersByPrice(idOrder string) ([]any, error) {

	pk := "ORDER#" + idOrder

	output, err := client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(tableName),
		IndexName:              aws.String("OrderByPriceIndex"),
		KeyConditionExpression: aws.String("PK = :pk AND Price > :price"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk":    &types.AttributeValueMemberS{Value: pk},
			":price": &types.AttributeValueMemberN{Value: "80"},
		},
		ReturnConsumedCapacity: types.ReturnConsumedCapacityTotal,
	})

	if err != nil {
		log.Printf("failed to get orders by price: %v", err)
		return nil, err
	}

	fmt.Printf("Count: %d\n", output.Count)
	fmt.Printf("ScannedCount: %d\n", output.ScannedCount)
	fmt.Printf("Consumed Capacity: %f\n", *output.ConsumedCapacity.CapacityUnits)

	var result []any
	err = attributevalue.UnmarshalListOfMaps(output.Items, &result)

	if err != nil {
		log.Printf("failed to unmarshal orders: %v", err)
		return nil, err
	}

	for _, item := range result {
		fmt.Printf("%+v\n", item)
	}

	return result, nil
}
