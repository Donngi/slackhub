package main

import (
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/nicoJN/slackhub/auth"
	"github.com/nicoJN/slackhub/tool"
	"github.com/slack-go/slack"
)

var (
	region    = os.Getenv("REGION")
	tokenKey  = os.Getenv("PARAMKEY_BOT_USER_AUTH_TOKEN")
	secretKey = os.Getenv("PARAMKEY_SIGNING_SECRET")
	dynamodb  = os.Getenv("DYNAMO_DB_NAME")
)

type privateMeta struct {
	ResponseURL string `json:"response_url"`
	ChannelID   string `json:"channel_id"`
}

func main() {
	lambda.Start(handleRequest)
}

func handleRequest(message slack.InteractionCallback) {

	var privateMeta privateMeta
	if err := json.Unmarshal([]byte(message.View.PrivateMetadata), &privateMeta); err != nil {
		log.Printf("[ERROR] Failed to unmarshal json: %v", err)
		return
	}

	// Create a slack client
	api, err := auth.New(region).Client(tokenKey)
	if err != nil {
		log.Printf("[ERROR] Failed to create slack client: %v", err)
		return
	}

	// Validate ID.
	id := message.View.State.Values[blockIDToolID][actionIDtoolID].Value
	if err := validateID(id); err != nil {
		if _, _, err := api.PostMessage(privateMeta.ChannelID, slack.MsgOptionText("[ERROR] ID you sent contains \":\" character.", false)); err != nil {
			log.Printf("[ERROR] Failed to send a message to Slack: %v", err)
		}
		log.Printf("[ERROR] ID in the request contains \":\" character")
		return
	}

	// Validate modalJSON. Register tool only checks if the json format is correct.
	modalJSON := message.View.State.Values[blockIDModalJSON][actionIDModalJSON].Value
	if err := validateModalJSON(modalJSON); err != nil {
		if _, _, err := api.PostMessage(privateMeta.ChannelID, slack.MsgOptionText("[ERROR] modalJSON you sent is not correct json format. Please check its syntax.", false)); err != nil {
			log.Printf("[ERROR] Failed to send a message to Slack: %v", err)
		}
		log.Printf("[ERROR] modalJSON in the request is not correct json format.")
		return
	}

	// Create new tool's metadata set.
	new := tool.NewTool(
		message.View.State.Values[blockIDToolID][actionIDtoolID].Value,
		message.View.State.Values[blockIDDisplayName][actionIDDisplayName].Value,
		message.View.State.Values[blockIDDescription][actionIDDescription].Value,
		message.View.State.Values[blockIDModalJSON][actionIDModalJSON].Value,
		message.View.State.Values[blockIDCalleeArn][actionIDCalleeArn].Value,
		message.View.State.Values[blockIDAdministrators][actionIDAdministrators].SelectedUsers,
		message.View.State.Values[blockIDAuthorizedUsers][actionIDAuthorizedUsers].SelectedUsers,
		message.View.State.Values[blockIDBootMode][actionIDBootMode].SelectedOption.Value)

	// Register
	if err := new.Register(region, dynamodb); err != nil {
		if _, _, err := api.PostMessage(privateMeta.ChannelID, slack.MsgOptionText("[ERROR] Failed to register new tool", false)); err != nil {
			log.Printf("[ERROR] Failed to send a message to Slack: %v", err)
		}
		log.Printf("[ERROR] Failed to register new tool: %v", err)
		return
	}

	// Send success message
	if _, _, err := api.PostMessage(privateMeta.ChannelID, slack.MsgOptionText(":confetti_ball::confetti_ball: Successfully registered \"*"+new.DisplayName+"*\" :confetti_ball::confetti_ball:", false)); err != nil {
		log.Printf("[ERROR] Failed to send a message to Slack: %v", err)
	}
	return

}

// validateID checks whether tool id contains a colon.
func validateID(id string) error {
	if strings.Contains(id, ":") {
		return errInvalidIDException
	}
	return nil
}

// validateModalJSON checks if the json format is correct. (CAUNTION: This method doesn't check if the json meet slack's modal format)
func validateModalJSON(j string) error {
	if !json.Valid([]byte(j)) && len(j) != 0 {
		return errInvalidModalJSONException
	}
	return nil
}
