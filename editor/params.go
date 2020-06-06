package main

import "errors"

type requestType string

var (
	errInvalidIDException          = errors.New("invalid id. id must not contain a colon (\":\")")
	errInvalidModalJSONException   = errors.New("invalid json. it doesn't meet json format")
	errUnknownRequestTypeException = errors.New("unknown request type")
)

const (
	requestToolSelection = requestType("editor_tool_selection")
	requestChangeRequest = requestType("editor_change_request")
)

const (
	actionIDToolSelection = "editor_action_tool_selection"
	blockIDToolSelection  = "editor_block_tool_selection"

	actionIDtoolID          = "editor_action_id"
	actionIDDisplayName     = "editor_action_display_name"
	actionIDDescription     = "editor_action_description"
	actionIDCalleeArn       = "editor_action_callee_arn"
	actionIDAdministrators  = "editor_action_administrators"
	actionIDAuthorizedUsers = "editor_action_authorized_users"
	actionIDModalJSON       = "editor_action_modal_json"
	actionIDBootMode        = "editor_action_boot_mode"

	blockIDToolID          = "editor_block_id"
	blockIDDisplayName     = "editor_block_display_name"
	blockIDDescription     = "editor_block_description"
	blockIDCalleeArn       = "editor_block_callee_arn"
	blockIDAdministrators  = "editor_block_administrators"
	blockIDAuthorizedUsers = "editor_block_authorized_users"
	blockIDModalJSON       = "editor_block_modal_json"
	blockIDBootMode        = "editor_block_boot_mode"
	blockIDReference       = "editor_block_reference"
	blockIDAdvanced        = "editor_block_advanced"
)
