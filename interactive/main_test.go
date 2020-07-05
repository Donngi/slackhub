package main

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/nicoJN/slackhub/tool"
	"github.com/slack-go/slack"
)

func TestidentifyRequestTypeToolSelection(t *testing.T) {
	message := slack.InteractionCallback{
		ActionCallback: slack.ActionCallbacks{
			BlockActions: []*slack.BlockAction{
				&slack.BlockAction{
					ActionID: "slackhub_tool_selection",
				},
			},
		},
	}

	actual := identifyRequestType(message)
	expect := requestSlackHubToolSelection

	if actual != expect {
		t.Fatalf("Failed test. expected is \n%v \nbut actual is \n%v", expect, actual)
	}
}

func TestidentifyRequestTypeToolSelectionCancel(t *testing.T) {
	message := slack.InteractionCallback{
		ActionCallback: slack.ActionCallbacks{
			BlockActions: []*slack.BlockAction{
				&slack.BlockAction{
					ActionID: "slackhub_tool_selection_cancel",
				},
			},
		},
	}

	actual := identifyRequestType(message)
	expect := requestSlackHubToolSelectionCancel

	if actual != expect {
		t.Fatalf("Failed test. expected is \n%v \nbut actual is \n%v", expect, actual)
	}
}

func TestidentifyRequestTypeOther(t *testing.T) {
	message := slack.InteractionCallback{
		ActionCallback: slack.ActionCallbacks{
			BlockActions: []*slack.BlockAction{
				&slack.BlockAction{
					ActionID: "",
				},
			},
		},
	}

	actual := identifyRequestType(message)
	expect := requestOthers

	if actual != expect {
		t.Fatalf("Failed test. expected is \n%v \nbut actual is \n%v", expect, actual)
	}
}

func TestIdentifyToolIDWithSingleColon(t *testing.T) {
	actual := identifyToolID("testtool:xxxx")
	expect := "testtool"

	if actual != expect {
		t.Fatalf("Failed test. expected is \n%v \nbut actual is \n%v", expect, actual)
	}
}

func TestIdentifyToolIDWithoutColon(t *testing.T) {
	actual := identifyToolID("testtool")
	expect := "testtool"

	if actual != expect {
		t.Fatalf("Failed test. expected is \n%v \nbut actual is \n%v", expect, actual)
	}
}

func TestIdentifyToolIDWithDoubleColon(t *testing.T) {
	actual := identifyToolID("testtool:xxxx:yyyy")
	expect := "testtool"

	if actual != expect {
		t.Fatalf("Failed test. expected is \n%v \nbut actual is \n%v", expect, actual)
	}
}

func TestInvokeLambdaNormal(t *testing.T) {

	testTool := tool.Tool{
		ID:              "test",
		DisplayName:     "test DisplayName",
		Description:     "test Description",
		ModalJSON:       "test ModalJSON",
		CalleeArn:       "test CalleeArn",
		Administrators:  []string{"test_admin1", "test_admin2"},
		AuthorizedUsers: []string{"test_user1", "test_user2"},
		BootMode:        "Normal",
	}

	message := slack.InteractionCallback{}

	l := lambdaClient{
		client: &MockLambdaClientGolden{},
	}

	actual, err := l.invokeLambda(&testTool, message)
	if err != nil {
		t.Fatalf("Failed test. Unexpected error is occured: %v", err)
	}

	expect := &lambda.InvokeOutput{}

	if !reflect.DeepEqual(actual, expect) {
		t.Fatalf("Failed test. expect is \n%v \nbut actual is \n%v", expect, actual)
	}
}

func TestInvokeLambdaAdvanced(t *testing.T) {

	testTool := tool.Tool{
		ID:              "test",
		DisplayName:     "test DisplayName",
		Description:     "test Description",
		ModalJSON:       "test ModalJSON",
		CalleeArn:       "test CalleeArn",
		Administrators:  []string{"test_admin1", "test_admin2"},
		AuthorizedUsers: []string{"test_user1", "test_user2"},
		BootMode:        "Advanced",
	}

	message := slack.InteractionCallback{}

	l := lambdaClient{
		client: &MockLambdaClientGolden{},
	}

	actual, err := l.invokeLambda(&testTool, message)
	if err != nil {
		t.Fatalf("Failed test. Unexpected error is occured: %v", err)
	}

	expect := &lambda.InvokeOutput{}

	if !reflect.DeepEqual(actual, expect) {
		t.Fatalf("Failed test. expect is \n%v \nbut actual is \n%v", expect, actual)
	}
}

func TestInvokeLambdaUnknownMode(t *testing.T) {

	testTool := tool.Tool{
		ID:              "test",
		DisplayName:     "test DisplayName",
		Description:     "test Description",
		ModalJSON:       "test ModalJSON",
		CalleeArn:       "test CalleeArn",
		Administrators:  []string{"test_admin1", "test_admin2"},
		AuthorizedUsers: []string{"test_user1", "test_user2"},
		BootMode:        "unknown",
	}

	message := slack.InteractionCallback{}

	l := lambdaClient{
		client: &MockLambdaClientGolden{},
	}

	if _, err := l.invokeLambda(&testTool, message); err != nil {
		actual := err
		expect := errUnknownBootModeException

		if actual != expect {
			t.Fatalf("Failed test. expect is \n%v \nbut actual is \n%v", expect, actual)
		}
	} else {
		t.Fatalf("Failed test. Expected error didn't occur.")
	}
}

func TestInvokeLambdaError(t *testing.T) {

	testTool := tool.Tool{
		ID:              "test",
		DisplayName:     "test DisplayName",
		Description:     "test Description",
		ModalJSON:       "test ModalJSON",
		CalleeArn:       "test CalleeArn",
		Administrators:  []string{"test_admin1", "test_admin2"},
		AuthorizedUsers: []string{"test_user1", "test_user2"},
		BootMode:        "Normal",
	}

	message := slack.InteractionCallback{}

	l := lambdaClient{
		client: &MockLambdaClientError{},
	}

	if _, err := l.invokeLambda(&testTool, message); err != nil {
		actual := err.Error()
		expect := "Mock error"

		if actual != expect {
			t.Fatalf("Failed test. expect is \n%v \nbut actual is \n%v", expect, actual)
		}
	} else {
		t.Fatalf("Failed test. Expected error didn't occur.")
	}
}
