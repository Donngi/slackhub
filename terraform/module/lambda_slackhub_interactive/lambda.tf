data "archive_file" "zip" {
  type        = "zip"
  source_file = "${var.source_code_dir}/${var.source_code_file}"
  output_path = "${var.source_code_dir}/upload_terraform.zip"
}

resource "aws_lambda_function" "interactive" {
  filename      = "${var.source_code_dir}/upload_terraform.zip"
  function_name = "SlackHubInteractive"
  role          = aws_iam_role.interactive.arn
  handler       = "main"

  runtime = "go1.x"

  source_code_hash = data.archive_file.zip.output_base64sha256

  environment {
    variables = {
      REGION                       = var.region
      PARAMKEY_BOT_USER_AUTH_TOKEN = var.param_key_bot_user_auth_token
      PARAMKEY_SIGNING_SECRET      = var.param_key_signing_secret
      DYNAMO_DB_NAME               = var.dynamodb_table_name
    }
  }
}

resource "aws_lambda_permission" "interactive" {
  statement_id  = "AllowAPIGatewaySlackHub"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.interactive.arn
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${var.api_gateway_slackhub_execution_arn}/*/*"
}
