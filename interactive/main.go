package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/Jimon-s/slackhub/auth"
	"github.com/Jimon-s/slackhub/tool"
	"github.com/aws/aws-lambda-go/events"
	awslambda "github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/slack-go/slack"
)

var (
	region    = os.Getenv("REGION")
	tokenKey  = os.Getenv("PARAMKEY_BOT_USER_AUTH_TOKEN")
	secretKey = os.Getenv("PARAMKEY_SIGNING_SECRET")
	dynamodb  = os.Getenv("DYNAMO_DB_NAME")
)

func main() {
	awslambda.Start(handleRequest)
}

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// NOTE: In this app, we always send response with status 200 (not 401 or 500)
	// Slack offers us to send 200 response if we can receive message and send error messages in other way.
	// For more info, https://api.slack.com/interactivity/handling#responses

	// Decode request body
	body := request.Body
	if request.IsBase64Encoded {
		b, err := base64.StdEncoding.DecodeString(body)
		if err != nil {
			log.Printf("[ERROR] Failed to decode request body: %v", err)
			return createResponseStatus200(""), nil
		}
		body = string(b)
	}

	// Slack authentication.
	if err := auth.New(region).Authorize(body, request.Headers, secretKey); err != nil {
		log.Printf("[ERROR] Failed to verify SigningSecret: %v", err)
		return createResponseStatus200(""), nil
	}

	// Parse request
	payload, err := url.QueryUnescape(body)
	if err != nil {
		log.Printf("[ERROR] Failed to unescape: %v", err)
		return createResponseStatus200(""), nil
	}
	payload = strings.Replace(payload, "payload=", "", 1)

	var message slack.InteractionCallback
	if err := json.Unmarshal([]byte(payload), &message); err != nil {
		log.Printf("[ERROR] Failed to unmarshal json: %v", err)
		return createResponseStatus200(""), nil
	}

	// Create slack client
	api, err := auth.New(region).Client(tokenKey)
	if err != nil {
		log.Printf("[ERROR] Failed to create slack client: %v", err)
		return createResponseStatus200(""), nil
	}

	// Switch case based on requestType
	switch identifyRequestType(message) {
	case requestSlackHubToolSelection:
		// Get selected tool's metadata
		tool := tool.New()
		if err := tool.GetItem(message.ActionCallback.BlockActions[0].SelectedOption.Value, region, dynamodb); err != nil {
			if _, _, err := api.PostMessage(message.Channel.ID, slack.MsgOptionText("[ERROR] Failed to get tool's metadata from DynamoDB.", false)); err != nil {
				log.Printf("[ERROR] Failed to send a message to Slack: %v", err)
			}
			log.Printf("[ERROR] Failed to get tool's metadata from DynamoDB: %v", err)
			return createResponseStatus200(""), nil
		}

		// Check whether a user is permitted to use the selected tool
		if !tool.IsAuthorizedUsers(message.User.ID) {
			resText := "[ERROR] You are NOT permitted to use " + tool.DisplayName + ". Please contact administrators."
			if _, _, err := api.PostMessage(message.Channel.ID, slack.MsgOptionText(resText, false)); err != nil {
				log.Printf("[ERROR] Failed to send a message to Slack: %v", err)
			}
			return createResponseStatus200(""), nil
		}

		// Check whether the tool uses modal or not
		if !tool.IsUseModal() {
			// Invoke another lambda
			l := newLambda()
			res, err := l.invokeLambda(tool, message)
			if err != nil {
				if _, _, err := api.PostMessage(message.Channel.ID, slack.MsgOptionText("[ERROR] Something occured when calling your tool's lambda.", false)); err != nil {
					log.Printf("[ERROR] Failed to send a message to Slack: %v", err)
				}
				log.Printf("[ERROR] Failed to invoke lambda: %v", err)
				return createResponseStatus200(""), err
			}

			// Delete the original message
			if _, _, err := api.DeleteMessage(message.Channel.ID, message.Message.Timestamp); err != nil {
				if _, _, err := api.PostMessage(message.Channel.ID, slack.MsgOptionText("[WARNING] Successfully called your tool, but failed to delete the original message. Please see CloudWatch.", false)); err != nil {
					log.Printf("[ERROR] Failed to send a message to Slack: %v", err)
				}
				log.Printf("[ERROR] Failed to delete the original message: %v", err)
				return createResponseStatus200(""), nil
			}

			// Create a response
			switch tool.BootMode {
			case "Normal":
				// Send success message
				resText := "<@" + message.User.ID + "> invokes \"" + tool.DisplayName + "\" :tophat:"
				if _, _, err := api.PostMessage(message.Channel.ID, slack.MsgOptionText(resText, false)); err != nil {
					log.Printf("[ERROR] Failed to send a message to Slack: %v", err)
					return createResponseStatus200(""), nil
				}

				// Close the view (NOTE: We can close the view by only sending HTTP 200 response.)
				return createResponseStatus200(""), nil

			case "Advanced":
				var proxyRes events.APIGatewayProxyResponse
				if err != json.Unmarshal(res.Payload, &proxyRes) {
					return createResponseStatus200(""), err
				}
				return proxyRes, nil

			default:
				return createResponseStatus200(""), nil
			}
		}

		// Create a modal
		view, err := createModal(tool, message.User.ID, message.ResponseURL, message.Channel.ID)
		if err != nil {
			switch err {
			case errNonEditableToolsException:
				if _, _, err := api.PostMessage(message.Channel.ID, slack.MsgOptionText(":warning: You don't have any editable tools. :warning:", false)); err != nil {
					log.Printf("[ERROR] Failed to send success message: %v", err)
					return createResponseStatus200(""), nil
				}

				// Delete the original message
				if _, _, err := api.DeleteMessage(message.Channel.ID, message.Message.Timestamp); err != nil {
					if _, _, err := api.PostMessage(message.Channel.ID, slack.MsgOptionText("[ERROR] Failed to delete the original message. Please see CloudWatch.", false)); err != nil {
						log.Printf("[ERROR] Failed to send a message to Slack: %v", err)
					}
					log.Printf("[ERROR] Failed to delete the original message: %v", err)
					return createResponseStatus200(""), nil
				}
				return createResponseStatus200(""), nil
			case errNonDeletableToolsException:
				if _, _, err := api.PostMessage(message.Channel.ID, slack.MsgOptionText(":warning: You don't have any deletable tools. :warning:", false)); err != nil {
					log.Printf("[ERROR] Failed to send success message: %v", err)
					return createResponseStatus200(""), nil
				}

				// Delete the original message
				if _, _, err := api.DeleteMessage(message.Channel.ID, message.Message.Timestamp); err != nil {
					if _, _, err := api.PostMessage(message.Channel.ID, slack.MsgOptionText("[ERROR] Failed to delete the original message. Please see CloudWatch.", false)); err != nil {
						log.Printf("[ERROR] Failed to send a message to Slack: %v", err)
					}
					log.Printf("[ERROR] Failed to delete the original message: %v", err)
					return createResponseStatus200(""), nil
				}
				return createResponseStatus200(""), nil
			default:
				log.Printf("[ERROR] Failed to create modal: %v", err)
				return createResponseStatus200(""), nil
			}
		}

		// Send the view to slack
		if _, err := api.OpenView(message.TriggerID, *view); err != nil {
			if _, _, err := api.PostMessage(message.Channel.ID, slack.MsgOptionText("[ERROR] Failed to open modal. Please see CloudWatch.", false)); err != nil {
				log.Printf("[ERROR] Failed to send a message to Slack: %v", err)
			}
			log.Printf("[ERROR] Failed to open modal: err %v", err)
			return createResponseStatus200(""), nil
		}

		// Delete the original message
		if _, _, err := api.DeleteMessage(message.Channel.ID, message.Message.Timestamp); err != nil {
			if _, _, err := api.PostMessage(message.Channel.ID, slack.MsgOptionText("[ERROR] Failed to delete the original message. Please see CloudWatch.", false)); err != nil {
				log.Printf("[ERROR] Failed to send a message to Slack: %v", err)
			}
			log.Printf("[ERROR] Failed to delete the original message: %v", err)
			return createResponseStatus200(""), nil
		}

		return createResponseStatus200(""), nil

	case requestSlackHubToolSelectionCancel:
		// Delete the original message
		if _, _, err := api.DeleteMessage(message.Channel.ID, message.Message.Timestamp); err != nil {
			if _, _, err := api.PostMessage(message.Channel.ID, slack.MsgOptionText("[ERROR] Failed to delete the original message. Please see CloudWatch.", false)); err != nil {
				log.Printf("[ERROR] Failed to send a message to Slack: %v", err)
			}
			log.Printf("[ERROR] Failed to delete the original message: %v", err)
			return createResponseStatus200(""), nil
		}
		return createResponseStatus200(""), nil

	case requestOthers:
		// Get PrivateMeta
		privateMeta := privateMeta{
			ResponseURL: "Unknown",
			ChannelID:   "Unknown",
		}
		json.Unmarshal([]byte(message.View.PrivateMetadata), &privateMeta)

		// Identify tool ID
		id := identifyToolID(message.View.ExternalID)

		// Get selected tool's metadata
		tool := tool.New()
		if err := tool.GetItem(id, region, dynamodb); err != nil {
			log.Printf("[ERROR] Failed to get tool's metadata from DynamoDB: %v", err)
			return createResponseStatus200(""), nil
		}

		// Invoke another lambda
		l := newLambda()
		res, err := l.invokeLambda(tool, message)
		if err != nil {
			if privateMeta.ChannelID != "Unknown" {
				if _, _, err := api.PostMessage(privateMeta.ChannelID, slack.MsgOptionText("[ERROR] Something occured when calling your tool's lambda.", false)); err != nil {
					log.Printf("[ERROR] Failed to send a message to Slack: %v", err)
				}
			}
			log.Printf("[ERROR] Failed to invoke lambda: %v", err)
			return createResponseStatus200(""), err
		}

		// Create a response
		switch tool.BootMode {
		case "Normal":
			// Send success message
			resText := "<@" + message.User.ID + "> invokes \"" + tool.DisplayName + "\" :tophat:"
			if _, _, err := api.PostMessage(privateMeta.ChannelID, slack.MsgOptionText(resText, false)); err != nil {
				return createResponseStatus200(""), err
			}

			// Close the view (NOTE: We can close the view by only sending HTTP 200 response.)
			return createResponseStatus200(""), nil

		case "Advanced":
			var proxyRes events.APIGatewayProxyResponse
			if err != json.Unmarshal(res.Payload, &proxyRes) {
				return createResponseStatus200(""), err
			}
			return proxyRes, nil

		default:
			return createResponseStatus200(""), nil
		}
	}
	return createResponseStatus200(""), nil
}

