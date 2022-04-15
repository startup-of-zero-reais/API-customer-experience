terraform {
  required_version = "~> 1.1.7"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.9.0"
    }
  }

  backend "s3" {
    key     = "application-name/lambda/terraform.tfstate"
    encrypt = false
  }
}

provider "aws" {
  region = "us-east-1"
}

resource "aws_lambda_function" "this" {
  for_each = local.lambda_functions

  filename      = format("./functions/%s.zip", each.key)
  function_name = format("%s", lookup(each.value, "function_name", ""))
  role          = data.aws_iam_role.lambda_iam_role.arn
  handler       = each.key
  description   = lookup(each.value, "description", each.key)
  timeout       = 10

  source_code_hash = filebase64sha256(format("./functions/%s.zip", each.key))

  runtime = lookup(each.value, "runtime", "go1.x")

  environment {
    variables = lookup(each.value, "env_vars", {
      ENV = "lab"
    })
  }
}

resource "aws_lambda_permission" "this" {
  for_each = local.lambda_functions

  statement_id  = format("AllowExecutionFromAPIGateway-for-%s", each.key)
  function_name = lookup(each.value, "function_name")
  action        = "lambda:InvokeFunction"
  principal     = "apigateway.amazonaws.com"

  source_arn = "${data.aws_apigatewayv2_api.this.execution_arn}/*/*"
}

resource "aws_cloudwatch_log_group" "this" {
  for_each = local.lambda_functions

  name              = format("/aws/lambda/%s", lookup(each.value, "function_name"))
  retention_in_days = 7

  tags = { "Name" = format("%s-%s", var.project_name, each.key) }
}

resource "aws_cloudwatch_log_stream" "foo" {
  for_each = local.lambda_functions

  name           = format("/first-log/group/%s", lookup(each.value, "function_name"))
  log_group_name = aws_cloudwatch_log_group.this[each.key].name
}
resource "aws_cloudwatch_log_subscription_filter" "lambda_logs" {
  for_each = local.lambda_functions

  name            = lookup(each.value, "function_name")
  log_group_name  = aws_cloudwatch_log_group.this[each.key].name
  filter_pattern  = ""
  destination_arn = data.aws_lambda_function.logs_destination.arn
}
