package main

import (
	"errors"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/lambda/lambdaiface"
)

type MockLambdaClientGolden struct {
	lambdaiface.LambdaAPI
}

func (m *MockLambdaClientGolden) Invoke(input *lambda.InvokeInput) (*lambda.InvokeOutput, error) {
	res := &lambda.InvokeOutput{}
	return res, nil
}

type MockLambdaClientError struct {
	lambdaiface.LambdaAPI
}

func (m *MockLambdaClientError) Invoke(input *lambda.InvokeInput) (*lambda.InvokeOutput, error) {
	return nil, errors.New("Mock error")
}

func TestNewLambda(t *testing.T) {

	actual := reflect.TypeOf(newLambda())
	expect := reflect.TypeOf(&lambdaClient{
		client: lambda.New(session.New()),
	})
	if actual != expect {
		t.Fatalf("Failed test. expect is \n%v \nbut actual is \n%v", expect, actual)
	}
}

func TestInvokeGolden(t *testing.T) {

	l := lambdaClient{
		client: &MockLambdaClientGolden{},
	}

	input := &lambda.InvokeInput{
		FunctionName:   aws.String("TestFunc"),
		Payload:        []byte(""),
		InvocationType: aws.String("Event"),
	}

	actual, err := l.invoke(input)
	if err != nil {
		t.Fatalf("Failed test. Unexpected error is occured: %v", err)
	}

	expect := &lambda.InvokeOutput{}

	if !reflect.DeepEqual(actual, expect) {
		t.Fatalf("Failed test. expect is \n%v \nbut actual is \n%v", expect, actual)
	}
}

func TestInvokeError(t *testing.T) {
	l := lambdaClient{
		client: &MockLambdaClientError{},
	}

	input := &lambda.InvokeInput{
		FunctionName:   aws.String("TestFunc"),
		Payload:        []byte(""),
		InvocationType: aws.String("Event"),
	}

	if _, err := l.invoke(input); err != nil {
		actual := err.Error()
		expect := "Mock error"

		if actual != expect {
			t.Fatalf("Failed test. expect is \n%v \nbut actual is \n%v", expect, actual)
		}
	} else {
		t.Fatalf("Failed test. Expected error didn't occur.")
	}
}
