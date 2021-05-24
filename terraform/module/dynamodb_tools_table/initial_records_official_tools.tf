resource "aws_dynamodb_table_item" "register" {
  table_name = aws_dynamodb_table.tools.name
  hash_key   = aws_dynamodb_table.tools.hash_key

  item = jsonencode({
    "id" : { "S" : "register" },
    "displayName" : { "S" : "SlackHub - Register" },
    "description" : { "S" : "Register is a SlackHub's official tool for registering your new tools." },
    "calleeArn" : { "S" : var.lambda_register_arn },
    "modalJSON" : { "S" : "SlackHub's official tool Register dinamically creates modal json." },
    "bootMode" : { "S" : "Normal" }
  })
}

resource "aws_dynamodb_table_item" "editor" {
  table_name = aws_dynamodb_table.tools.name
  hash_key   = aws_dynamodb_table.tools.hash_key

  item = jsonencode({
    "id" : { "S" : "editor" },
    "displayName" : { "S" : "SlackHub - Editor" },
    "description" : { "S" : "Editor is a SlackHub's official tool for editing your tools." },
    "calleeArn" : { "S" : var.lambda_editor_arn },
    "modalJSON" : { "S" : "SlackHub's official tool Editor dinamically creates modal json." },
    "bootMode" : { "S" : "Advanced" }
  })
}

resource "aws_dynamodb_table_item" "catalog" {
  table_name = aws_dynamodb_table.tools.name
  hash_key   = aws_dynamodb_table.tools.hash_key

  item = jsonencode({
    "id" : { "S" : "catalog" },
    "displayName" : { "S" : "SlackHub - Catalog" },
    "description" : { "S" : "Catalog is a SlackHub's official tool for listing your tools." },
    "calleeArn" : { "S" : var.lambda_catalog_arn },
    "modalJSON" : { "S" : "SlackHub's official tool Catalog dinamically creates modal json." },
    "bootMode" : { "S" : "Normal" }
  })
}

resource "aws_dynamodb_table_item" "eraser" {
  table_name = aws_dynamodb_table.tools.name
  hash_key   = aws_dynamodb_table.tools.hash_key

  item = jsonencode({
    "id" : { "S" : "eraser" },
    "displayName" : { "S" : "SlackHub - Eraser" },
    "description" : { "S" : "Eraser is a SlackHub's official tool for listing your tools." },
    "calleeArn" : { "S" : var.lambda_eraser_arn },
    "modalJSON" : { "S" : "SlackHub's official tool Eraser dinamically creates modal json." },
    "bootMode" : { "S" : "Normal" }
  })
}
