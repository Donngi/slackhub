package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
	"github.com/nicoJN/slackhub/tool"
	"github.com/slack-go/slack"
)

// createMenu returns menu of tools.
func createMenu() (slack.MsgOption, error) {
	// Section - Text
	headerText := slack.NewTextBlockObject("mrkdwn", "*Welcome!*\nLet's select the tool:+1:", false, false)
	headerTextSection := slack.NewSectionBlock(headerText, nil, nil)

	// Action - Cancel button
	cancelButtonText := slack.NewTextBlockObject("plain_text", "Cancel", true, false)
	cancelButtonElement := slack.NewButtonBlockElement("slackhub_tool_selection_cancel", "ACTION_CANCEL", cancelButtonText)

	// Action - Select
	selectElement, err := createSelectElement()
	if err != nil {
		return nil, err
	}

	// Actions - Combine cancel button and select
	actions := slack.NewActionBlock("", selectElement, cancelButtonElement)

	// Blocks
	blocks := slack.MsgOptionBlocks(headerTextSection, actions)

	return blocks, nil
}

// createSelectElement returns a select block element which contains a tool list.
func createSelectElement() (*slack.SelectBlockElement, error) {

	toolList, err := getToolList()
	if err != nil {
		return nil, err
	}

	var optObjs []*slack.OptionBlockObject
	for _, tool := range toolList {
		optText := slack.NewTextBlockObject("plain_text", tool.DisplayName, false, false)
		optObj := slack.NewOptionBlockObject(tool.ID, optText)
		optObjs = append(optObjs, optObj)
	}

	selectText := slack.NewTextBlockObject("plain_text", "Select", false, false)
	selectElement := slack.NewOptionsSelectBlockElement("static_select", selectText, "slackhub_tool_selection", optObjs...)
	return selectElement, nil
}

// getToolList loads tool list from DynamoDB.
// NOTE: The number of tools doesn't increase so much, so this function uses DynamoDB full-scan.
func getToolList() ([]tool.Tool, error) {

	// Get all tool's metadata from DynamoDB.
	toolList, err := tool.GetAllItems(region, dynamodb)
	if err != nil {
		return toolList, err
	}

	// Sort tool list to show slackhub's official tool at bottom.
	toolList = tool.SortTools(toolList)

	// If dynamodb doesn't have any records, initialize it with register tool.
	if len(toolList) == 0 {

		toolList, err = initializeToolList()
		if err != nil {
			return toolList, err
		}
	}

	return toolList, nil
}

// initializeToolList puts an initial record to DynamoDB.
func initializeToolList() ([]tool.Tool, error) {
	var db = dynamo.New(session.New(), &aws.Config{
		Region: aws.String(region),
	})
	var table = db.Table(dynamodb)

	registerTool := tool.Tool{
		ID:          "register",
		DisplayName: "SlackHub - Register",
		Description: "Register is a SlackHub's official tool for registering your new tools.",
		CalleeArn:   registerArn,
		ModalJSON:   "SlackHub's official tool Register dinamically creates modal json.",
		BootMode:    "Normal",
	}

	editorTool := tool.Tool{
		ID:          "editor",
		DisplayName: "SlackHub - Editor",
		Description: "Editor is a SlackHub's official tool for editing your tools.",
		CalleeArn:   editorArn,
		ModalJSON:   "SlackHub's official tool Editor dinamically creates modal json.",
		BootMode:    "Advanced",
	}

	catalogTool := tool.Tool{
		ID:          "catalog",
		DisplayName: "SlackHub - Catalog",
		Description: "Catalog is a SlackHub's official tool for listing your tools.",
		CalleeArn:   catalogArn,
		ModalJSON:   "SlackHub's official tool Catalog dinamically creates modal json.",
		BootMode:    "Normal",
	}

	eraserTool := tool.Tool{
		ID:          "eraser",
		DisplayName: "SlackHub - Eraser",
		Description: "Eraser is a SlackHub's official tool for deleting your tools.",
		CalleeArn:   eraserArn,
		ModalJSON:   "SlackHub's official tool Eraser dinamically creates modal json.",
		BootMode:    "Normal",
	}

	sampleGoTool := tool.Tool{
		ID:          "sample_go",
		DisplayName: "Sample Tool",
		Description: "Sample tool written in Go. Simply returns user's inputs.",
		CalleeArn:   sampleGoArn,
		ModalJSON:   sampleGoJSON,
		BootMode:    "Normal",
	}

	if err := table.Put(registerTool).Run(); err != nil {
		return nil, err
	}

	if err := table.Put(editorTool).Run(); err != nil {
		return nil, err
	}

	if err := table.Put(catalogTool).Run(); err != nil {
		return nil, err
	}

	if err := table.Put(eraserTool).Run(); err != nil {
		return nil, err
	}

	if err := table.Put(sampleGoTool).Run(); err != nil {
		return nil, err
	}

	return []tool.Tool{sampleGoTool, registerTool, editorTool, catalogTool, eraserTool}, nil
}
