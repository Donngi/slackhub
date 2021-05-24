output "dynamodb_table_name" {
  value = aws_dynamodb_table.tools.id
}

output "dynamodb_table_arn" {
  value = aws_dynamodb_table.tools.arn
}
