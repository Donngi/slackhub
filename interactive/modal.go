package main

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/Jimon-s/slackhub/tool"
	"github.com/slack-go/slack"
)

type privateMeta struct {
	ResponseURL string `json:"response_url"`
	ChannelID   string `json:"channel_id"`
}

// createView returns a modal view specified to tool id.
func createModal(tool *tool.Tool, userID string, responseURL string, channelID string) (*slack.ModalViewRequest, error) {

	var modal *slack.ModalViewRequest
	var err error

	switch tool.ID {
	case "register":
		modal = createRegisterModal()
	case "editor":
		modal, err = createEditorModal()
	case "catalog":
		modal, err = createCatalogModal()
	case "eraser":
		modal, err = createEraserModal()
	default:
		modal, err = createStaticModal(tool)
	}
	if err != nil {
		return nil, err
	}

	// NOTE: ExternalID is an invisible field and it must be unique throughout your app.
	// In slackhub, we use this field to identify which tool was selected.
	// The value is {tool.ID}:{UserID}{CurrentTimeStamp} format
	modal.ExternalID = tool.ID + ":" + userID + strconv.FormatInt(time.Now().UTC().UnixNano(), 10)

	// NOTE: PrivateMeta is an invisible field.
	// In slackhub, we use this field to pass convenient parameters to your app.
	// Currently, we pass your app
	// 1. responseURL
	// 2. channelID
	// Both params are useful when you want to send a message to slack.
	params := privateMeta{
		ResponseURL: responseURL,
		ChannelID:   channelID,
	}
	bytes, _ := json.Marshal(params)
	modal.PrivateMetadata = string(bytes)

	return modal, nil
}

