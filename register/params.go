package main

import "errors"

var (
	errInvalidIDException        = errors.New("invalid id. id must not contain a colon (\":\")")
	errInvalidModalJSONException = errors.New("invalid json. it doesn't meet json format")
)

const (
	actionIDtoolID          = "register_action_id"
	actionIDDisplayName     = "register_action_display_name"
	actionIDDescription     = "register_action_description"
	actionIDCalleeArn       = "register_action_callee_arn"
	actionIDAdministrators  = "register_action_administrators"
	actionIDAuthorizedUsers = "register_action_authorized_users"
	actionIDModalJSON       = "register_action_modal_json"
	actionIDBootMode        = "register_action_boot_mode"

	blockIDToolID          = "register_block_id"
	blockIDDisplayName     = "register_block_display_name"
	blockIDDescription     = "register_block_description"
	blockIDCalleeArn       = "register_block_callee_arn"
	blockIDAdministrators  = "register_block_administrators"
	blockIDAuthorizedUsers = "register_block_authorized_users"
	blockIDModalJSON       = "register_block_modal_json"
	blockIDBootMode        = "register_block_boot_mode"
)
