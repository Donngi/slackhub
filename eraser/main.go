package main

import (
	"encoding/json"
	"log"
	"os"

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
		log.Printf("[ERROR] Failed to create a slack client: %v", err)
		return
	}

	// Get selected tool's metadata
	tool := tool.New()
	if err := tool.GetItem(message.View.State.Values[blockIDToolSelection][actionIDToolSelection].SelectedOption.Value, region, dynamodb); err != nil {
		log.Printf("[ERROR] Failed to get tool's metadata from DynamoDB: %v", err)
		return
	}

	// Validate confirmation message
	confirmation := message.View.State.Values[blockIDConfirmation][actionIDConfirmation].Value
	if err := validateConfirmation(confirmation, tool); err != nil {
		resText := ":warning: Confirmation you sent and target tool's display name are different. Delete process was canceled.\nConfirmation: " + confirmation + "\nTool ID: " + tool.ID
		if _, _, err := api.PostMessage(privateMeta.ChannelID, slack.MsgOptionText(resText, false)); err != nil {
			log.Printf("[ERROR] Failed to send a message to Slack: %v", err)
		}
		return
	}

	// Delete the tool
	if err := tool.Delete(region, dynamodb); err != nil {
		if _, _, err := api.PostMessage(privateMeta.ChannelID, slack.MsgOptionText("[ERROR] Failed to delete tool's metadata from DynamoDB", false)); err != nil {
			log.Printf("[ERROR] Failed to send a message to Slack: %v", err)
		}
		log.Printf("[ERROR] Failed to delete tool's metadata from DynamoDB: %v", err)
		return
	}

	// Send success message
	if _, _, err := api.PostMessage(privateMeta.ChannelID, slack.MsgOptionText(":confetti_ball::confetti_ball: Successfully delete \"*"+tool.DisplayName+"*\" :confetti_ball::confetti_ball:", false)); err != nil {
		log.Printf("[ERROR] Failed to send a message to Slack: %v", err)
	}
	return

}

// validateConfirmation checks if a confirmation input is same as tool's display name.
func validateConfirmation(confirmation string, t *tool.Tool) error {
	if confirmation != t.DisplayName {
		return errInvalidConfirmationException
	}
	return nil
}
