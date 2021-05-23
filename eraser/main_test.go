package main

import (
	"testing"

	"github.com/Jimon-s/slackhub/tool"
)

func TestValidateConfirmation(t *testing.T) {

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
	confirmation := "test DisplayName"

	if err := validateConfirmation(confirmation, &testTool); err != nil {
		t.Fatalf("Failed test. Unexpected error occured: %v", err)
	}
}

func TestValidateIDError(t *testing.T) {
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
	confirmation := "Wrong confirmation"

	if err := validateConfirmation(confirmation, &testTool); err != nil {
		actual := err
		expect := errInvalidConfirmationException

		if actual != expect {
			t.Fatalf("Failed test. expect is \n%v \nbut actual is \n%v", expect, actual)
		}
	} else {
		t.Fatalf("Failed test. Expected error didn't occur.")
	}
}
