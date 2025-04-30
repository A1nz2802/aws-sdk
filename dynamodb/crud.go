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
		Addresses: map[string][]Address{},
	}

	item, err := attributevalue.MarshalMap(userFormated)
	if err != nil {
		panic(err)
	}
	_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(tableName), Item: item,
	})
	if err != nil {
		log.Printf("Couldn't add item to table. Here's why: %v\n", err)
	}
	return err
}

// USER + ADDRESS
func AddAddress(data AddressDto) error {
	pkValue := "USER#" + data.Username
	skValue := "PROFILE#" + data.Username
	attibuteName := "addresses."

	upd := expression.NewBuilder().
		WithUpdate(expression.
			Set(expression.
				Name(attibuteName+data.Label), expression.Value(data.Address)))

	expr, err := upd.Build()
	if err != nil {
		return fmt.Errorf("error construyendo la expresión: %w", err)
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

	log.Printf("Dirección agrega para el usuario %v", data.Username)

	return err
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
