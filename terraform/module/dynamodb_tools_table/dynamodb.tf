resource "aws_dynamodb_table" "tools" {
  name           = "SlackHubToolsTable"
  billing_mode   = "PROVISIONED"
  read_capacity  = 5
  write_capacity = 5
  
  hash_key       = "id"

  attribute {
    name = "id"
    type = "S"
  }
}