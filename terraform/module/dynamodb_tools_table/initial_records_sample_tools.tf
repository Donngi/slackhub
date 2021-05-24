resource "aws_dynamodb_table_item" "sample_go" {
  table_name = aws_dynamodb_table.tools.name
  hash_key   = aws_dynamodb_table.tools.hash_key

  item = jsonencode({
    "id" : { "S" : "sample_go" },
    "displayName" : { "S" : "Sample Tool Go" },
    "description" : { "S" : "Sample tool written in Go. Simply returns user's inputs." },
    "calleeArn" : { "S" : var.lambda_sample_go_arn },
    "modalJSON" : { "S" : file("${path.module}/sample_tool_go_modal.json") },
    "bootMode" : { "S" : "Normal" }
  })
}
