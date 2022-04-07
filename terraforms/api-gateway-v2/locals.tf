locals {
  api_stage = "api"

  routes = {
    create_user = {
      version     = "v1"
      prefix      = "customer-experience/"
      path        = "user"
      method      = "POST"
      invoke_arn  = data.aws_lambda_function.user.invoke_arn
      description = "User resources lambda - Method: POST"
    }

    get_user = {
      version     = "v1"
      prefix      = "customer-experience/"
      path        = "user"
      method      = "GET"
      invoke_arn  = data.aws_lambda_function.user.invoke_arn
      description = "User resources lambda - Method: GET"
    }

    put_user = {
      version     = "v1"
      prefix      = "customer-experience/"
      path        = "user"
      method      = "PUT"
      invoke_arn  = data.aws_lambda_function.user.invoke_arn
      description = "User resources lambda - Method: PUT"
    }

    delete_user = {
      version     = "v1"
      prefix      = "customer-experience/"
      path        = "user"
      method      = "DELETE"
      invoke_arn  = data.aws_lambda_function.user.invoke_arn
      description = "User resources lambda - Method: DELETE"
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

