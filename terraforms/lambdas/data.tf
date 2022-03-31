data "aws_iam_role" "lambda_iam_role" {
  name = format("%s-iam-for-lambda", var.project_name)
}

data "aws_apigatewayv2_api" "this" {
  api_id = var.api_id
}
