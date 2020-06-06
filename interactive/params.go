package main

import "errors"

type requestType string

const (
	requestSlackHubToolSelection       = requestType("slackhub_tool_selection")
	requestSlackHubToolSelectionCancel = requestType("slackhub_tool_selection_cancel")
	requestOthers                      = requestType("others")
)

var (
	errUnknownBootModeException   = errors.New("unknown boot mode")
	errNonEditableToolsException  = errors.New("no editable tools exist")
	errNonDeletableToolsException = errors.New("no deletable tools exist")
)

const (
	actionIDRegistertoolID          = "register_action_id"
	actionIDRegisterDisplayName     = "register_action_display_name"
	actionIDRegisterDescription     = "register_action_description"
	actionIDRegisterCalleeArn       = "register_action_callee_arn"
	actionIDRegisterAdministrators  = "register_action_administrators"
	actionIDRegisterAuthorizedUsers = "register_action_authorized_users"
	actionIDRegisterModalJSON       = "register_action_modal_json"
	actionIDRegisterBootMode        = "register_action_boot_mode"

	blockIDRegisterToolID          = "register_block_id"
	blockIDRegisterDisplayName     = "register_block_display_name"
	blockIDRegisterDescription     = "register_block_description"
	blockIDRegisterCalleeArn       = "register_block_callee_arn"
	blockIDRegisterAdministrators  = "register_block_administrators"
	blockIDRegisterAuthorizedUsers = "register_block_authorized_users"
	blockIDRegisterModalJSON       = "register_block_modal_json"
	blockIDRegisterReference       = "register_block_reference_context"
	blockIDRegisterBootMode        = "register_block_boot_mode"
	blockIDRegisterAdvanced        = "register_block_advanced"

	actionIDEditorToolSelection = "editor_action_tool_selection"
	blockIDEditorToolSelection  = "editor_block_tool_selection"
	blockIDEditorReference      = "editor_block_reference_context"

	actionIDCatalogToolSelection = "catalog_action_tool_selection"
	blockIDCatalogToolSelection  = "catalog_block_tool_selection"
	blockIDCatalogReference      = "catalog_block_reference_context"

	actionIDEraserToolSelection = "eraser_action_tool_selection"
	actionIDEraserConfirmation  = "eraser_action_cofirmation"
	blockIDEraserToolSelection  = "eraser_block_tool_selection"
	blockIDEraserConfirmation   = "eraser_block_cofirmation"
	blockIDEraserReference      = "eraser_block_reference_context"
)
