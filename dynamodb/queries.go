package dynamodb

import (
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func GetUserByPk(username string) (User, error) {
	pk := "USER#" + username
	sk := "PROFILE#" + username

	output, err := client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: pk},
			"SK": &types.AttributeValueMemberS{Value: sk},
		},
		ProjectionExpression: aws.String("#u, #e"),
		ExpressionAttributeNames: map[string]string{
			"#u": "username",
			"#e": "email",
		},
	})

	if err != nil {
		log.Printf("failed to get user: %v", err)
		return User{}, err
	}

	if output.Item == nil {
		log.Printf("user with PK=%s and SK=%s doens't exist", pk, sk)
		return User{}, nil
	}

	var user User
	err = attributevalue.UnmarshalMap(output.Item, &user)

	if err != nil {
		log.Printf("failed to unmarshal user: %v", err)
		return User{}, err
	}

	return user, nil
}

func GetOrdersByUserAndStatus(username string, status string) ([]Order, error) {
	userOrderGSI := "UserOrdersGSI"

	output, err := client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(tableName),
		IndexName:              aws.String(userOrderGSI),
		KeyConditionExpression: aws.String("Username = :user"),
		FilterExpression:       aws.String("#s = :status"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":user":   &types.AttributeValueMemberS{Value: username},
			":status": &types.AttributeValueMemberS{Value: status},
		},
		ExpressionAttributeNames: map[string]string{"#s": "State"},
	})

	if err != nil {
		log.Printf("failed to query orders by user and status: %v", err)
		return nil, err
	}

	var orders []Order
	err = attributevalue.UnmarshalListOfMaps(output.Items, &orders)

	if err != nil {
		log.Printf("failed to unmarshal orders: %v", err)
		return nil, err
	}

	return orders, nil
}
