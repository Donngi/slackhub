package tool

import (
	"errors"
	"reflect"
	"testing"

	"github.com/guregu/dynamo"
)

var sampleModalJSON = `
{
	"type": "modal",
	"title": {
		"type": "plain_text",
		"text": "Sample Tool",
		"emoji": true
	},
	"submit": {
		"type": "plain_text",
		"text": "Submit",
		"emoji": true
	},
	"close": {
		"type": "plain_text",
		"text": "Cancel",
		"emoji": true
	},
	"blocks": [
		{
			"type": "section",
			"text": {
				"type": "plain_text",
				"text": ":wave: Hi! This is a sample tool written in Go.",
				"emoji": true
			}
		},
		{
			"type": "divider"
		},
		{
            "block_id":"lunch_block",
			"type": "input",
			"element": {
                "action_id":"lunch_action",
				"type": "plain_text_input"
			},
			"label": {
				"type": "plain_text",
				"text": "What is your favorite lunch?",
				"emoji": true
			}
		},
		{
            "block_id":"detail_block",
			"type": "input",
			"label": {
				"type": "plain_text",
				"text": "Tell us more!",
				"emoji": true
			},
			"element": {
                "action_id":"detail_action",
				"type": "plain_text_input",
				"multiline": true
			}
		}
	]
}`

// Mock for golden case
type mockDbClientGolden struct{}

func (m *mockDbClientGolden) getOne(table dynamo.Table, key string, value string, t *Tool) error {
	*t = Tool{
		ID:              "test",
		DisplayName:     "test DisplayName",
		Description:     "test Description",
		ModalJSON:       sampleModalJSON,
		CalleeArn:       "test CalleeArn",
		Administrators:  []string{"test_admin1", "test_admin2"},
		AuthorizedUsers: []string{"test_user1", "test_user2"},
		BootMode:        "Normal",
	}
	return nil
}

func (m *mockDbClientGolden) put(table dynamo.Table, t *Tool) error {
	return nil
}

func (m *mockDbClientGolden) putIf(table dynamo.Table, t *Tool, condition string) error {
	return nil
}

func (m *mockDbClientGolden) delete(table dynamo.Table, key string, value string, t *Tool) error {
	return nil
}

// Mock for error case
type mockDbClientError struct{}

func (m *mockDbClientError) getOne(table dynamo.Table, key string, value string, t *Tool) error {
	return errors.New("Mock error")
}

func (m *mockDbClientError) put(table dynamo.Table, t *Tool) error {
	return errors.New("Mock error")
}

func (m *mockDbClientError) putIf(table dynamo.Table, t *Tool, condition string) error {
	return errors.New("Mock error")
}

func (m *mockDbClientError) delete(table dynamo.Table, key string, value string, t *Tool) error {
	return errors.New("Mock error")
}

func TestNew(t *testing.T) {
	actual := reflect.TypeOf(New().dbClient)
	expect := reflect.TypeOf(&dbClient{})
	if actual != expect {
		t.Fatalf("Failed test. expect is \n%v \nbut actual is \n%v", expect, actual)
	}
}

func TestGetItemGolden(t *testing.T) {
	testTool := Tool{}
	testTool.dbClient = &mockDbClientGolden{}

	if err := testTool.GetItem("test", "ap-northeast-1", "TestTable"); err != nil {
		t.Fatalf("Failed test. Unexpected error is occured: %v", err)
	}

	actual := testTool
	expect := Tool{
		ID:              "test",
		DisplayName:     "test DisplayName",
		Description:     "test Description",
		ModalJSON:       sampleModalJSON,
		CalleeArn:       "test CalleeArn",
		Administrators:  []string{"test_admin1", "test_admin2"},
		AuthorizedUsers: []string{"test_user1", "test_user2"},
		BootMode:        "Normal",
		dbClient:        &dbClient{},
	}

	if !reflect.DeepEqual(actual, expect) {
		t.Fatalf("Failed test. expect is \n%v \nbut actual is \n%v", expect, actual)
	}
}

func TestGetItemError(t *testing.T) {
	testTool := Tool{
		ID:              "test",
		DisplayName:     "test DisplayName",
		Description:     "test Description",
		ModalJSON:       sampleModalJSON,
		CalleeArn:       "test CalleeArn",
		Administrators:  []string{"test_admin1", "test_admin2"},
		AuthorizedUsers: []string{"test_user1", "test_user2"},
		BootMode:        "Normal",
	}

	testTool.dbClient = &mockDbClientError{}

	if err := testTool.GetItem("test", "ap-northeast-1", "TestTable"); err != nil {
		actual := err.Error()
		expect := "Mock error"

		if actual != expect {
			t.Fatalf("Failed test. expect is \n%v \nbut actual is \n%v", expect, actual)
		}
	} else {
		t.Fatalf("Failed test. Expected error didn't occur.")
	}
}

