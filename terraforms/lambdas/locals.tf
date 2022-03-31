locals {
  lambda_functions = {
    your-lambda-1 = {
      function_name = format("%s-%s", var.project_name, "your-lambda-1")
      description   = "Your lambda function 1"
      runtime       = "go1.x"
      env_vars = {
        ENV = "prod"
      }
    },

    your-lambda-2 = {
      function_name = format("%s-%s", var.project_name, "your-lambda-2")
      description   = "Your lambda function 2"
      runtime       = "go1.x"
      env_vars = {
        ENV = "prod"
      }
    }
  }
}
