package auth

import (
	"errors"
	"reflect"
	"testing"

	"github.com/slack-go/slack"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
)

type MockClientGolden struct {
	ssmiface.SSMAPI
}

func (m *MockClientGolden) GetParameter(*ssm.GetParameterInput) (*ssm.GetParameterOutput, error) {
	res := &ssm.GetParameterOutput{
		Parameter: &ssm.Parameter{Value: aws.String("TestValue")},
	}
	return res, nil
}

type MockClientError struct {
	ssmiface.SSMAPI
}

func (m *MockClientError) GetParameter(*ssm.GetParameterInput) (*ssm.GetParameterOutput, error) {
	return nil, errors.New("Mock error")
}

func TestNew(t *testing.T) {

	actual := reflect.TypeOf(New("ap-northeast-1"))
	expect := reflect.TypeOf(&Auth{})
	if actual != expect {
		t.Fatalf("Failed test. expect is \n%v \nbut actual is \n%v", expect, actual)
	}
}

func TestGetParamGolden(t *testing.T) {
	a := Auth{
		client: &MockClientGolden{},
	}

	actual, err := a.getParam("key")
	if err != nil {
		t.Fatalf("Failed test. Unexpected error is occured: %v", err)
	}

	expect := "TestValue"

	if actual != expect {
		t.Fatalf("Failed test. expect is \n%v \nbut actual is \n%v", expect, actual)
	}
}

func TestGetParamError(t *testing.T) {
	a := Auth{
		client: &MockClientError{},
	}

	if _, err := a.getParam("key"); err != nil {
		actual := err.Error()
		expect := "Mock error"

		if actual != expect {
			t.Fatalf("Failed test. expect is \n%v \nbut actual is \n%v", expect, actual)
		}
	} else {
		t.Fatalf("Failed test. Expected error didn't occur.")
	}
}

func TestClientGolden(t *testing.T) {
	a := Auth{
		client: &MockClientGolden{},
	}

	actual, err := a.Client("test")
	if err != nil {
		t.Fatalf("Failed test. Unexpected error is occured: %v", err)
	}

	expect := slack.New("TestValue")

	if !reflect.DeepEqual(actual, expect) {
		t.Fatalf("Failed test. expect is \n%v \nbut actual is \n%v", expect, actual)
	}
}

func TestClientError(t *testing.T) {
	a := Auth{
		client: &MockClientError{},
	}

	if _, err := a.Client("key"); err != nil {
		actual := err.Error()
		expect := "Mock error"

		if actual != expect {
			t.Fatalf("Failed test. expect is \n%v \nbut actual is \n%v", expect, actual)
		}
	} else {
		t.Fatalf("Failed test. Expected error didn't occur.")
	}
}
