package main

import (
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/nicoJN/chatops-slackhub/auth"
	"github.com/nicoJN/chatops-slackhub/tool"
	"github.com/slack-go/slack"
)

type privateMeta struct {
	ResponseURL  string `json:"response_url,omitempty"`
	ChannelID    string `json:"channel_id,omitempty"`
	TargetToolID string `json:"target_tool_id,omitempty"`
}

var (
	region    = os.Getenv("REGION")
	tokenKey  = os.Getenv("PARAMKEY_BOT_USER_AUTH_TOKEN")
	secretKey = os.Getenv("PARAMKEY_SIGNING_SECRET")
	dynamodb  = os.Getenv("DYNAMO_DB_NAME")
)

func main() {
	lambda.Start(handleRequest)
}

func handleRequest(message slack.InteractionCallback) (events.APIGatewayProxyResponse, error) {

	api, err := auth.New(region).Client(tokenKey)
	if err != nil {
		log.Printf("[ERROR] Failed to create slack client: %v", err)
		return createResponseStatus200(""), nil
	}

	var privateMeta privateMeta
	if err := json.Unmarshal([]byte(message.View.PrivateMetadata), &privateMeta); err != nil {
		log.Printf("[ERROR] Failed to unmarshal json: %v", err)
		return createResponseStatus200(""), nil
	}

	// Identify request type
	reqType, err := identifyRequestType(message)
	if err != nil {
		if _, _, err := api.PostMessage(privateMeta.ChannelID, slack.MsgOptionText("[ERROR] Unknown requestType.", false)); err != nil {
			log.Printf("[ERROR] Failed to send a message to Slack: %v", err)
		}
		log.Printf("[ERROR] Unknown requestType: %v", message.View.CallbackID)
		return createResponseStatus200(""), nil
	}

	switch reqType {
	case requestToolSelection:
		// Get selected tool's metadata
		tool := tool.New()
		if err := tool.GetItem(message.View.State.Values[blockIDToolSelection][actionIDToolSelection].SelectedOption.Value, region, dynamodb); err != nil {
			if _, _, err := api.PostMessage(privateMeta.ChannelID, slack.MsgOptionText("[ERROR] Failed to get tool's metadata from DynamoDB.", false)); err != nil {
				log.Printf("[ERROR] Failed to send a message to Slack: %v", err)
			}
			log.Printf("[ERROR] Failed to get tool's metadata from DynamoDB: %v", err)
			return createResponseStatus200(""), nil
		}

		// Check whether a user is permitted to edit the selected tool
		if !tool.IsAdministrators(message.User.ID) {
			if _, _, err := api.PostMessage(privateMeta.ChannelID, slack.MsgOptionText("[ERROR] You are NOT permitted to edit "+tool.DisplayName+". Please contact administrators.", false)); err != nil {
				log.Printf("[ERROR] Failed to send a message to Slack: %v", err)
			}
			return createResponseStatus200(""), nil
		}

		// Create modal
		view := createModal(tool, message.User.ID, privateMeta.ResponseURL, privateMeta.ChannelID)

		// Create response
		resAction := slack.NewUpdateViewSubmissionResponse(view)
		byte, _ := json.Marshal(resAction)
		res := createResponseStatus200(string(byte))

		return res, nil

	case requestChangeRequest:
		action := message.View.State.Values

		// Validate modalJSON. Editor tool only checks if the json format is correct.
		modalJSON := action[blockIDModalJSON][actionIDModalJSON].Value
		if err := validateModalJSON(modalJSON); err != nil {
			log.Printf("[ERROR] Validation error: %v", err)
			errors := map[string]string{blockIDModalJSON: "[ERROR] Not correct json format. Please check its syntax."}
			return createErrorResponseStatus200(errors), nil
		}

		// Create tool's new metadata set.
		new := tool.NewTool(
			privateMeta.TargetToolID,
			message.View.State.Values[blockIDDisplayName][actionIDDisplayName].Value,
			message.View.State.Values[blockIDDescription][actionIDDescription].Value,
			message.View.State.Values[blockIDModalJSON][actionIDModalJSON].Value,
			message.View.State.Values[blockIDCalleeArn][actionIDCalleeArn].Value,
			message.View.State.Values[blockIDAdministrators][actionIDAdministrators].SelectedUsers,
			message.View.State.Values[blockIDAuthorizedUsers][actionIDAuthorizedUsers].SelectedUsers,
			message.View.State.Values[blockIDBootMode][actionIDBootMode].SelectedOption.Value)

		// Update tool's metadata
		if err := new.RegisterForce(region, dynamodb); err != nil {
			if _, _, err := api.PostMessage(privateMeta.ChannelID, slack.MsgOptionText("[ERROR] Failed to register new tool.", false)); err != nil {
				log.Printf("[ERROR] Failed to send a message to Slack: %v", err)
			}
			log.Printf("[ERROR] Failed to register new tool: %v", err)
			return createResponseStatus200(""), nil
		}

		// Send success message
		if _, _, err := api.PostMessage(privateMeta.ChannelID, slack.MsgOptionText(":confetti_ball::confetti_ball: Successfully edited \"*"+new.DisplayName+"*\" :confetti_ball::confetti_ball:", false)); err != nil {
			log.Printf("[ERROR] Failed to send a message to Slack: %v", err)
		}

		return createResponseStatus200(""), nil

	default:
		return createResponseStatus200(""), nil
	}
}

// identifyRequestType returns the result of determining request type.
func identifyRequestType(message slack.InteractionCallback) (requestType, error) {
	switch message.View.CallbackID {
	case string(requestToolSelection):
		return requestToolSelection, nil
	case string(requestChangeRequest):
		return requestChangeRequest, nil
	default:
		return requestType(""), errUnknownRequestTypeException
	}
}

// createErrorResponseStatus200 returns a APIGatewayProxyResponse with status 200 OK.
func createErrorResponseStatus200(errors map[string]string) events.APIGatewayProxyResponse {
	resAction := slack.NewErrorsViewSubmissionResponse(errors)
	byte, _ := json.Marshal(resAction)

	return events.APIGatewayProxyResponse{
		StatusCode:      200,
		IsBase64Encoded: false,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(byte),
	}
}

// createResponseStatus200 returns a APIGatewayProxyResponse with status 200 OK.
func createResponseStatus200(body string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode:      200,
		IsBase64Encoded: false,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: body,
	}
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
