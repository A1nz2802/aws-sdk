package sqs

import (
	"aws-sdk/iam"
	"aws-sdk/sns"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	awssns "github.com/aws/aws-sdk-go-v2/service/sns"
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

func AttachSendMessagePolicy(queueUrl string, queueArn string, topicArn string) error {
	policyDoc := iam.PolicyDocument{
		Version: "2012-10-17",
		Statement: []iam.PolicyStatement{{
			Effect:    "Allow",
			Action:    []string{"sqs:SendMessage"},
			Principal: map[string]string{"Service": "sns.amazonaws.com"},
			Resource:  aws.String(queueArn),
			Condition: iam.PolicyCondition{"ArnEquals": map[string]string{"aws:SourceArn": topicArn}},
		}},
	}

	policyBytes, err := json.Marshal(policyDoc)

	if err != nil {
		return fmt.Errorf("❌ failed to marshal policy document for queue %s: %w", queueUrl, err)
	}

	_, err = client.SetQueueAttributes(ctx, &sqs.SetQueueAttributesInput{
		Attributes: map[string]string{
			string(types.QueueAttributeNamePolicy): string(policyBytes),
		},
		QueueUrl: aws.String(queueUrl),
	})

	if err != nil {
		return fmt.Errorf("❌ failed to set send message policy on queue %s from topic %s: %w", queueUrl, topicArn, err)
	}

	return nil
}

func SubscribeQueue(topicArn string, queueArn string) (string, error) {
	var subscriptionArn string
	var attributes map[string]string

	output, err := sns.GetClient().Subscribe(ctx, &awssns.SubscribeInput{
		Protocol:              aws.String("sqs"),
		TopicArn:              aws.String(topicArn),
		Attributes:            attributes,
		Endpoint:              aws.String(queueArn),
		ReturnSubscriptionArn: true,
	})

	if err != nil {
		return "", fmt.Errorf("❌ failed to subscribe queue %s to topic %s: %w", queueArn, topicArn, err)
	}

	subscriptionArn = *output.SubscriptionArn

	return subscriptionArn, err
}

func SubscribeQueueToTopic(queueName string, queueUrl string, topicName string, topicArn string) (string, error) {

	queueArn, err := GetQueueArn(queueUrl)
	if err != nil {
		return "", fmt.Errorf("❌ failed to retrieve ARN for queue %s: %w", queueUrl, err)
	}

	err = AttachSendMessagePolicy(queueUrl, queueArn, topicArn)
	if err != nil {
		return "", fmt.Errorf("❌ failed to attach send message policy to queue %s from topic %s: %w", queueUrl, topicArn, err)
	}

	log.Println("✅ Successfully attached an IAM policy to the queue allowing messages from the topic.")

	subscriptionArn, err := SubscribeQueue(topicArn, queueArn)
	if err != nil {
		return "", fmt.Errorf("❌ failed to subscribe queue %s to topic %s: %w", queueUrl, topicArn, err)
	}

	log.Printf("✅ Queue '%s' is now successfully subscribed to topic '%s' with subscription ARN '%s'.\n",
		queueName, topicName, subscriptionArn)

	return subscriptionArn, nil
}

func DeleteQueue(queueUrl string) error {
	_, err := client.DeleteQueue(ctx, &sqs.DeleteQueueInput{QueueUrl: aws.String(queueUrl)})

	if err != nil {
		log.Printf("Couldn't delete queue %v. Here's why: %v\n", queueUrl, err)
	}

	return err
}
