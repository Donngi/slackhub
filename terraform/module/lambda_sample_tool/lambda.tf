data "archive_file" "zip" {
  type        = "zip"
  source_file = "${var.source_code_dir}/${var.source_code_file}"
  output_path = "${var.source_code_dir}/upload_terraform.zip"
}

resource "aws_lambda_function" "sample" {
  filename      = "${var.source_code_dir}/upload_terraform.zip"
  function_name = var.function_name
  role          = aws_iam_role.sample_lambda.arn
  handler       = "main"

  runtime = "go1.x"

  source_code_hash = data.archive_file.zip.output_base64sha256
}
