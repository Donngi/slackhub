package main

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/lambda/lambdaiface"
)

// All methods in this file is lapper of Lambda access.
// It makes unit test easy.

type lambdaClient struct {
	client lambdaiface.LambdaAPI
}

func newLambda() *lambdaClient {
	return &lambdaClient{
		client: lambda.New(session.New()),
	}
}

// invoke calls other lambda.
func (l *lambdaClient) invoke(input *lambda.InvokeInput) (*lambda.InvokeOutput, error) {
	res, err := l.client.Invoke(input)
	if err != nil {
		return res, err
	}
	return res, nil
}
