package iam

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
)

func CreatePolicy(policyName string, actions []string,
	resourceArn string) (*types.Policy, error) {
	var policy *types.Policy

	policyDoc := PolicyDocument{
		Version: "2012-10-17",
		Statement: []PolicyStatement{{
			Effect:   "Allow",
			Action:   actions,
			Resource: aws.String(resourceArn),
		}},
	}

	policyBytes, err := json.Marshal(policyDoc)

	if err != nil {
		log.Printf("Couldn't create policy document for %v. Here's why: %v\n", resourceArn, err)
		return nil, err
	}

	result, err := client.CreatePolicy(ctx, &iam.CreatePolicyInput{
		PolicyDocument: aws.String(string(policyBytes)),
		PolicyName:     aws.String(policyName),
	})

	if err != nil {
		log.Printf("Couldn't create policy %v. Here's why: %v\n", policyName, err)
	}

	policy = result.Policy

	return policy, err
}
