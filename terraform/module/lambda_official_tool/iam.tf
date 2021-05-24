# ----------------------------------------------------------
# IAM Role for Lambda
# ----------------------------------------------------------

resource "aws_iam_role" "official_lambda" {
  name = "${var.function_name}Role"

  assume_role_policy = jsonencode(
    {
      "Version" : "2012-10-17",
      "Statement" : [
        {
          "Action" : "sts:AssumeRole",
          "Principal" : {
            "Service" : "lambda.amazonaws.com"
          },
          "Effect" : "Allow",
          "Sid" : ""
        }
      ]
    }
  )
}

# ----------------------------------------------------------
# Allow access to SlackHub tools table
# ----------------------------------------------------------
resource "aws_iam_policy" "allow_access_to_tools_table" {
  name        = "${var.function_name}SlackHubToolsTablePolicy"
  description = "Allow lambda to access to specific DynamoDB table"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = [
          "dynamodb:GetItem",
          "dynamodb:PutItem",
          "dynamodb:Scan",
          "dynamodb:DeleteItem"
        ]
        Effect   = "Allow"
        Resource = var.dynamodb_table_arn
      },
      {
        Action = [
          "dynamodb:ListTables"
        ]
        Effect   = "Allow"
        Resource = "*"
      },
    ]
  })
}

resource "aws_iam_role_policy_attachment" "allow_access_to_tools_table" {
  role       = aws_iam_role.official_lambda.name
  policy_arn = aws_iam_policy.allow_access_to_tools_table.arn
}

# ----------------------------------------------------------
# LambdaBasicExecutionRole
# ----------------------------------------------------------
resource "aws_iam_role_policy_attachment" "lambda_basic_execution" {
  role       = aws_iam_role.official_lambda.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

# ----------------------------------------------------------
# Allow access to SSM
# ----------------------------------------------------------
resource "aws_iam_role_policy_attachment" "ssm_read_only" {
  role       = aws_iam_role.official_lambda.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonSSMReadOnlyAccess"
}
