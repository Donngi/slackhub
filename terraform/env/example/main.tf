data "aws_region" "current" {}

module "dymanodb_tools_table" {
  source = "../../module/dynamodb_tools_table"
}

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
