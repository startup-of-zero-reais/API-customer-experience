locals {
  lambda_functions = {
    # company = {
    #   function_name = format("%s-%s", var.project_name, "company")
    #   description   = "Company resources"
    #   runtime       = "go1.x"
    #   env_vars = {
    #     ENV = "prod"
    #   }
    # }

    # favorites = {
    #   function_name = format("%s-%s", var.project_name, "favorites")
    #   description   = "Favorites resources"
    #   runtime       = "go1.x"
    #   env_vars = {
    #     ENV = "prod"
    #   }
    # }

    # orders = {
    #   function_name = format("%s-%s", var.project_name, "orders")
    #   description   = "Orders resources"
    #   runtime       = "go1.x"
    #   env_vars = {
    #     ENV = "prod"
    #   }
    # }

    session = {
      function_name = format("%s-%s", var.project_name, "session")
      description   = "Session resources"
      runtime       = "go1.x"
      env_vars = {
        ENVIRONMENT = "prod"
        JWT_SERVICE = "DEV"
        LOG_LEVEL   = 2
      }
    }
    user = {
      function_name = format("%s-%s", var.project_name, "user")
      description   = "User resources"
      runtime       = "go1.x"
      env_vars = {
        ENVIRONMENT = "prod"
        JWT_SERVICE = "DEV"
        LOG_LEVEL   = 2
      }
    }
  }
}
