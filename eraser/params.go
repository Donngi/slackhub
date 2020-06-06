package main

import "errors"

var errInvalidConfirmationException = errors.New("invalid confirmation. confirmation and target tool's display name are different")

const (
	actionIDToolSelection = "eraser_action_tool_selection"
	actionIDConfirmation  = "eraser_action_cofirmation"
	blockIDToolSelection  = "eraser_block_tool_selection"
	blockIDConfirmation   = "eraser_block_cofirmation"
)