func TestGetAllItemsGoldem(t *testing.T) {

	t1 := Tool{
		ID:              "test1",
		DisplayName:     "test DisplayName",
		Description:     "test Description",
		ModalJSON:       sampleModalJSON,
		CalleeArn:       "test CalleeArn",
		Administrators:  []string{"test_admin1", "test_admin2"},
		AuthorizedUsers: []string{"test_user1", "test_user2"},
		BootMode:        "Normal",
	}

	t2 := Tool{
		ID:              "test2",
		DisplayName:     "test DisplayName",
		Description:     "test Description",
		ModalJSON:       sampleModalJSON,
		CalleeArn:       "test CalleeArn",
		Administrators:  []string{"test_admin1", "test_admin2"},
		AuthorizedUsers: []string{"test_user1", "test_user2"},
		BootMode:        "Normal",
	}

	t1updated := Tool{
		ID:              "test1",
		DisplayName:     "test DisplayName",
		Description:     "test Description",
		ModalJSON:       sampleModalJSON,
		CalleeArn:       "test CalleeArn",
		Administrators:  []string{"test_admin1", "test_admin2"},
		AuthorizedUsers: []string{"test_user1", "test_user2"},
		BootMode:        "Normal",
		dbClient:        &dbClient{},
	}

	t2updated := Tool{
		ID:              "test2",
		DisplayName:     "test DisplayName",
		Description:     "test Description",
		ModalJSON:       sampleModalJSON,
		CalleeArn:       "test CalleeArn",
		Administrators:  []string{"test_admin1", "test_admin2"},
		AuthorizedUsers: []string{"test_user1", "test_user2"},
		BootMode:        "Normal",
		dbClient:        &dbClient{},
	}

	// Mock scanAll method
	scanAll = func(table dynamo.Table, toolList *[]Tool) error {
		*toolList = []Tool{t1, t2}
		return nil
	}

	actual, err := GetAllItems("ap-north-east1", "TestTable")
	if err != nil {
		t.Fatalf("Failed test. Unexpected error occured: %v", err)
	}

	expect := []Tool{t1updated, t2updated}

	if !reflect.DeepEqual(actual, expect) {
		t.Fatalf("Failed test. expect is \n%v \nbut actual is \n%v", expect, actual)
	}
}

func TestGetAllItemsError(t *testing.T) {

	// Mock scanAll method
	scanAll = func(table dynamo.Table, toolList *[]Tool) error {
		return errors.New("Mock error")
	}

	if _, err := GetAllItems("ap-north-east1", "TestTable"); err != nil {
		actual := err.Error()
		expect := "Mock error"

		if actual != expect {
			t.Fatalf("Failed test. expect is \n%v \nbut actual is \n%v", expect, actual)
		}
	} else {
		t.Fatalf("Failed test. Expected error didn't occur.")
	}
}

func TestIsAdministratorsTrue(t *testing.T) {
	testTool := Tool{
		ID:              "test",
		DisplayName:     "test DisplayName",
		Description:     "test Description",
		ModalJSON:       sampleModalJSON,
		CalleeArn:       "test CalleeArn",
		Administrators:  []string{"test_admin1", "test_admin2"},
		AuthorizedUsers: []string{"test_user1", "test_user2"},
		BootMode:        "Normal",
	}

	actual := testTool.IsAdministrators("test_admin1")
	expect := true

	if actual != expect {
		t.Fatalf("Failed test. expect is \n%v \nbut actual is \n%v", expect, actual)
	}
}

func TestIsAdministratorsTrueEmpty(t *testing.T) {
	testTool := Tool{
		ID:              "test",
		DisplayName:     "test DisplayName",
		Description:     "test Description",
		ModalJSON:       sampleModalJSON,
		CalleeArn:       "test CalleeArn",
		Administrators:  []string{},
		AuthorizedUsers: []string{"test_user1", "test_user2"},
		BootMode:        "Normal",
	}

	actual := testTool.IsAdministrators("test_admin1")
	expect := true

	if actual != expect {
		t.Fatalf("Failed test. expect is \n%v \nbut actual is \n%v", expect, actual)
	}
}

func TestIsAdministratorsFalse(t *testing.T) {
	testTool := Tool{
		ID:              "test",
		DisplayName:     "test DisplayName",
		Description:     "test Description",
		ModalJSON:       sampleModalJSON,
		CalleeArn:       "test CalleeArn",
		Administrators:  []string{"test_admin1", "test_admin2"},
		AuthorizedUsers: []string{"test_user1", "test_user2"},
		BootMode:        "Normal",
	}

	actual := testTool.IsAdministrators("unknown_user")
	expect := false

	if actual != expect {
		t.Fatalf("Failed test. expect is \n%v \nbut actual is \n%v", expect, actual)
	}
}

