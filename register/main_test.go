package main

import "testing"

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
