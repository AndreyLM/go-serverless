// aws lambda list-functions
package main

import (
	"encoding/json"
)

// FunctionList - functins list
type FunctionList struct {
	Functions []Function
}

// Function - aws function
type Function struct {
	FunctionName string
}

// NewFunctionList - get list from aws
func NewFunctionList() (list FunctionList, err error) {
	dt, err := run("aws", "lambda", "list-functions")
	if err != nil {
		return
	}
	err = json.Unmarshal(dt, &list)
	return
}

// HasFunction - ckecks if function exist
func (fl FunctionList) HasFunction(fname string) bool {
	for _, v := range fl.Functions {
		if v.FunctionName == fname {
			return true
		}
	}

	return false
}

// CreateLambdaFunction - run command for  creating lambda function
func CreateLambdaFunction(lambdaName, roleName, handlerName, backetName string) ([]byte, error) {
	return run("aws", "lambda", "create-function",
		"--function-name", lambdaName,
		"--runtime", "go1.x",
		"--role", roleName,
		"--handler", handlerName,
		"--code", "S3Bucket="+backetName+", S3Key="+lambdaName+".zip")
}

// UpdateLambdaFunction - run command for updating lambda function
func UpdateLambdaFunction(lambdaName, backetName string) ([]byte, error) {
	return run("aws", "lambda", "update-function-code",
		"--function-name", lambdaName,
		"--s3-bucket", backetName,
		"--s3-key", lambdaName+".zip")
}