func TestIsAuthorizedUsersTrue(t *testing.T) {
	testTool := Tool{
		ID:              "test",
		DisplayName:     "test DisplayName",
		Description:     "test Description",
		ModalJSON:       sampleModalJSON,
		CalleeArn:       "test CalleeArn",
		Administrators:  []string{"test_admin1", "test_admin2"},
		AuthorizedUsers: []string{"test_user1", "test_user2"},
		BootMode:        "Normal",
	}

	actual := testTool.IsAuthorizedUsers("test_user1")
	expect := true

	if actual != expect {
		t.Fatalf("Failed test. expect is \n%v \nbut actual is \n%v", expect, actual)
	}
}

func TestIsAuthorizedUsersTrueEmpty(t *testing.T) {
	testTool := Tool{
		ID:              "test",
		DisplayName:     "test DisplayName",
		Description:     "test Description",
		ModalJSON:       sampleModalJSON,
		CalleeArn:       "test CalleeArn",
		Administrators:  []string{"test_admin1", "test_admin2"},
		AuthorizedUsers: []string{},
		BootMode:        "Normal",
	}

	actual := testTool.IsAuthorizedUsers("test_user1")
	expect := true

	if actual != expect {
		t.Fatalf("Failed test. expect is \n%v \nbut actual is \n%v", expect, actual)
	}
}

func TestIsAuthorizedUsersFalse(t *testing.T) {
	testTool := Tool{
		ID:              "test",
		DisplayName:     "test DisplayName",
		Description:     "test Description",
		ModalJSON:       sampleModalJSON,
		CalleeArn:       "test CalleeArn",
		Administrators:  []string{"test_admin1", "test_admin2"},
		AuthorizedUsers: []string{"test_user1", "test_user2"},
		BootMode:        "Normal",
	}

	actual := testTool.IsAuthorizedUsers("unknown_user")
	expect := false

	if actual != expect {
		t.Fatalf("Failed test. expect is \n%v \nbut actual is \n%v", expect, actual)
	}
}

func TestIsUseModalTrue(t *testing.T) {
	testTool := Tool{
		ID:              "test",
		DisplayName:     "test DisplayName",
		Description:     "test Description",
		ModalJSON:       sampleModalJSON,
		CalleeArn:       "test CalleeArn",
		Administrators:  []string{"test_admin1", "test_admin2"},
		AuthorizedUsers: []string{"test_user1", "test_user2"},
		BootMode:        "Normal",
	}

	actual := testTool.IsUseModal()
	expect := true

	if actual != expect {
		t.Fatalf("Failed test. expect is \n%v \nbut actual is \n%v", expect, actual)
	}
}

func TestIsUseModalFalse(t *testing.T) {
	testTool := Tool{
		ID:              "test",
		DisplayName:     "test DisplayName",
		Description:     "test Description",
		ModalJSON:       "",
		CalleeArn:       "test CalleeArn",
		Administrators:  []string{"test_admin1", "test_admin2"},
		AuthorizedUsers: []string{"test_user1", "test_user2"},
		BootMode:        "Normal",
	}

	actual := testTool.IsUseModal()
	expect := false

	if actual != expect {
		t.Fatalf("Failed test. expect is \n%v \nbut actual is \n%v", expect, actual)
	}
}

func TestResisterGolden(t *testing.T) {
	testTool := Tool{
		ID:              "test",
		DisplayName:     "test DisplayName",
		Description:     "test Description",
		ModalJSON:       sampleModalJSON,
		CalleeArn:       "test CalleeArn",
		Administrators:  []string{"test_admin1", "test_admin2"},
		AuthorizedUsers: []string{"test_user1", "test_user2"},
		BootMode:        "Normal",
	}

	testTool.dbClient = &mockDbClientGolden{}

	if err := testTool.Register("ap-northeast-1", "TestTable"); err != nil {
		t.Fatalf("Failed test. Unexpected error is occured: %v", err)
	}
}

func TestResisterError(t *testing.T) {
	testTool := Tool{
		ID:              "test",
		DisplayName:     "test DisplayName",
		Description:     "test Description",
		ModalJSON:       sampleModalJSON,
		CalleeArn:       "test CalleeArn",
		Administrators:  []string{"test_admin1", "test_admin2"},
		AuthorizedUsers: []string{"test_user1", "test_user2"},
		BootMode:        "Normal",
	}

	testTool.dbClient = &mockDbClientError{}

	if err := testTool.Register("ap-northeast-1", "TestTable"); err != nil {
		actual := err.Error()
		expect := "Mock error"

		if actual != expect {
			t.Fatalf("Failed test. expect is \n%v \nbut actual is \n%v", expect, actual)
		}
	} else {
		t.Fatalf("Failed test. Expected error didn't occur.")
	}
}

