# ----------------------------------------------------------
# API Gateway basic config
# ----------------------------------------------------------

resource "aws_apigatewayv2_api" "slackhub" {
  name          = "SlackHubAPI"
  protocol_type = "HTTP"
}

resource "aws_apigatewayv2_stage" "default" {
  api_id      = aws_apigatewayv2_api.slackhub.id
  name        = "$default"
  auto_deploy = true
  access_log_settings {
    destination_arn = aws_cloudwatch_log_group.api_gateway_slackhub.arn
    format          = jsonencode({ "requestId" : "$context.requestId", "ip" : "$context.identity.sourceIp", "requestTime" : "$context.requestTime", "httpMethod" : "$context.httpMethod", "routeKey" : "$context.routeKey", "status" : "$context.status", "protocol" : "$context.protocol", "responseLength" : "$context.responseLength" })
  }
}

resource "aws_cloudwatch_log_group" "api_gateway_slackhub" {
  name              = "/aws/apigateway/SlackHubAPI"
  retention_in_days = 7
}

# ----------------------------------------------------------
# Endpoint - Initial
# ----------------------------------------------------------

resource "aws_apigatewayv2_integration" "initial" {
  api_id                 = aws_apigatewayv2_api.slackhub.id
  integration_type       = "AWS_PROXY"
  integration_uri        = var.lambda_initial_invoke_arn
  payload_format_version = "2.0"
}

resource "aws_apigatewayv2_route" "initial" {
  api_id    = aws_apigatewayv2_api.slackhub.id
  route_key = "POST /initial"
  target    = "integrations/${aws_apigatewayv2_integration.initial.id}"
}

# ----------------------------------------------------------
# Endpoint - Interactive
# ----------------------------------------------------------

resource "aws_apigatewayv2_integration" "interactive" {
  api_id                 = aws_apigatewayv2_api.slackhub.id
  integration_type       = "AWS_PROXY"
  integration_uri        = var.lambda_interactive_invoke_arn
  payload_format_version = "2.0"
}

resource "aws_apigatewayv2_route" "interactive" {
  api_id    = aws_apigatewayv2_api.slackhub.id
  route_key = "POST /interactive"
  target    = "integrations/${aws_apigatewayv2_integration.interactive.id}"
}
