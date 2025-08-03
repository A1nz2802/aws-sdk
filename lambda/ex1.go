package lambda

import (
	"aws-sdk/iam"
	"aws-sdk/sns"
	"aws-sdk/sqs"
	"bytes"
	"fmt"
	"log"
	"os"
	"time"
)

func CreateExample1() error {

	var err error
	const roleName = "my-role-01"
	const policyARN = "arn:aws:iam::aws:policy/service-role/AWSLambdaSQSQueueExecutionRole"
	servicePrincipal := map[string]string{"Service": "lambda.amazonaws.com"}

	const lambdaName = "first-ex-02"
	const handlerName = "handleRequest"

	var topicArn string
	var queueUrl string
	const topicName = "sns-test"
	const queueName = "my-std-queue"

	defer func() {
		CleanupResourcesForEx1(DataForCleanUp{
			lambdaName: lambdaName,
			roleName:   roleName,
			policyArn:  policyARN,
			topicName:  topicName,
			topicArn:   topicArn,
			queueName:  queueName,
			queueUrl:   queueUrl,
		}, err)
	}()

	// Create Policy Document
	trustPolicy := iam.PolicyDocument{
		Version: string(iam.DefaultVersionPolicy),
		Statement: []iam.PolicyStatement{{
			Effect:    "Allow",
			Principal: servicePrincipal, // "AWS": "trustedUserArn"
			Action:    []string{"sts:AssumeRole"},
		}},
	}

	// Create Role
	role, err := iam.CreateRole(roleName, trustPolicy)
	if err != nil {
		return fmt.Errorf("❌ failed to create or retrieve IAM role '%s': %w", roleName, err)
	}

	log.Printf("✅ IAM Role '%s' created/retrieved successfully. ARN: %s\n", roleName, *role.Arn)

	err = iam.AttachRolePolicy(policyARN, roleName)
	if err != nil {
		return fmt.Errorf("❌ couldn't attach policy %s to role %s: %w", policyARN, roleName, err)
	}

	log.Printf("✅ Policy %s attached to role %s.\n", policyARN, roleName)

	zipContent, err := os.ReadFile("lambda.zip")
	if err != nil {
		return fmt.Errorf("❌ failed to read lambda zip file: %w", err)
	}

	zipPackage := bytes.NewBuffer(zipContent)

	log.Println("⏳ Waiting 10 seconds for IAM role policy propagation...")
	time.Sleep(10 * time.Second)

	// Create lambda function
	functionState, err := CreateFunction(lambdaName, handlerName, role.Arn, zipPackage)
	if err != nil {
		return fmt.Errorf("❌ failed to create Lambda function '%s': %w", lambdaName, err)
	}

	log.Printf("✅ Lambda function '%s' creation process finished. Current state: %s\n", lambdaName, functionState)

	// Create sns topic
	topicArn, err = sns.CreateTopic(topicName, false, false)
	if err != nil {
		return fmt.Errorf("❌ failed to create SNS topic '%s': %w", topicName, err)
	}

	log.Printf("✅ SNS Topic '%s' created successfully. ARN: %s\n", topicName, topicArn)

	// Create sqs queue
	queueUrl, err = sqs.CreateQueue(false, queueName)
	if err != nil {
		return fmt.Errorf("❌ failed to create SQS queue '%s': %w", queueName, err)
	}

	log.Printf("✅ SQS Queue '%s' created successfully. URL: %s\n", queueName, queueUrl)

	_, err = sqs.SubscribeQueueToTopic(queueName, queueUrl, topicName, topicArn)
	if err != nil {
		return fmt.Errorf("❌ failed to subscribe SQS queue '%s' to SNS topic '%s': %w", queueName, topicName, err)
	}

	log.Printf("✅ SQS Queue '%s' subscribed to SNS Topic '%s'\n", queueName, topicName)

	return nil
}

type DataForCleanUp struct {
	lambdaName string

	roleName  string
	policyArn string

	topicName string
	topicArn  string

	queueName string
	queueUrl  string
}

func CleanupResourcesForEx1(data DataForCleanUp, opErr error) {
	if opErr == nil {
		log.Println("✅ No error detected. Skipping resource cleanup.")
		return
	}

	// Remove lambda function
	log.Printf("ERROR DETECTED. Attempting cleanup for resources")

	log.Printf("⏳ Attempting to delete Lambda function '%s'...", data.lambdaName)

	err := DeleteFunction(data.lambdaName)

	if err != nil {
		log.Printf("⚠️ WARNING: Error during cleanup of Lambda function '%s': %v", data.lambdaName, err)
	}

	log.Println("✅ Lambda function deleted sucessfully")

	// Detach role to policy
	log.Printf("⏳ Trying to detach role to policy")

	err = iam.DetachRolePolicy(data.roleName, data.policyArn)

	if err != nil {
		log.Printf("⚠️ WARNING: Error during detaching role to '%s' policy: %v", data.roleName, err)
	}

	log.Println("✅ Detach role to policy successfully")

	// Remove IAM Role
	log.Printf("⏳ Attempting to delete IAM role '%s'...", data.roleName)

	err = iam.DeleteRole(data.roleName)

	if err != nil {
		log.Printf("⚠️ WARNING: Error during cleanup of IAM role '%s': %v", data.roleName, err)
	}

	log.Printf("✅ IAM role '%s' remove sucessfully", data.roleName)

	// Remove SNS Topic
	log.Printf("⏳ Attempting to delete sns topic")

	err = sns.DeleteTopic(data.topicArn)

	if err != nil {
		log.Printf("⚠️ WARNING: Error during cleanup sns topic '%s': %v", data.topicArn, err)
	}

	log.Printf("✅ sns topic '%s' remove successfully", data.topicName)

	// Remove SQS queue
	err = sqs.DeleteQueue(data.queueUrl)

	if err != nil {
		log.Printf("⚠️ WARNING: Error during cleanup sqs queue '%s': %v", data.queueName, err)
	}

	log.Printf("✅ sqs queue '%s' remove successfully", data.queueName)
}