func TestResisterForceGolden(t *testing.T) {
	testTool := Tool{
		ID:              "test",
		DisplayName:     "test DisplayName",
		Description:     "test Description",
		ModalJSON:       sampleModalJSON,
		CalleeArn:       "test CalleeArn",
		Administrators:  []string{"test_admin1", "test_admin2"},
		AuthorizedUsers: []string{"test_user1", "test_user2"},
		BootMode:        "Normal",
	}

	testTool.dbClient = &mockDbClientGolden{}

	if err := testTool.RegisterForce("ap-northeast-1", "TestTable"); err != nil {
		t.Fatalf("Failed test. Unexpected error is occured: %v", err)
	}
}

func TestResisterForceError(t *testing.T) {
	testTool := Tool{
		ID:              "test",
		DisplayName:     "test DisplayName",
		Description:     "test Description",
		ModalJSON:       sampleModalJSON,
		CalleeArn:       "test CalleeArn",
		Administrators:  []string{"test_admin1", "test_admin2"},
		AuthorizedUsers: []string{"test_user1", "test_user2"},
		BootMode:        "Normal",
	}

	testTool.dbClient = &mockDbClientError{}

	if err := testTool.RegisterForce("ap-northeast-1", "TestTable"); err != nil {
		actual := err.Error()
		expect := "Mock error"

		if actual != expect {
			t.Fatalf("Failed test. expect is \n%v \nbut actual is \n%v", expect, actual)
		}
	} else {
		t.Fatalf("Failed test. Expected error didn't occur.")
	}
}

func TestDeleteGolden(t *testing.T) {
	testTool := Tool{
		ID:              "test",
		DisplayName:     "test DisplayName",
		Description:     "test Description",
		ModalJSON:       sampleModalJSON,
		CalleeArn:       "test CalleeArn",
		Administrators:  []string{"test_admin1", "test_admin2"},
		AuthorizedUsers: []string{"test_user1", "test_user2"},
		BootMode:        "Normal",
	}

	testTool.dbClient = &mockDbClientGolden{}

	if err := testTool.Delete("ap-northeast-1", "TestTable"); err != nil {
		t.Fatalf("Failed test. Unexpected error is occured: %v", err)
	}
}

func TestDeleteError(t *testing.T) {
	testTool := Tool{
		ID:              "test",
		DisplayName:     "test DisplayName",
		Description:     "test Description",
		ModalJSON:       sampleModalJSON,
		CalleeArn:       "test CalleeArn",
		Administrators:  []string{"test_admin1", "test_admin2"},
		AuthorizedUsers: []string{"test_user1", "test_user2"},
		BootMode:        "Normal",
	}

	testTool.dbClient = &mockDbClientError{}

	if err := testTool.Delete("ap-northeast-1", "TestTable"); err != nil {
		actual := err.Error()
		expect := "Mock error"

		if actual != expect {
			t.Fatalf("Failed test. expect is \n%v \nbut actual is \n%v", expect, actual)
		}
	} else {
		t.Fatalf("Failed test. Expected error didn't occur.")
	}
}

func TestRemoveDuplication(t *testing.T) {
	s := []string{"A", "A", "B"}

	actual := removeDuplication(s)
	expect := []string{"A", "B"}

	if reflect.DeepEqual(actual, expect) == false {
		t.Fatalf("Failed test. expect is \n%v \nbut actual is \n%v", expect, actual)
	}
}

func TestRemoveElement(t *testing.T) {
	s := []string{"A", "B", "C"}

	actual := removeElement(s, "C")
	expect := []string{"A", "B"}

	if reflect.DeepEqual(actual, expect) == false {
		t.Fatalf("Failed test. expect is \n%v \nbut actual is \n%v", expect, actual)
	}
}

func TestRemoveElements(t *testing.T) {
	s := []string{"A", "B", "C"}
	e := []string{"A", "B"}

	actual := removeElements(s, e)
	expect := []string{"C"}

	if reflect.DeepEqual(actual, expect) == false {
		t.Fatalf("Failed test. expect is \n%v \nbut actual is \n%v", expect, actual)
	}
}

func TestSortTools(t *testing.T) {
	register := Tool{
		ID: "register",
	}

	editor := Tool{
		ID: "editor",
	}

	catalog := Tool{
		ID: "catalog",
	}

	eraser := Tool{
		ID: "eraser",
	}

	other := Tool{
		ID: "other",
	}

	raw := []Tool{other, editor, register, eraser, catalog}

	actual := SortTools(raw)
	expect := []Tool{other, register, editor, catalog, eraser}

	if reflect.DeepEqual(actual, expect) == false {
		t.Fatalf("Failed test. expect is \n%v \nbut actual is \n%v", expect, actual)
	}
}