// createRegisterModal returns a modal specific to the Register Tool
func createRegisterModal() *slack.ModalViewRequest {

	// Blocks
	descText := slack.NewTextBlockObject("mrkdwn", ":ballot_box_with_ballot: Register is a SlackHub's official tool for registering your new tools.", false, false)
	descTextSection := slack.NewSectionBlock(descText, nil, nil)

	dividerBlock := slack.NewDividerBlock()

	idInputTitle := slack.NewTextBlockObject("plain_text", "Tool ID", false, false)
	idInputElement := slack.NewPlainTextInputBlockElement(nil, actionIDRegistertoolID)
	idInputBlock := slack.NewInputBlock(blockIDRegisterToolID, idInputTitle, idInputElement)
	idInputBlock.Hint = slack.NewTextBlockObject("plain_text", "You should set a unique id to your tool (Can't be changed). NOTE: You can't cantain \":\" character.", false, false)

	displayNameInputTitle := slack.NewTextBlockObject("plain_text", "Display Name", false, false)
	displayNameInputElement := slack.NewPlainTextInputBlockElement(nil, actionIDRegisterDisplayName)
	displayNameInputBlock := slack.NewInputBlock(blockIDRegisterDisplayName, displayNameInputTitle, displayNameInputElement)
	displayNameInputBlock.Hint = slack.NewTextBlockObject("plain_text", "This name will appear in the menu.", false, false)

	descInputTitle := slack.NewTextBlockObject("plain_text", "Description", false, false)
	descInputElement := slack.NewPlainTextInputBlockElement(nil, actionIDRegisterDescription)
	descInputBlock := slack.NewInputBlock(blockIDRegisterDescription, descInputTitle, descInputElement)
	descInputBlock.Hint = slack.NewTextBlockObject("plain_text", "This name will appear when using the catalog tool", false, false)

	calleeArnInputTitle := slack.NewTextBlockObject("plain_text", "Arn of callee lambda", false, false)
	calleeArnInputElement := slack.NewPlainTextInputBlockElement(nil, actionIDRegisterCalleeArn)
	calleeArnInputBlock := slack.NewInputBlock(blockIDRegisterCalleeArn, calleeArnInputTitle, calleeArnInputElement)
	calleeArnInputBlock.Hint = slack.NewTextBlockObject("plain_text", "Set tool lambda's arn. SlackHub can only call lambda.", false, false)

	administratorInputTitle := slack.NewTextBlockObject("plain_text", ":bust_in_silhouette: Administrators", false, false)
	administratorInputElement := slack.NewOptionsMultiSelectBlockElement("multi_users_select", nil, actionIDRegisterAdministrators)
	administratorInputBlock := slack.NewInputBlock(blockIDRegisterAdministrators, administratorInputTitle, administratorInputElement)
	administratorInputBlock.Optional = true
	administratorInputBlock.Hint = slack.NewTextBlockObject("plain_text", "Optional: If you select administrators, no one else will be able to edit the tool. Default is none.", false, false)

	authorizedInputTitle := slack.NewTextBlockObject("plain_text", ":busts_in_silhouette: Authorized users", false, false)
	authorizedInputElement := slack.NewOptionsMultiSelectBlockElement("multi_users_select", nil, actionIDRegisterAuthorizedUsers)
	authorizedInputBlock := slack.NewInputBlock(blockIDRegisterAuthorizedUsers, authorizedInputTitle, authorizedInputElement)
	authorizedInputBlock.Optional = true
	authorizedInputBlock.Hint = slack.NewTextBlockObject("plain_text", "Optional: If you select authorized users, no one else will be able to use the tool. Default is none.", false, false)

	modalJSONInputTitle := slack.NewTextBlockObject("plain_text", "Modal JSON", false, false)
	modalJSONInputElement := slack.NewPlainTextInputBlockElement(nil, actionIDRegisterModalJSON)
	modalJSONInputElement.Multiline = true
	modalJSONInputBlock := slack.NewInputBlock(blockIDRegisterModalJSON, modalJSONInputTitle, modalJSONInputElement)
	modalJSONInputBlock.Optional = true
	modalJSONInputBlock.Hint = slack.NewTextBlockObject("plain_text", "Optional: Source of a modal view appearance. If you don't register any modal, SlackHub won't send modal to Slack and will simply pass the request to your tool's lambda.", false, false)

	referenceText := slack.NewTextBlockObject("mrkdwn", "You can make json with *Official GUI Tool*. For more info, see <https://github.com/Jimon-s/slackhub/blob/master/documents/guide_for_developer/step1_create_modal_view|SlackHub user guide>.", false, false)
	referenceContextBlock := slack.NewContextBlock(blockIDRegisterReference, []slack.MixedElement{referenceText}...)

	advText := slack.NewTextBlockObject("mrkdwn", ":zap: *ADVANCED SETTING*", false, false)
	advTextSection := slack.NewSectionBlock(advText, nil, nil)

	advContextText := slack.NewTextBlockObject("mrkdwn", "If you learn Slack's modal lifecycle sequence, you can create more flexible tools. By changing this param, you can fully control your event. For more info, see <https://github.com/Jimon-s/slackhub/blob/master/documents/guide_for_developer/stepx_use_advanced_mode|SlackHub user advanced guide>.", false, false)
	advContextBlock := slack.NewContextBlock(blockIDRegisterAdvanced, []slack.MixedElement{advContextText}...)

	optNormalText := slack.NewTextBlockObject("plain_text", "Normal", false, false)
	optNormalObj := slack.NewOptionBlockObject("Normal", optNormalText, nil)
	optAdvancedText := slack.NewTextBlockObject("plain_text", "Advanced", false, false)
	optAdvancedObj := slack.NewOptionBlockObject("Advanced", optAdvancedText, nil)
	optObjs := []*slack.OptionBlockObject{
		optNormalObj,
		optAdvancedObj,
	}
	bootInputText := slack.NewTextBlockObject("plain_text", "Normal", false, false)
	bootInputElement := slack.NewOptionsSelectBlockElement("static_select", bootInputText, actionIDRegisterBootMode, optObjs...)
	bootInputElement.InitialOption = optNormalObj
	bootInputTitle := slack.NewTextBlockObject("plain_text", "Boot Mode", false, false)
	bootInputBlock := slack.NewInputBlock(blockIDRegisterBootMode, bootInputTitle, bootInputElement)
	bootInputBlock.Hint = slack.NewTextBlockObject("plain_text", "Default value is \"Normal\" and we strictly recommend you to keep \"Normal\" if you are NOT familiar with Slack's lifecycle sequence. When you set \"Advanced\", we will call your lambda synchronously, and you must manage all request/response by yourself.", false, false)

	blocks := slack.Blocks{
		BlockSet: []slack.Block{
			descTextSection,
			dividerBlock,
			idInputBlock,
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
		Title:  slack.NewTextBlockObject("plain_text", "SlackHub - Register", false, false),
		Close:  slack.NewTextBlockObject("plain_text", "Cancel", false, false),
		Submit: slack.NewTextBlockObject("plain_text", "Submit", false, false),
		Blocks: blocks,
	}

	return &modal
}

// createEditorModal returns a dialog specific to the Editor Tool
func createEditorModal() (*slack.ModalViewRequest, error) {

	// Get all tool's metadata from DynamoDB.
	toolList, err := tool.GetAllItems(region, dynamodb)
	if err != nil {
		return nil, err
	}

	// Sort tool list
	toolList = tool.SortTools(toolList)

	// Text Section
	descText := slack.NewTextBlockObject("mrkdwn", ":wrench: Editor is a SlackHub's official tool for editing your tools. ", false, false)
	descTextSection := slack.NewSectionBlock(descText, nil, nil)

	dividerBlock := slack.NewDividerBlock()

	// Static select
	var optObjs []*slack.OptionBlockObject
	for _, tool := range toolList {

		// remove slackhub official tools
		if tool.ID == "register" || tool.ID == "editor" || tool.ID == "catalog" || tool.ID == "eraser" {
			continue
		}

		optText := slack.NewTextBlockObject("plain_text", tool.DisplayName, false, false)
		optObj := slack.NewOptionBlockObject(tool.ID, optText, nil)
		optObjs = append(optObjs, optObj)
	}

	// Check if there are some tools to be edited
	if len(optObjs) == 0 {
		return nil, errNonEditableToolsException
	}

	selectInputText := slack.NewTextBlockObject("plain_text", "Tools ...", false, false)
	selectInputElement := slack.NewOptionsSelectBlockElement("static_select", selectInputText, actionIDEditorToolSelection, optObjs...)
	selectInputTitle := slack.NewTextBlockObject("plain_text", "Select the tool to be edited", false, false)
	selectInputBlock := slack.NewInputBlock(blockIDEditorToolSelection, selectInputTitle, selectInputElement)

	// Reference
	referenceText := slack.NewTextBlockObject("mrkdwn", "NOTE: You can't edit SlackHub's official tools (Register, Editor, Catalog, Eraser).\nNOTE: You can't edit ID of tools.", false, false)
	referenceContextBlock := slack.NewContextBlock(blockIDEditorReference, []slack.MixedElement{referenceText}...)

	blocks := slack.Blocks{
		BlockSet: []slack.Block{
			descTextSection,
			dividerBlock,
			selectInputBlock,
			referenceContextBlock,
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

	// Set Callback ID to recognize request type of editor (This is a specific approach of editor)
	modal.CallbackID = "editor_tool_selection"

	return &modal, nil
}

// createCatalogModal returns a dialog specific to the Editor Tool
func createCatalogModal() (*slack.ModalViewRequest, error) {
	// Get all tool's metadata from DynamoDB.
	toolList, err := tool.GetAllItems(region, dynamodb)
	if err != nil {
		return nil, err
	}

	// Sort tool list to show slackhub's official tool at bottom.
	toolList = tool.SortTools(toolList)

	// Text Section
	descText := slack.NewTextBlockObject("mrkdwn", ":green_book: Catalog is a SlackHub's official tool for listing your tools.", false, false)
	descTextSection := slack.NewSectionBlock(descText, nil, nil)

	dividerBlock := slack.NewDividerBlock()

	// Static select
	var optObjs []*slack.OptionBlockObject
	allText := slack.NewTextBlockObject("plain_text", "Please show all tools!", false, false)
	allObj := slack.NewOptionBlockObject("all", allText, nil)
	optObjs = append(optObjs, allObj)

	for _, tool := range toolList {
		optText := slack.NewTextBlockObject("plain_text", tool.DisplayName, false, false)
		optObj := slack.NewOptionBlockObject(tool.ID, optText, nil)
		optObjs = append(optObjs, optObj)
	}

	selectInputElement := slack.NewOptionsSelectBlockElement("static_select", nil, actionIDCatalogToolSelection, optObjs...)
	selectInputElement.InitialOption = allObj
	selectInputTitle := slack.NewTextBlockObject("plain_text", "Select the tool you want to know more", false, false)
	selectInputBlock := slack.NewInputBlock(blockIDCatalogToolSelection, selectInputTitle, selectInputElement)

	// Reference
	referenceText := slack.NewTextBlockObject("mrkdwn", "NOTE: If you have a lot of tools, the catalog will be too large size, and your channel will be got buried under the catalog. We recommend selecting a tool or delete a catalog after you read.", false, false)
	referenceContextBlock := slack.NewContextBlock(blockIDCatalogReference, []slack.MixedElement{referenceText}...)

	blocks := slack.Blocks{
		BlockSet: []slack.Block{
			descTextSection,
			dividerBlock,
			selectInputBlock,
			referenceContextBlock,
		},
	}

	// ModalView
	modal := slack.ModalViewRequest{
		Type:   slack.ViewType("modal"),
		Title:  slack.NewTextBlockObject("plain_text", "SlackHub - Catalog", false, false),
		Close:  slack.NewTextBlockObject("plain_text", "Cancel", false, false),
		Submit: slack.NewTextBlockObject("plain_text", "Submit", false, false),
		Blocks: blocks,
	}

	return &modal, nil
}

// createEraserModal returns a dialog specific to the Editor Tool
func createEraserModal() (*slack.ModalViewRequest, error) {

	// Get all tool's metadata from DynamoDB.
	toolList, err := tool.GetAllItems(region, dynamodb)
	if err != nil {
		return nil, err
	}

	// Text Section
	descText := slack.NewTextBlockObject("mrkdwn", ":warning: Eraser is a SlackHub's official tool for deleting your tools.", false, false)
	descTextSection := slack.NewSectionBlock(descText, nil, nil)

	dividerBlock := slack.NewDividerBlock()

	// Static select
	var optObjs []*slack.OptionBlockObject
	for _, tool := range toolList {

		// remove slackhub official tools
		if tool.ID == "register" || tool.ID == "editor" || tool.ID == "catalog" || tool.ID == "eraser" {
			continue
		}

		optText := slack.NewTextBlockObject("plain_text", tool.DisplayName, false, false)
		optObj := slack.NewOptionBlockObject(tool.ID, optText, nil)
		optObjs = append(optObjs, optObj)
	}

	// Check if there are some tools to be deleted
	if len(optObjs) == 0 {
		return nil, errNonDeletableToolsException
	}

	selectInputText := slack.NewTextBlockObject("plain_text", "Tools ...", false, false)
	selectInputElement := slack.NewOptionsSelectBlockElement("static_select", selectInputText, actionIDEraserToolSelection, optObjs...)
	selectInputTitle := slack.NewTextBlockObject("plain_text", "Select the tool to be deleted", false, false)
	selectInputBlock := slack.NewInputBlock(blockIDEraserToolSelection, selectInputTitle, selectInputElement)

	// Reference
	referenceText := slack.NewTextBlockObject("mrkdwn", "NOTE: You can't delete SlackHub's official tools (Register, Editor, Catalog, Eraser).", false, false)
	referenceContextBlock := slack.NewContextBlock(blockIDEraserReference, []slack.MixedElement{referenceText}...)

	// Confirmation
	confirmationInputTitle := slack.NewTextBlockObject("plain_text", "Enter the \"Display Name\" of the tool to confirm. (like \"Eraser\")", false, false)
	confirmationInputElement := slack.NewPlainTextInputBlockElement(nil, actionIDEraserConfirmation)
	confirmationInputBlock := slack.NewInputBlock(blockIDEraserConfirmation, confirmationInputTitle, confirmationInputElement)

	blocks := slack.Blocks{
		BlockSet: []slack.Block{
			descTextSection,
			dividerBlock,
			selectInputBlock,
			referenceContextBlock,
			confirmationInputBlock,
		},
	}

	// ModalView
	modal := slack.ModalViewRequest{
		Type:   slack.ViewType("modal"),
		Title:  slack.NewTextBlockObject("plain_text", "SlackHub - Eraser", false, false),
		Close:  slack.NewTextBlockObject("plain_text", "Cancel", false, false),
		Submit: slack.NewTextBlockObject("plain_text", "Submit", false, false),
		Blocks: blocks,
	}

	return &modal, nil
}

// createStaticModal returns a dialog which is based on selected tool's ModalJSON value.
func createStaticModal(tool *tool.Tool) (*slack.ModalViewRequest, error) {

	// Create a modal view
	var modal slack.ModalViewRequest
	if err := json.Unmarshal([]byte(tool.ModalJSON), &modal); err != nil {
		return nil, err
	}

	return &modal, nil
}
