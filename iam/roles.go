package iam

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
)

func CreateRole(roleName string, trustPolicy PolicyDocument) (*types.Role, error) {
	var role *types.Role

	policyBytes, err := json.Marshal(trustPolicy)

	if err != nil {
		return nil, fmt.Errorf("couldn't marshal trust policy for role %s: %w", roleName, err)
	}

	result, err := client.CreateRole(ctx, &iam.CreateRoleInput{
		AssumeRolePolicyDocument: aws.String(string(policyBytes)),
		RoleName:                 aws.String(roleName),
	})

	if err != nil {
		return &types.Role{}, fmt.Errorf("couldn't create role %v. here's why: %v", roleName, err)
	}

	role = result.Role

	return role, err
}

func AttachRolePolicy(policyArn string, roleName string) error {
	_, err := client.AttachRolePolicy(ctx, &iam.AttachRolePolicyInput{
		PolicyArn: aws.String(policyArn),
		RoleName:  aws.String(roleName),
	})

	if err != nil {
		log.Printf("Couldn't attach policy %v to role %v. Here's why: %v\n", policyArn, roleName, err)
	}

	return err
}

func DetachRolePolicy(roleName string, policyArn string) error {
	_, err := client.DetachRolePolicy(ctx, &iam.DetachRolePolicyInput{
		PolicyArn: aws.String(policyArn),
		RoleName:  aws.String(roleName),
	})

	if err != nil {
		log.Printf("Couldn't detach policy from role %v. Here's why: %v\n", roleName, err)
	}

	return err
}

func GetRole(roleName string) (*types.Role, error) {
	var role *types.Role
	result, err := client.GetRole(ctx, &iam.GetRoleInput{RoleName: aws.String(roleName)})

	if err != nil {
		log.Printf("Couldn't get role %v. Here's why: %v\n", roleName, err)
	}

	role = result.Role

	return role, err
}

func DeleteRole(roleName string) error {

	_, err := client.DeleteRole(ctx, &iam.DeleteRoleInput{
		RoleName: aws.String(roleName),
	})

	if err != nil {
		return fmt.Errorf("couldn't delete function %v. Here's why: %w", roleName, err)
	}

	return nil
}
