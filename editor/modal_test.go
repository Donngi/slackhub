package main

import (
	"testing"

	"github.com/Jimon-s/slackhub/tool"
)

func TestCreateModal(t *testing.T) {
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

	view := createModal(&testTool, "test", "test_url", "test_channelID")

	actual := view.Title.Text
	expect := "SlackHub - Editor"

	if actual != expect {
		t.Fatalf("Failed test. expected is \n%v \nbut actual is \n%v", expect, actual)
	}
}
