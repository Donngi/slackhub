package main

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/slack-go/slack"
)

func TestIdentifyRequestTypeToolSelection(t *testing.T) {
	view := slack.View{
		CallbackID: "editor_tool_selection",
	}
	message := slack.InteractionCallback{
		View: view,
	}

	actual, err := identifyRequestType(message)
	if err != nil {
		t.Fatalf("Failed test. Unexpected error occured: %v", err)
	}

	expect := requestToolSelection

	if actual != expect {
		t.Fatalf("Failed test. expect is \n%v \nbut actual is \n%v", expect, actual)
	}
}

func TestIdentifyRequestTypeChangeRequest(t *testing.T) {
	view := slack.View{
		CallbackID: "editor_change_request",
	}
	message := slack.InteractionCallback{
		View: view,
	}

	actual, err := identifyRequestType(message)
	if err != nil {
		t.Fatalf("Failed test. Unexpected error occured: %v", err)
	}

	expect := requestChangeRequest

	if actual != expect {
		t.Fatalf("Failed test. expect is \n%v \nbut actual is \n%v", expect, actual)
	}
}

func TestIdentifyRequestTypeError(t *testing.T) {
	view := slack.View{
		CallbackID: "unknown",
	}
	message := slack.InteractionCallback{
		View: view,
	}

	if _, err := identifyRequestType(message); err != nil {
		actual := err
		expect := errUnknownRequestTypeException

		if actual != expect {
			t.Fatalf("Failed test. expect is \n%v \nbut actual is \n%v", expect, actual)
		}
	} else {
		t.Fatalf("Failed test. Expected error didn't occur.")
	}
}

func TestCreateErrorResponseStatus200(t *testing.T) {
	m := map[string]string{"test-block-id": "test-error"}
	actual := createErrorResponseStatus200(m)

	resAction := slack.NewErrorsViewSubmissionResponse(m)
	byte, _ := json.Marshal(resAction)

	expect := events.APIGatewayProxyResponse{
		StatusCode:      200,
		IsBase64Encoded: false,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(byte),
	}

	if !reflect.DeepEqual(actual, expect) {
		t.Fatalf("Failed test. expect is \n%v \nbut actual is \n%v", expect, actual)
	}
}

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

func TestValidateIDGolden(t *testing.T) {
	id := "testid"
	if err := validateID(id); err != nil {
		t.Fatalf("Failed test. Unexpected error occured: %v", err)
	}
}

func TestValidateIDError(t *testing.T) {
	id := "test:id"

	if err := validateID(id); err != nil {
		actual := err
		expect := errInvalidIDException

		if actual != expect {
			t.Fatalf("Failed test. expect is \n%v \nbut actual is \n%v", expect, actual)
		}
	} else {
		t.Fatalf("Failed test. Expected error didn't occur.")
	}
}

func TestValidateModalJSONGolden(t *testing.T) {
	json := "{\"test\":\"value\"}"
	if err := validateModalJSON(json); err != nil {
		t.Fatalf("Failed test. Unexpected error occured: %v", err)
	}
}

func TestValidateModalJSONEmpty(t *testing.T) {
	json := ""
	if err := validateModalJSON(json); err != nil {
		t.Fatalf("Failed test. Unexpected error occured: %v", err)
	}
}

func TestValidateModalJSONError(t *testing.T) {
	json := "invalid format"

	if err := validateModalJSON(json); err != nil {
		actual := err
		expect := errInvalidModalJSONException

		if actual != expect {
			t.Fatalf("Failed test. expect is \n%v \nbut actual is \n%v", expect, actual)
		}
	} else {
		t.Fatalf("Failed test. Expected error didn't occur.")
	}
}
