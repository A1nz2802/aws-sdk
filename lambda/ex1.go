package lambda

import (
	"aws-sdk/iam"
	"bytes"
	"fmt"
	"log"
	"os"
	"time"
)

type Data struct {
	RoleName   string
	LambdaName string
}

func CreateExample1() (Data, error) {

	var err error
	const roleName = "my-role-01"
	const policyARN = "arn:aws:iam::aws:policy/service-role/AWSLambdaSQSQueueExecutionRole"
	servicePrincipal := map[string]string{"Service": "lambda.amazonaws.com"}

	const lambdaName = "first-ex"
	const handlerName = "handleRequest"

	defer func() {
		CleanupResourcesForEx1(roleName, lambdaName, policyARN, err)
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
		return Data{}, fmt.Errorf("❌ failed to create or retrieve IAM role '%s': %w", roleName, err)
	}

	log.Printf("✅ IAM Role '%s' created/retrieved successfully. ARN: %s\n", roleName, *role.Arn)

	err = iam.AttachRolePolicy(policyARN, roleName)
	if err != nil {
		return Data{}, fmt.Errorf("❌ couldn't attach policy %s to role %s: %w", policyARN, roleName, err)
	}

	log.Printf("✅ Policy %s attached to role %s.\n", policyARN, roleName)

	zipContent, err := os.ReadFile("lambda.zip")
	if err != nil {
		return Data{}, fmt.Errorf("❌ failed to read lambda zip file: %w", err)
	}

	zipPackage := bytes.NewBuffer(zipContent)

	log.Println("⏳ Waiting 10 seconds for IAM role policy propagation...")
	time.Sleep(10 * time.Second)

	// Create lambda function
	functionState, err := CreateFunction(lambdaName, handlerName, role.Arn, zipPackage)
	if err != nil {
		return Data{}, fmt.Errorf("❌ failed to create Lambda function '%s': %w", lambdaName, err)
	}

	log.Printf("✅ Lambda function '%s' creation process finished. Current state: %s\n", lambdaName, functionState)

	return Data{
		RoleName:   roleName,
		LambdaName: lambdaName,
	}, nil
}

func CleanupResourcesForEx1(roleName, lambdaName string, policyArn string, opErr error) {
	if opErr == nil {
		log.Println("✅ No error detected. Skipping resource cleanup.")
		return
	}

	log.Printf("ERROR DETECTED. Attempting cleanup for resources '%s' and '%s'. Original error: %v", roleName, lambdaName, opErr)
	log.Printf("⏳ Attempting to delete Lambda function '%s'...", lambdaName)

	err := DeleteFunction(lambdaName)

	if err != nil {
		log.Printf("⚠️ WARNING: Error during cleanup of Lambda function '%s': %v", lambdaName, err)
	}

	if err == nil {
		log.Println("✅ Lambda function deleted sucessfully")
	}

	log.Printf("⏳ Trying to detach role to policy")

	err = iam.DetachRolePolicy(roleName, policyArn)

	if err != nil {
		log.Printf("⚠️ WARNING: Error during detaching role to '%s' policy: %v", roleName, err)
	}

	log.Printf("Attempting to delete IAM role '%s'...", roleName)

	if delErr := iam.DeleteRole(roleName); delErr != nil {
		log.Printf("⚠️ WARNING: Error during cleanup of IAM role '%s': %v", roleName, delErr)
	}

	log.Printf("✅ IAM role '%s' delete initiated (cleanup).", roleName)
}
