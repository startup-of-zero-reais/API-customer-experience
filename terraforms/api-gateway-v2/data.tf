data "aws_apigatewayv2_api" "this" {
  api_id = var.api_id
}

# data "aws_lambda_function" "company" {
#   function_name = format("%s-%s", var.project_name, "company")
# }

# data "aws_lambda_function" "favorites" {
#   function_name = format("%s-%s", var.project_name, "favorites")
# }

# data "aws_lambda_function" "orders" {
#   function_name = format("%s-%s", var.project_name, "orders")
# }

data "aws_lambda_function" "company" {
  function_name = format("%s-%s", var.project_name, "company")
}

data "aws_lambda_function" "session" {
  function_name = format("%s-%s", var.project_name, "session")
}

data "aws_lambda_function" "user" {
  function_name = format("%s-%s", var.project_name, "user")
}


