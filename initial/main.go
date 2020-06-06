package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/nicoJN/chatops-slackhub/auth"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

var (
	region      = os.Getenv("REGION")
	tokenKey    = os.Getenv("PARAMKEY_BOT_USER_AUTH_TOKEN")
	secretKey   = os.Getenv("PARAMKEY_SIGNING_SECRET")
	dynamodb    = os.Getenv("DYNAMO_DB_NAME")
	registerArn = os.Getenv("REGISTER_TOOL_ARN")
	editorArn   = os.Getenv("EDITOR_TOOL_ARN")
	catalogArn  = os.Getenv("CATALOG_TOOL_ARN")
	eraserArn   = os.Getenv("ERASER_TOOL_ARN")
	sampleGoArn = os.Getenv("SAMPLE_TOOL_GO_ARN")
)

func main() {
	lambda.Start(handleRequest)
}

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// NOTE: In this app, we always send response with status 200 (not 401 or 500)
	// Slack offers us to send 200 response if we can receive message and send error messages in other way.
	// For more info, https://api.slack.com/interactivity/handling#responses

	// Slack authentication.
	if err := auth.New(region).Authorize(request, secretKey); err != nil {
		log.Printf("[ERROR] Failed to verify SigningSecret: %v", err)
		return createResponseStatus200(""), nil
	}

	// Analyze request's event type
	eventsAPIEvent, err := slackevents.ParseEvent(json.RawMessage(request.Body), slackevents.OptionNoVerifyToken())
	if err != nil {
		log.Printf("[ERROR] Failed to parse request body: %v", err)
		return createResponseStatus200(""), nil
	}

	switch eventsAPIEvent.Type {
	case slackevents.URLVerification:
		// Case: URL Verification. This logic is only called from slack developer's console when you set up your app.
		var r *slackevents.ChallengeResponse
		err := json.Unmarshal([]byte(request.Body), &r)
		if err != nil {
			log.Printf("[ERROR] Failed to unmarchal json: %v", err)
			return events.APIGatewayProxyResponse{Body: "[ERROR] Failed to unmarshal json", StatusCode: 401}, nil
		}
		return createResponseStatus200(r.Challenge), nil

	case slackevents.CallbackEvent:
		// Case: CallbackEvent.

		switch ev := eventsAPIEvent.InnerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:

			api, err := auth.New(region).Client(tokenKey)
			if err != nil {
				log.Printf("[ERROR] Failed to create slack client: %v", err)
				return createResponseStatus200(""), nil
			}

			// Create a menu.
			menu, err := createMenu()
			if err != nil {
				if _, _, err := api.PostMessage(ev.Channel, slack.MsgOptionText("[ERROR] Failed to create menu.", false)); err != nil {
					log.Printf("[ERROR] Failed to send a message to Slack: %v", err)
				}
				log.Printf("[ERROR] Failed to create menu: %v", err)
				return createResponseStatus200(""), nil
			}

			// Send a menu to slack channel.
			if _, _, err := api.PostMessage(ev.Channel, menu); err != nil {
				log.Printf("[ERROR] Failed to send a message to Slack: %v", err)
				return createResponseStatus200(""), nil
			}

			return createResponseStatus200(""), nil

		default:
			log.Printf("[ERROR] Unknown eventAPIEvent.InnerEvent: %T", eventsAPIEvent.InnerEvent.Data)
			return createResponseStatus200(""), nil
		}

	default:
		log.Printf("[ERROR] Unknown eventAPIEvent.Type: %v", eventsAPIEvent.Type)
		return createResponseStatus200(""), nil
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
