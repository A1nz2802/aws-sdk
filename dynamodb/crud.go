package dynamodb

import (
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// USER + PROFILE
func AddUser(user UserDto) error {
	userFormated := User{
		PK:        "USER#" + user.Username,
		SK:        "PROFILE#" + user.Username,
		Username:  user.Username,
		Fullname:  user.Fullname,
		Email:     user.Email,
		CreatedAt: time.Now().Format(time.RFC3339),
		Addresses: map[string]Address{},
	}

	item, err := attributevalue.MarshalMap(userFormated)
	if err != nil {
		log.Printf("failed to marshal user data for %s: %v\n", user.Username, err)
		return err
	}
	_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(tableName), Item: item,
	})

	if err != nil {
		log.Printf("couldn't add item to table: %v\n", err)
		return err
	}

	return nil
}

// USER + ADDRESS
func AddAddress(data AddressDto) error {
	username := data.Username
	pkValue := "USER#" + username
	skValue := "PROFILE#" + username
	attibuteName := "addresses."

	upd := expression.NewBuilder().
		WithUpdate(expression.
			Set(expression.
				Name(attibuteName+data.Label), expression.Value(data.Address)))

	expr, err := upd.Build()
	if err != nil {
		log.Printf("failed to build expression for user %s: %v", username, err)
		return err
	}

	_, err = client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: pkValue},
			"SK": &types.AttributeValueMemberS{Value: skValue},
		},
		UpdateExpression:          expr.Update(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	})

	if err != nil {
		log.Printf("failed to update address for user %s: %v", username, err)
		return err
	}

	log.Printf("address added successfully for user: %v", username)
	return nil
}

// USER + ORDER
// ORDER + ITEM
func AddOrder(data OrderDto) error {
	order := Order{
		PK:        "USER#" + data.Username,
		SK:        "ORDER#" + data.OrderId,
		Username:  data.Username,
		IdOrder:   data.OrderId,
		Status:    data.Status,
		Address:   data.Address,
		CreatedAt: time.Now().Format(time.RFC3339),
	}

	orderAv, err := attributevalue.MarshalMap(order)

	if err != nil {
		return fmt.Errorf("failed to marshal order: %w", err)
	}

	var txnItems []types.TransactWriteItem

	txnItems = append(txnItems, types.TransactWriteItem{
		Put: &types.Put{
			TableName: aws.String(tableName),
			Item:      orderAv,
		},
	})

	for _, orderItemDto := range data.Items {
		itemCreated := Item{
			PK:          "ITEM#" + orderItemDto.IdItem,
			SK:          "ORDER#" + data.OrderId,
			IdItem:      orderItemDto.IdItem,
			IdOrder:     data.OrderId,
			ProductName: orderItemDto.ProductName,
			Price:       orderItemDto.Price,
			Status:      orderItemDto.Status,
		}

		itemAv, err := attributevalue.MarshalMap(itemCreated)

		if err != nil {
			return fmt.Errorf("marshal order item %s: %w", orderItemDto.IdItem, err)
		}

		txnItems = append(txnItems, types.TransactWriteItem{
			Put: &types.Put{
				TableName: aws.String(tableName),
				Item:      itemAv,
			},
		})
	}

	_, err = client.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{
		TransactItems: txnItems,
	})

	if err != nil {
		return fmt.Errorf("transaction failed: %w", err)
	}

	println("Order created successufully :)")

	return nil
}
