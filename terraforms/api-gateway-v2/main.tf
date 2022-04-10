
terraform {
  required_version = "~> 1.1.7"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.8.0"
    }
  }

  backend "s3" {
    key     = "application-name/api-gateway-v2/terraform.tfstate"
    encrypt = false
  }
}

provider "aws" {
  region = "us-east-1"
}

resource "aws_apigatewayv2_integration" "this" {
  for_each = local.routes

  api_id = data.aws_apigatewayv2_api.this.id

  integration_method     = "POST"
  integration_type       = "AWS_PROXY"
  integration_uri        = lookup(each.value, "invoke_arn")
  description            = lookup(each.value, "description")
  payload_format_version = "2.0"
}

resource "aws_apigatewayv2_route" "this" {
  for_each = local.routes

  api_id = data.aws_apigatewayv2_api.this.id
  route_key = format(
    "%s /%s/%s%s",
    lookup(each.value, "method"),
    lookup(each.value, "version"),
    lookup(each.value, "prefix"),
    lookup(each.value, "path"),
  )

  target = "integrations/${aws_apigatewayv2_integration.this[each.key].id}"
}
