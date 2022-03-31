locals {
  api_stage = "api"

  routes = {
    route-1 = {
      version     = "v1"
      prefix      = "application-resource/"
      path        = "ping"
      method      = "GET"
      invoke_arn  = data.aws_lambda_function.ping.invoke_arn
      description = "Ping lambda"
    }

    # route-2 = {
    #   version       = "v1"
    #   prefix        = "application-resource/"
    #   path          = "ping-2"
    #   method        = "GET"
    #   invoke_arn    = data.aws_lambda_function.ping2.invoke_arn
    #   description   = "Ping lambda"
    # }
  }
}

