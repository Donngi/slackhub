package main

import (
	"encoding/json"
	"log"
	"os"
	"strconv"

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
	var toolList []tool.Tool
	id := message.View.State.Values[blockIDToolSelection][actionIDToolSelection].SelectedOption.Value
	switch id {
	case "all":
		// Get all tool's metadata from DynamoDB.
		toolList, err = tool.GetAllItems(region, dynamodb)
		if err != nil {
			if _, _, err := api.PostMessage(privateMeta.ChannelID, slack.MsgOptionText("[ERROR] Failed to get tool's metadata from DynamoDB", false)); err != nil {
				log.Printf("[ERROR] Failed to send a message to Slack: %v", err)
			}
			log.Printf("Failed to get tool's metadata from DynamoDB")
			return
		}

		// Sort
		toolList = tool.SortTools(toolList)

	default:
		tool := tool.New()
		if err := tool.GetItem(id, region, dynamodb); err != nil {
			if _, _, err := api.PostMessage(privateMeta.ChannelID, slack.MsgOptionText("[ERROR] Failed to get tool's metadata from DynamoDB", false)); err != nil {
				log.Printf("[ERROR] Failed to send a message to Slack: %v", err)
			}
			log.Printf("Failed to get tool's metadata from DynamoDB")
			return
		}
		toolList = append(toolList, *tool)
	}

	// Send catalog
	if err = sendCatalog(toolList, api, privateMeta.ChannelID, message.User.ID); err != nil {
		if _, _, err := api.PostMessage(privateMeta.ChannelID, slack.MsgOptionText("[ERROR] Failed to send a catalog to slack", false)); err != nil {
			log.Printf("[ERROR] Failed to send a message to Slack: %v", err)
		}
		log.Printf("[ERROR] Failed to send a catalog to slack")
	}

}

// sendCatalog sends some catalogs of tools to slack.
// Each message contains up to 10 tools information
func sendCatalog(toolList []tool.Tool, api *slack.Client, channelID string, userID string) error {
	// Divide toolList
	divided := divideToolList(toolList, 10)

	// Send catalog
	for i, list := range divided {
		option, err := createOption(list, api, i == 0 && len(toolList) != 1)
		if err != nil {
			return err
		}

		// Send a catalog to slack
		if _, err := api.PostEphemeral(channelID, userID, option); err != nil {
			return err
		}
	}

	return nil
}

// divideToolList divide slice into slice of slice.
func divideToolList(toolList []tool.Tool, num int) [][]tool.Tool {
	var divided [][]tool.Tool
	for i := 0; i < len(toolList); i += num {
		end := i + num
		if len(toolList) < end {
			end = len(toolList)
		}
		divided = append(divided, toolList[i:end])
	}
	return divided
}

// createOption returns a catalog option.
func createOption(toolList []tool.Tool, api *slack.Client, withHeader bool) (slack.MsgOption, error) {

	var buf []slack.Block

	// Divider
	dividerBlock := slack.NewDividerBlock()

	if withHeader {
		// Section - Text
		headerText := slack.NewTextBlockObject("mrkdwn", "*Hi!*\nYou have "+strconv.Itoa(len(toolList))+" tools! :tada:", false, false)
		headerTextSection := slack.NewSectionBlock(headerText, nil, nil)

		// Store blocks
		buf = append(buf, headerTextSection, dividerBlock)
	}

	for _, tool := range toolList {
		// Display Name, Tool ID
		displayNameField := slack.NewTextBlockObject("mrkdwn", "*Display Name "+getIcon(tool.ID)+":*\n"+tool.DisplayName, false, false)
		idField := slack.NewTextBlockObject("mrkdwn", "*Tood ID:*\n"+tool.ID, false, false)

		var topFieldSlice []*slack.TextBlockObject
		topFieldSlice = append(topFieldSlice, displayNameField)
		topFieldSlice = append(topFieldSlice, idField)

		topFieldsSection := slack.NewSectionBlock(nil, topFieldSlice, nil)

		// Description
		descText := slack.NewTextBlockObject("mrkdwn", "*Description:*\n"+tool.Description, false, false)
		descTextSection := slack.NewSectionBlock(descText, nil, nil)

		// Administrators, Authorized Users
		admins, err := getNameList(tool.Administrators, api)
		if err != nil {
			return nil, err
		}
		adminField := slack.NewTextBlockObject("mrkdwn", "*Administrator:*\n"+sliceToString(admins, " "), false, false)
		authorizedUsers, err := getNameList(tool.AuthorizedUsers, api)
		if err != nil {
			return nil, err
		}
		authorizedUsersField := slack.NewTextBlockObject("mrkdwn", "*Authorized Users:*\n"+sliceToString(authorizedUsers, " "), false, false)

		var bottomFieldSlice []*slack.TextBlockObject
		bottomFieldSlice = append(bottomFieldSlice, adminField)
		bottomFieldSlice = append(bottomFieldSlice, authorizedUsersField)

		bottomFieldsSection := slack.NewSectionBlock(nil, bottomFieldSlice, nil)

		// Store
		buf = append(buf, topFieldsSection, descTextSection, bottomFieldsSection, dividerBlock)

	}

	// Blocks
	blocks := slack.MsgOptionBlocks(buf...)

	return blocks, nil
}

// getIcon returns a emoji specified to the tool.
func getIcon(id string) string {
	switch id {
	case "register":
		return ":ballot_box_with_ballot:"
	case "editor":
		return ":wrench:"
	case "catalog":
		return ":green_book:"
	case "eraser":
		return ":warning:"
	default:
		return ":name_badge:"
	}
}

// getNameList returns display name of slack users. If users are empty, "All users" is set instead.
func getNameList(users []string, api *slack.Client) ([]string, error) {
	var res []string
	if len(users) != 0 {
		admins, err := api.GetUsersInfo(users...)
		if err != nil {
			return nil, err
		}

		for _, user := range *admins {
			res = append(res, user.Name)
		}

	} else {
		res = append(res, "All users")
	}
	return res, nil
}

// sliceToString combine elements in the slice with separator.
func sliceToString(s []string, separator string) string {
	var res string
	for i, v := range s {
		res += v
		if i != len(s)-1 {
			res += separator
		}
	}
	return res
}