// identifyRequestType returns the result of determining request type.
func identifyRequestType(message slack.InteractionCallback) requestType {

	// Check if the request comes from modal view or not.
	if actions := message.ActionCallback.BlockActions; len(actions) == 1 && message.View.Hash == "" {
		if actions[0].ActionID == string(requestSlackHubToolSelection) {
			return requestSlackHubToolSelection
		}
		if actions[0].ActionID == string(requestSlackHubToolSelectionCancel) {
			return requestSlackHubToolSelectionCancel
		}
	}
	return requestOthers
}

// identifyToolID
func identifyToolID(externalID string) string {
	return strings.Split(externalID, ":")[0]
}

// invokeLambda invokes another lambda function which is tied to the selected tool.
func (l *lambdaClient) invokeLambda(tool *tool.Tool, message slack.InteractionCallback) (*lambda.InvokeOutput, error) {

	// Create input jsonbytes
	jsonBytes, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}

	// Set invocation type
	var mode string
	switch tool.BootMode {
	case "Normal":
		mode = "Event"
	case "Advanced":
		mode = "RequestResponse"
	default:
		return nil, errUnknownBootModeException
	}

	// Create input data
	input := &lambda.InvokeInput{
		FunctionName:   aws.String(tool.CalleeArn),
		Payload:        jsonBytes,
		InvocationType: aws.String(mode),
	}

	// Invoke a tool lambda
	res, err := l.client.Invoke(input)
	if err != nil {
		return res, err
	}

	return res, nil
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
