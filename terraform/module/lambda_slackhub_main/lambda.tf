data "archive_file" "zip" {
  type        = "zip"
  source_file = "${var.source_code_dir}/${var.source_code_file}"
  output_path = "${var.source_code_dir}/upload_terraform.zip"
}

resource "aws_lambda_function" "slackhub" {
  filename      = "${var.source_code_dir}/upload_terraform.zip"
  function_name = var.function_name
  role          = aws_iam_role.slackhub_lambda.arn
  handler       = "main"

  runtime = "go1.x"

  environment {
    variables = {
      REGION                       = var.region
      PARAMKEY_BOT_USER_AUTH_TOKEN = var.param_key_bot_user_auth_token
      PARAMKEY_SIGNING_SECRET      = var.param_key_signing_secret
      DYNAMO_DB_NAME               = var.dynamodb_table_name
    }
  }
}

resource "aws_lambda_permission" "slackhub" {
  statement_id  = "AllowAPIGatewaySlackHub"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.slackhub.arn
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${var.api_gateway_slackhub_execution_arn}/*/*"
}
