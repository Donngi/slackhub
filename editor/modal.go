package main

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/nicoJN/chatops-slackhub/tool"
	"github.com/slack-go/slack"
)

// createModal returns a modal view specified to the Editor tool.
func createModal(tool *tool.Tool, userID string, responseURL string, channelID string) *slack.ModalViewRequest {

	// Blocks
	descText := slack.NewTextBlockObject("mrkdwn", ":wrench: Let's edit your tool!", false, false)
	descTextSection := slack.NewSectionBlock(descText, nil, nil)

	dividerBlock := slack.NewDividerBlock()

	displayNameInputTitle := slack.NewTextBlockObject("plain_text", "Display Name", false, false)
	displayNameInputElement := slack.NewPlainTextInputBlockElement(nil, actionIDDisplayName)
	displayNameInputElement.InitialValue = tool.DisplayName
	displayNameInputBlock := slack.NewInputBlock(blockIDDisplayName, displayNameInputTitle, displayNameInputElement)
	displayNameInputBlock.Hint = slack.NewTextBlockObject("plain_text", "This name will appear in the menu.", false, false)

	descInputTitle := slack.NewTextBlockObject("plain_text", "Description", false, false)
	descInputElement := slack.NewPlainTextInputBlockElement(nil, actionIDDescription)
	descInputElement.InitialValue = tool.Description
	descInputBlock := slack.NewInputBlock(blockIDDescription, descInputTitle, descInputElement)
	descInputBlock.Hint = slack.NewTextBlockObject("plain_text", "This name will appear when using the catalog tool.", false, false)

	calleeArnInputTitle := slack.NewTextBlockObject("plain_text", "Arn of callee lambda", false, false)
	calleeArnInputElement := slack.NewPlainTextInputBlockElement(nil, actionIDCalleeArn)
	calleeArnInputElement.InitialValue = tool.CalleeArn
	calleeArnInputBlock := slack.NewInputBlock(blockIDCalleeArn, calleeArnInputTitle, calleeArnInputElement)
	calleeArnInputBlock.Hint = slack.NewTextBlockObject("plain_text", "Set tool lambda's arn. SlackHub can only call lambda.", false, false)

	administratorInputTitle := slack.NewTextBlockObject("plain_text", ":bust_in_silhouette: Administrators", false, false)
	administratorInputElement := slack.NewOptionsMultiSelectBlockElement("multi_users_select", nil, actionIDAdministrators)
	administratorInputElement.InitialUsers = tool.Administrators
	administratorInputBlock := slack.NewInputBlock(blockIDAdministrators, administratorInputTitle, administratorInputElement)
	administratorInputBlock.Optional = true
	administratorInputBlock.Hint = slack.NewTextBlockObject("plain_text", "Optional: If you select administrators, no one else will be able to edit the tool. Default is none.", false, false)

	authorizedInputTitle := slack.NewTextBlockObject("plain_text", ":busts_in_silhouette: Authorized users", false, false)
	authorizedInputElement := slack.NewOptionsMultiSelectBlockElement("multi_users_select", nil, actionIDAuthorizedUsers)
	authorizedInputElement.InitialUsers = tool.AuthorizedUsers
	authorizedInputBlock := slack.NewInputBlock(blockIDAuthorizedUsers, authorizedInputTitle, authorizedInputElement)
	authorizedInputBlock.Optional = true
	authorizedInputBlock.Hint = slack.NewTextBlockObject("plain_text", "Optional: If you select authorized users, no one else will be able to use the tool. Default is none.", false, false)

	modalJSONInputTitle := slack.NewTextBlockObject("plain_text", "Modal JSON", false, false)
	modalJSONInputElement := slack.NewPlainTextInputBlockElement(nil, actionIDModalJSON)
	modalJSONInputElement.InitialValue = tool.ModalJSON
	modalJSONInputElement.Multiline = true
	modalJSONInputBlock := slack.NewInputBlock(blockIDModalJSON, modalJSONInputTitle, modalJSONInputElement)
	modalJSONInputBlock.Optional = true
	modalJSONInputBlock.Hint = slack.NewTextBlockObject("plain_text", "Optional: Source of a modal view appearance. If you don't register any modal, SlackHub won't send modal to Slack and will simply pass the request to your tool's lambda.", false, false)

	referenceText := slack.NewTextBlockObject("mrkdwn", "You can make json with *Official GUI Tool*. For more info, see <https://github.com/nicoJN/slackhub/blob/master/documents/guide_for_developer/step1_create_modal_view|SlackHub user guide>.", false, false)
	referenceContextBlock := slack.NewContextBlock(blockIDReference, []slack.MixedElement{referenceText}...)

	advText := slack.NewTextBlockObject("mrkdwn", ":zap: *ADVANCED SETTING*", false, false)
	advTextSection := slack.NewSectionBlock(advText, nil, nil)

	advContextText := slack.NewTextBlockObject("mrkdwn", "If you learn Slack's modal lifecycle sequence, you can create more flexible tools. By changing this param, you can fully control your event. For more info, see <https://github.com/nicoJN/slackhub/blob/master/documents/guide_for_developer/stepx_use_advanced_mode|SlackHub user advanced guide>.", false, false)
	advContextBlock := slack.NewContextBlock(blockIDAdvanced, []slack.MixedElement{advContextText}...)

	optNormalText := slack.NewTextBlockObject("plain_text", "Normal", false, false)
	optNormalObj := slack.NewOptionBlockObject("Normal", optNormalText)
	optAdvancedText := slack.NewTextBlockObject("plain_text", "Advanced", false, false)
	optAdvancedObj := slack.NewOptionBlockObject("Advanced", optAdvancedText)
	optObjs := []*slack.OptionBlockObject{
		optNormalObj,
		optAdvancedObj,
	}
	bootInputText := slack.NewTextBlockObject("plain_text", "Normal", false, false)
	bootInputElement := slack.NewOptionsSelectBlockElement("static_select", bootInputText, actionIDBootMode, optObjs...)
	switch tool.BootMode {
	case "Normal":
		bootInputElement.InitialOption = optNormalObj
	case "Advanced":
		bootInputElement.InitialOption = optAdvancedObj
	}
	bootInputTitle := slack.NewTextBlockObject("plain_text", "Boot Mode", false, false)
	bootInputBlock := slack.NewInputBlock(blockIDBootMode, bootInputTitle, bootInputElement)
	bootInputBlock.Hint = slack.NewTextBlockObject("plain_text", "Default value is \"Normal\" and we strictly recommend you to keep \"Normal\" if you are NOT familiar with Slack's lifecycle sequence. When you set \"Advanced\", we will call your lambda synchronously, and you must manage all request/response by yourself.", false, false)

	blocks := slack.Blocks{
		BlockSet: []slack.Block{
			descTextSection,
			dividerBlock,
			displayNameInputBlock,
			descInputBlock,
			calleeArnInputBlock,
			administratorInputBlock,
			authorizedInputBlock,
			modalJSONInputBlock,
			referenceContextBlock,
			dividerBlock,
			advTextSection,
			advContextBlock,
			bootInputBlock,
		},
	}

	// ModalView
	modal := slack.ModalViewRequest{
		Type:   slack.ViewType("modal"),
		Title:  slack.NewTextBlockObject("plain_text", "SlackHub - Editor", false, false),
		Close:  slack.NewTextBlockObject("plain_text", "Cancel", false, false),
		Submit: slack.NewTextBlockObject("plain_text", "Submit", false, false),
		Blocks: blocks,
	}

	// External ID
	modal.ExternalID = "editor" + ":" + userID + strconv.FormatInt(time.Now().UTC().UnixNano(), 10)

	// Callback ID
	modal.CallbackID = string(requestChangeRequest)

	// PrivateMera
	params := privateMeta{
		ResponseURL:  responseURL,
		ChannelID:    channelID,
		TargetToolID: tool.ID,
	}
	bytes, _ := json.Marshal(params)
	modal.PrivateMetadata = string(bytes)

	return &modal
}
