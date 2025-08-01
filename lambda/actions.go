package lambda

import (
	"bytes"
	"errors"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
)

func CreateFunction(functionName string, handlerName string, iamRoleArn *string, zipPackage *bytes.Buffer) (types.State, error) {
	var state types.State

	_, err := client.CreateFunction(ctx, &lambda.CreateFunctionInput{
		Code:         &types.FunctionCode{ZipFile: zipPackage.Bytes()},
		FunctionName: aws.String(functionName),
		Role:         iamRoleArn,
		Handler:      aws.String(handlerName),
		Publish:      true,
		Runtime:      types.RuntimeGo1x,
	})

	if err != nil {
		var resConflict *types.ResourceConflictException

		if errors.As(err, &resConflict) {
			log.Printf("Function %v already exists.\n", functionName)
			return types.StateActive, nil
		}

		return "", fmt.Errorf("couldn't create function %v. Here's why: %w", functionName, err)
	}

	return state, nil
}

func DeleteFunction(functionName string) error {

	_, err := client.DeleteFunction(ctx, &lambda.DeleteFunctionInput{
		FunctionName: aws.String(functionName),
	})

	if err != nil {
		return fmt.Errorf("couldn't delete function %v. Here's why: %w", functionName, err)
	}

	return nil
}
