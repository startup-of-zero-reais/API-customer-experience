data "aws_apigatewayv2_api" "this" {
  api_id = var.api_id
}

data "aws_lambda_function" "ping" {
  function_name = format("%s-%s", var.project_name, "your-lambda-1")
}

data "aws_lambda_function" "ping2" {
  function_name = format("%s-%s", var.project_name, "your-lambda-2")
}

