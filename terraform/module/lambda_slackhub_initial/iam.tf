# ----------------------------------------------------------
# IAM Role for Lambda
# ----------------------------------------------------------

resource "aws_iam_role" "initial" {
  name = "SlackHubInitialRole"

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
  name        = "SlackHubInitialSlackHubToolsTablePolicy"
  description = "Allow lambda to access to specific DynamoDB table"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = [
          "dynamodb:Scan"
        ]
        Effect   = "Allow"
        Resource = var.dynamodb_table_arn
      },
    ]
  })
}

resource "aws_iam_role_policy_attachment" "allow_access_to_tools_table" {
  role       = aws_iam_role.initial.name
  policy_arn = aws_iam_policy.allow_access_to_tools_table.arn
}

# ----------------------------------------------------------
# LambdaBasicExecutionRole
# ----------------------------------------------------------
resource "aws_iam_role_policy_attachment" "lambda_basic_execution" {
  role       = aws_iam_role.initial.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

# ----------------------------------------------------------
# Allow access to SSM
# ----------------------------------------------------------
resource "aws_iam_role_policy_attachment" "ssm_read_only" {
  role       = aws_iam_role.initial.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonSSMReadOnlyAccess"
}
