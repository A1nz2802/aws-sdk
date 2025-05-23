package dynamodb

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func AddUser(user UserDto) error {
	userFormatted := User{
		PK:        "USER#" + user.Username,
		SK:        "PROFILE#" + user.Username,
		Username:  user.Username,
		Fullname:  user.Fullname,
		Email:     user.Email,
		CreatedAt: time.Now().Format(time.RFC3339),
		Addresses: map[string]Address{},
	}

	item, err := attributevalue.MarshalMap(userFormatted)
	if err != nil {
		log.Printf("failed to marshal user data for %s: %v\n", user.Username, err)
		return err
	}

	_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(tableName), Item: item,
		ConditionExpression: aws.String("attribute_not_exists(PK) AND attribute_not_exists(SK)"),
	})

	if err != nil {
		var condCheckErr *types.ConditionalCheckFailedException
		if errors.As(err, &condCheckErr) {
			return fmt.Errorf("user %s already exists: %w", user.Username, err)
		}
		log.Printf("couldn't add item: %v", err)
		return err
	}

	return nil
}

// USER + ADDRESS
func AddAddress(data AddressDto) error {
	username := data.Username
	pkValue := "USER#" + username
	skValue := "PROFILE#" + username
	attibuteName := "Addresses."

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
			PK:          "ORDER#" + data.OrderId,
			SK:          "ITEM#" + orderItemDto.IdItem,
			IdOrder:     data.OrderId,
			IdItem:      orderItemDto.IdItem,
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
			"#u": "Username",
			"#e": "Email",
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

func UpdateEmailUser(username string, email string) error {
	pk := "USER#" + username
	sk := "PROFILE#" + username

	_, err := client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: pk},
			"SK": &types.AttributeValueMemberS{Value: sk},
		},
		UpdateExpression:         aws.String("SET #e = :email"),
		ExpressionAttributeNames: map[string]string{"#e": "Email"},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":email": &types.AttributeValueMemberS{Value: email},
		},
	})

	if err != nil {
		log.Printf("failed to update email user")
		return err
	}

	return nil
}

func RemoveCreatedAtAttribute(username string) (User, error) {
	pk := "USER#" + username
	sk := "PROFILE#" + username

	output, err := client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: pk},
			"SK": &types.AttributeValueMemberS{Value: sk},
		},
		UpdateExpression: aws.String("REMOVE #c"),
		ExpressionAttributeNames: map[string]string{
			"#c": "Created_at",
		},
		ReturnValues: types.ReturnValueAllNew,
	})

	if err != nil {
		log.Printf("failed to remove created_at: %v", err)
		return User{}, err
	}

	var user User

	if len(output.Attributes) == 0 {
		log.Printf("no attributes returned after removing created_at for user: %s", username)
		return User{}, nil
	}

	err = attributevalue.UnmarshalMap(output.Attributes, &user)

	if err != nil {
		log.Printf("failed to unmarshal updated user: %v", err)
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

func RemoveItemOrder(idItem string, idOrder string) error {
	pk := "ORDER#" + idOrder
	sk := "ITEM#" + idItem

	_, err := client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: pk},
			"SK": &types.AttributeValueMemberS{Value: sk},
		},
		ConditionExpression: aws.String("Price > :p"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":p": &types.AttributeValueMemberN{Value: "289"},
		},
	})

	if err != nil {
		log.Printf("failed to remove item order: %v", err)
	}

	return nil
}
