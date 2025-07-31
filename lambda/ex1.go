package lambda

import (
	"aws-sdk/iam"
	"bytes"
	"fmt"
	"log"
	"os"
)

func CreateExample1() error {

	const roleName = "my-role-01"
	const trustedUserArn = "lambda.amazonaws.com"
	const policyARN = "arn:aws:iam::aws:policy/service-role/AWSLambdaSQSQueueExecutionRole"

	// Create Role
	role, err := iam.CreateRole(roleName, trustedUserArn)

	if err != nil {

	}

	err = iam.AttachRolePolicy(policyARN, roleName)

	if err != nil {
		return fmt.Errorf("couldn't attach policy %s to role %s: %w", policyARN, roleName, err)
	}

	log.Printf("Policy %s attached to role %s.\n", policyARN, roleName)

	zipContent, err := os.ReadFile("lambda.zip")

	if err != nil {
		return fmt.Errorf("failed to read lambda zip: %w", err)
	}

	zipPackage := bytes.NewBuffer(zipContent)

	// Create function
	result, err := CreateFunction("first-ex", "some", role.Arn, zipPackage)

	if err != nil {

	}

	log.Println(result.Values())
	return nil
}
