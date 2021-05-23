data "aws_region" "current" {}

# ----------------------------------------------------------
# SlackHub main
# ----------------------------------------------------------

module "dymanodb_tools_table" {
  source              = "../../module/dynamodb_tools_table"
  lambda_register_arn = module.lambda_register.lambda_arn
  lambda_editor_arn   = module.lambda_editor.lambda_arn
  lambda_catalog_arn  = module.lambda_catalog.lambda_arn
  lambda_eraser_arn   = module.lambda_eraser.lambda_arn
}

module "api_gateway" {
  source                        = "../../module/api_gateway"
  lambda_initial_invoke_arn     = module.lambda_initial.lambda_invoke_arn
  lambda_interactive_invoke_arn = module.lambda_interactive.lambda_invoke_arn
}

module "lambda_initial" {
  source = "../../module/lambda_slackhub_main"

  function_name                      = "SlackHubInitial"
  source_code_dir                    = "../../../initial/bin"
  source_code_file                   = "main"
  region                             = data.aws_region.current.name
  param_key_bot_user_auth_token      = local.param_key_bot_user_auth_token
  param_key_signing_secret           = local.param_key_signing_secret
  dynamodb_table_name                = module.dymanodb_tools_table.dynamodb_table_name
  dynamodb_table_arn                 = module.dymanodb_tools_table.dynamodb_table_arn
  api_gateway_slackhub_execution_arn = module.api_gateway.api_gateway_slackhub_execution_arn
}

module "lambda_interactive" {
  source = "../../module/lambda_slackhub_main"

  function_name                      = "SlackHubInteractive"
  source_code_dir                    = "../../../interactive/bin"
  source_code_file                   = "main"
  region                             = data.aws_region.current.name
  param_key_bot_user_auth_token      = local.param_key_bot_user_auth_token
  param_key_signing_secret           = local.param_key_signing_secret
  dynamodb_table_name                = module.dymanodb_tools_table.dynamodb_table_name
  dynamodb_table_arn                 = module.dymanodb_tools_table.dynamodb_table_arn
  api_gateway_slackhub_execution_arn = module.api_gateway.api_gateway_slackhub_execution_arn
}

# ----------------------------------------------------------
# Official tools
# ----------------------------------------------------------

module "lambda_register" {
  source = "../../module/lambda_official_tool"

  function_name                 = "SlackHubRegister"
  source_code_dir               = "../../../register/bin"
  source_code_file              = "main"
  region                        = data.aws_region.current.name
  param_key_bot_user_auth_token = local.param_key_bot_user_auth_token
  param_key_signing_secret      = local.param_key_signing_secret
  dynamodb_table_name           = module.dymanodb_tools_table.dynamodb_table_name
  dynamodb_table_arn            = module.dymanodb_tools_table.dynamodb_table_arn
}

module "lambda_editor" {
  source = "../../module/lambda_official_tool"

  function_name                 = "SlackHubEditor"
  source_code_dir               = "../../../editor/bin"
  source_code_file              = "main"
  region                        = data.aws_region.current.name
  param_key_bot_user_auth_token = local.param_key_bot_user_auth_token
  param_key_signing_secret      = local.param_key_signing_secret
  dynamodb_table_name           = module.dymanodb_tools_table.dynamodb_table_name
  dynamodb_table_arn            = module.dymanodb_tools_table.dynamodb_table_arn
}

module "lambda_catalog" {
  source = "../../module/lambda_official_tool"

  function_name                 = "SlackHubCatalog"
  source_code_dir               = "../../../catalog/bin"
  source_code_file              = "main"
  region                        = data.aws_region.current.name
  param_key_bot_user_auth_token = local.param_key_bot_user_auth_token
  param_key_signing_secret      = local.param_key_signing_secret
  dynamodb_table_name           = module.dymanodb_tools_table.dynamodb_table_name
  dynamodb_table_arn            = module.dymanodb_tools_table.dynamodb_table_arn
}

module "lambda_eraser" {
  source = "../../module/lambda_official_tool"

  function_name                 = "SlackHubEraser"
  source_code_dir               = "../../../eraser/bin"
  source_code_file              = "main"
  region                        = data.aws_region.current.name
  param_key_bot_user_auth_token = local.param_key_bot_user_auth_token
  param_key_signing_secret      = local.param_key_signing_secret
  dynamodb_table_name           = module.dymanodb_tools_table.dynamodb_table_name
  dynamodb_table_arn            = module.dymanodb_tools_table.dynamodb_table_arn
}
