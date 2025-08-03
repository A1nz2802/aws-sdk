package sns

import (
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

func CreateTopic(topicName string, isFifoTopic bool, contentBasedDeduplication bool) (string, error) {
	var topicArn string
	topicAttributes := map[string]string{}

	if isFifoTopic {
		topicAttributes["FifoTopic"] = "true"
	}

	if contentBasedDeduplication {
		topicAttributes["ContentBasedDeduplication"] = "true"
	}

	topic, err := client.CreateTopic(ctx, &sns.CreateTopicInput{
		Name:       aws.String(topicName),
		Attributes: topicAttributes,
	})

	if err != nil {
		log.Printf("Couldn't create topic %v. Here's why: %v\n", topicName, err)
	}

	topicArn = *topic.TopicArn

	return topicArn, err
}

func DeleteTopic(topicArn string) error {
	_, err := client.DeleteTopic(ctx, &sns.DeleteTopicInput{
		TopicArn: aws.String(topicArn)})

	if err != nil {
		log.Printf("Couldn't delete topic %v. Here's why: %v\n", topicArn, err)
		return err
	}

	return nil
}
