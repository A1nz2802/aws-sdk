package sqs

import (
	"errors"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

func GetQueueURL(queueName string) (string, error) {

	result, err := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
		QueueName: aws.String(queueName),
	})

	if err != nil {
		return "", errors.New("Error getting URL queue " + queueName)
	}

	return aws.ToString(result.QueueUrl), nil
}

func ListQueues() error {
	var queueUrls []string

	fmt.Println("Let's list the queues for your account.")

	paginator := sqs.NewListQueuesPaginator(client, &sqs.ListQueuesInput{})

	for paginator.HasMorePages() {

		output, err := paginator.NextPage(ctx)

		if err != nil {
			log.Printf("Couldn't get queues. Here's why: %v\n", err)
			return err
		}

		queueUrls = append(queueUrls, output.QueueUrls...)

	}

	if len(queueUrls) == 0 {
		fmt.Println("You don't have any queues!")
		return nil
	}

	for _, queueUrl := range queueUrls {
		fmt.Printf("\t%v\n", queueUrl)
	}

	return nil
}

func GetQueueArn(queueUrl string) (string, error) {
	var queueArn string
	arnAttributeName := types.QueueAttributeNameQueueArn

	attribute, err := client.GetQueueAttributes(ctx, &sqs.GetQueueAttributesInput{
		QueueUrl:       aws.String(queueUrl),
		AttributeNames: []types.QueueAttributeName{arnAttributeName},
	})

	if err != nil {
		log.Printf("Couldn't get ARN for queue %v. Here's why: %v\n", queueUrl, err)
	}

	queueArn = attribute.Attributes[string(arnAttributeName)]

	return queueArn, err
}
