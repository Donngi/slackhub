package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/nicoJN/slackhub/tool"
	"github.com/slack-go/slack"
)

func TestSendCatalogGolden(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"test":"value"}`))
			return
		},
	))
	defer ts.Close()

	api := slack.New("testToken", slack.OptionAPIURL(ts.URL+"/"))

	var testTools []tool.Tool
	for i := 0; i < 10; i++ {
		to := tool.Tool{
			ID:              "test",
			DisplayName:     "test DisplayName",
			Description:     "test Description",
			ModalJSON:       "test ModalJSON",
			CalleeArn:       "test CalleeArn",
			Administrators:  []string{"test_admin1", "test_admin2"},
			AuthorizedUsers: []string{"test_user1", "test_user2"},
			BootMode:        "Normal",
		}
		testTools = append(testTools, to)
	}

	if err := sendCatalog(testTools, api, "test_channel", "test_user"); err != nil {
		t.Fatalf("Failed test. Unexpected error is occured: %v", err)
	}
}

func TestSendCatalogError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(500)
			return
		},
	))
	defer ts.Close()

	api := slack.New("testToken", slack.OptionAPIURL(ts.URL+"/"))

	var testTools []tool.Tool
	for i := 0; i < 10; i++ {
		to := tool.Tool{
			ID:              "test",
			DisplayName:     "test DisplayName",
			Description:     "test Description",
			ModalJSON:       "test ModalJSON",
			CalleeArn:       "test CalleeArn",
			Administrators:  []string{"test_admin1", "test_admin2"},
			AuthorizedUsers: []string{"test_user1", "test_user2"},
			BootMode:        "Normal",
		}
		testTools = append(testTools, to)
	}

	if err := sendCatalog(testTools, api, "test_channel", "test_user"); err != nil {
		actual := err.Error()
		expect := "slack server error: 500 Internal Server Error"

		if actual != expect {
			t.Fatalf("Failed test. expect is \n%v \nbut actual is \n%v", expect, actual)
		}
	} else {
		t.Fatalf("Failed test. Expected error didn't occur.")
	}
}

func TestDivideToolListSingle(t *testing.T) {
	var testTools []tool.Tool
	for i := 0; i < 10; i++ {
		to := tool.New()
		testTools = append(testTools, *to)
	}

	actual := divideToolList(testTools, 10)
	expect := [][]tool.Tool{testTools}

	if !reflect.DeepEqual(actual, expect) {
		t.Fatalf("Failed test. expect is \n%v \nbut actual is \n%v", expect, actual)
	}
}

func TestDivideToolListMultiple(t *testing.T) {
	var testTools []tool.Tool
	for i := 0; i < 15; i++ {
		to := tool.New()
		testTools = append(testTools, *to)
	}

	actual := divideToolList(testTools, 10)

	var item10 []tool.Tool
	for i := 0; i < 10; i++ {
		to := tool.New()
		item10 = append(item10, *to)
	}
	var item5 []tool.Tool
	for i := 0; i < 5; i++ {
		to := tool.New()
		item5 = append(item5, *to)
	}
	expect := [][]tool.Tool{item10, item5}

	if !reflect.DeepEqual(actual, expect) {
		t.Fatalf("Failed test. expect is \n%v \nbut actual is \n%v", expect, actual)
	}
}

func TestGetIcon(t *testing.T) {
	cases := []struct {
		id   string
		icon string
	}{
		{id: "register", icon: ":ballot_box_with_ballot:"},
		{id: "editor", icon: ":wrench:"},
		{id: "catalog", icon: ":green_book:"},
		{id: "eraser", icon: ":warning:"},
		{id: "other", icon: ":name_badge:"},
	}

	for _, c := range cases {
		actual := getIcon(c.id)
		expect := c.icon

		if actual != expect {
			t.Fatalf("Failed test. Case [id:%v, value: %v]. expect is \n%v \nbut actual is \n%v", c.id, c.icon, expect, actual)
		}
	}
}

func TestGetNameListGolden(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			var testUsers []slack.User
			for i := 0; i < 2; i++ {
				u := slack.User{
					Name: "Test user",
				}
				testUsers = append(testUsers, u)
			}

			w.Header().Set("Content-Type", "application/json")
			res, _ := json.Marshal(struct {
				Ok    bool         `json:"ok"`
				Users []slack.User `json:"users"`
			}{
				Ok:    true,
				Users: testUsers,
			})
			w.Write(res)
			return
		},
	))
	defer ts.Close()

	api := slack.New("testToken", slack.OptionAPIURL(ts.URL+"/"))

	actual, err := getNameList([]string{""}, api)
	if err != nil {
		t.Fatalf("Failed test. Unexpected error is occured: %v", err)
	}

	expect := []string{"Test user", "Test user"}
	if reflect.DeepEqual(actual, expect) == false {
		t.Fatalf("Failed test. expect is \n%v \nbut actual is \n%v", expect, actual)
	}

}

func TestGetNameListZero(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			var testUsers []slack.User
			w.Header().Set("Content-Type", "application/json")
			res, _ := json.Marshal(struct {
				Ok    bool         `json:"ok"`
				Users []slack.User `json:"users"`
			}{
				Ok:    true,
				Users: testUsers,
			})
			w.Write(res)
			return
		},
	))
	defer ts.Close()

	api := slack.New("testToken", slack.OptionAPIURL(ts.URL+"/"))

	actual, err := getNameList([]string{}, api)
	if err != nil {
		t.Fatalf("Failed test. Unexpected error is occured: %v", err)
	}

	expect := []string{"All users"}
	if reflect.DeepEqual(actual, expect) == false {
		t.Fatalf("Failed test. expect is \n%v \nbut actual is \n%v", expect, actual)
	}

}

func TestGetNameListError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(500)
			return
		},
	))
	defer ts.Close()

	api := slack.New("testToken", slack.OptionAPIURL(ts.URL+"/"))

	if _, err := getNameList([]string{""}, api); err != nil {
		actual := err.Error()
		expect := "slack server error: 500 Internal Server Error"

		if actual != expect {
			t.Fatalf("Failed test. expect is \n%v \nbut actual is \n%v", expect, actual)
		}
	} else {
		t.Fatalf("Failed test. Expected error didn't occur.")
	}

}

func TestSliceToString(t *testing.T) {
	s := []string{"user1", "user2", "user3"}
	separator := ","

	actual := sliceToString(s, separator)
	expect := "user1,user2,user3"

	if actual != expect {
		t.Fatalf("Failed test. expect is \n%v \nbut actual is \n%v", expect, actual)
	}
}
