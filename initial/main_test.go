package main

import (
	"reflect"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestCreateResponseStatus200(t *testing.T) {
	body := "test_body"

	actual := createResponseStatus200(body)
	expect := events.APIGatewayProxyResponse{
		StatusCode:      200,
		IsBase64Encoded: false,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: "test_body",
	}

	if !reflect.DeepEqual(actual, expect) {
		t.Fatalf("Failed test. expect is \n%v \nbut actual is \n%v", expect, actual)
	}
}
