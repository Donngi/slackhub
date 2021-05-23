package main

import (
	"github.com/Jimon-s/slackhub/tool"
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
		optObj := slack.NewOptionBlockObject(tool.ID, optText, nil)
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

	return toolList, nil
}
