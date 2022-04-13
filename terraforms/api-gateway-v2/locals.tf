locals {
  api_stage = "api"

  routes = {
    create_user = {
      version     = "v1"
      prefix      = "customer-experience/"
      path        = "user"
      method      = "POST"
      invoke_arn  = data.aws_lambda_function.user.invoke_arn
      description = "Customer user resources lambda - Method: POST"
    }

    get_user = {
      version     = "v1"
      prefix      = "customer-experience/"
      path        = "user"
      method      = "GET"
      invoke_arn  = data.aws_lambda_function.user.invoke_arn
      description = "Customer user resources lambda - Method: GET"
    }

    put_user = {
      version     = "v1"
      prefix      = "customer-experience/"
      path        = "user"
      method      = "PUT"
      invoke_arn  = data.aws_lambda_function.user.invoke_arn
      description = "Customer user resources lambda - Method: PUT"
    }

    delete_user = {
      version     = "v1"
      prefix      = "customer-experience/"
      path        = "user"
      method      = "DELETE"
      invoke_arn  = data.aws_lambda_function.user.invoke_arn
      description = "Customer user resources lambda - Method: DELETE"
    }

    sign_in = {
      version     = "v1"
      prefix      = "customer-experience/"
      path        = "sign-in"
      method      = "POST"
      invoke_arn  = data.aws_lambda_function.session.invoke_arn
      description = "SignIn resource lambda - Method: POST"
    }

    sign_out = {
      version     = "v1"
      prefix      = "customer-experience/"
      path        = "sign-out"
      method      = "POST"
      invoke_arn  = data.aws_lambda_function.session.invoke_arn
      description = "SignOut resource lambda - Method: POST"
    }

    recover_password = {
      version     = "v1"
      prefix      = "customer-experience/"
      path        = "recover-password"
      method      = "POST"
      invoke_arn  = data.aws_lambda_function.session.invoke_arn
      description = "Recover pass resource lambda - Method: POST"
    }

    reset_password = {
      version     = "v1"
      prefix      = "customer-experience/"
      path        = "reset-password"
      method      = "POST"
      invoke_arn  = data.aws_lambda_function.session.invoke_arn
      description = "Reset pass resource lambda - Method: POST"
    }

    company = {
      version     = "v1"
      prefix      = "customer-experience/"
      path        = "company/{slug}"
      method      = "GET"
      invoke_arn  = data.aws_lambda_function.company.invoke_arn
      description = "Company food menu resource lambda - Method: GET"
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

