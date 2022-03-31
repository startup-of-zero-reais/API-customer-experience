locals {
  api_stage = "api"

  routes = {
    user = {
      version     = "v1"
      prefix      = "customer-experience/"
      path        = "user"
      method      = "GET"
      invoke_arn  = data.aws_lambda_function.user.invoke_arn
      description = "User resources lambda"
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

