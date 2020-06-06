package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/slack-go/slack"
)

func main() {
	lambda.Start(handleRequest)
}

type privateMeta struct {
	ResponseURL string `json:"response_url"`
	ChannelID   string `json:"channel_id"`
}

func handleRequest(message slack.InteractionCallback) {
	// This sample lambda simply returns user's inputs.

	// Get private_metadata fields value.
	// NOTE: private_metadata field's type is simply "String". SlackHub stores some convenient values in the form of JSON.
	// In order to use them easily, you should unmarshal private_metadata field value to struct.
	var privateMeta privateMeta
	if err := json.Unmarshal([]byte(message.View.PrivateMetadata), &privateMeta); err != nil {
		log.Printf("[ERROR] Failed to unmarshal json: %v", err)
		return
	}

	// Get user input values.
	lunchInput := message.View.State.Values["lunch_block"]["lunch_action"].Value
	detailInput := message.View.State.Values["detail_block"]["detail_action"].Value

	msg := "What is your favorite lunch?\n - " + lunchInput + "\nTell us more!\n - " + detailInput

	// Send a message to Slack
	// Create request body
	param := `{"text": "` + msg + `"}`
	req, err := http.NewRequest("POST", privateMeta.ResponseURL, bytes.NewBuffer([]byte(param)))
	if err != nil {
		log.Printf("[ERROR] Failed to create request body: %v", err)
		return
	}

	// Create request header
	req.Header.Set("Content-Type", "application/json")

	// Send to slack
	c := &http.Client{}
	res, err := c.Do(req)
	if err != nil {
		log.Printf("[ERROR] Failed to send a message: %v", err)
		return
	}
	defer res.Body.Close()

	return
}
