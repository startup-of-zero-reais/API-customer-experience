locals {
  lambda_functions = {
    company = {
      function_name = format("%s-%s", var.project_name, "company")
      description   = "Company resources"
      runtime       = "go1.x"
      env_vars = {
        ENVIRONMENT = "lab"
        JWT_SERVICE = "DEV"
        LOG_LEVEL   = 2
        APP_URL     = "https://customer-experience.zero-reais-lab.cloud"
      }
    }

    favorites = {
      function_name = format("%s-%s", var.project_name, "favorites")
      description   = "Favorites resources"
      runtime       = "go1.x"
      env_vars = {
        ENVIRONMENT = "lab"
        JWT_SERVICE = "DEV"
        LOG_LEVEL   = 2
        APP_URL     = "https://customer-experience.zero-reais-lab.cloud"
      }
    }

    # orders = {
    #   function_name = format("%s-%s", var.project_name, "orders")
    #   description   = "Orders resources"
    #   runtime       = "go1.x"
    #   env_vars = {
    #     ENVIRONMENT = "lab"
    #   }
    # }

    session = {
      function_name = format("%s-%s", var.project_name, "session")
      description   = "Session resources"
      runtime       = "go1.x"
      env_vars = {
        ENVIRONMENT = "lab"
        JWT_SERVICE = "DEV"
        LOG_LEVEL   = 2
      }
    }
    user = {
      function_name = format("%s-%s", var.project_name, "user")
      description   = "User resources"
      runtime       = "go1.x"
      env_vars = {
        ENVIRONMENT = "lab"
        JWT_SERVICE = "DEV"
        LOG_LEVEL   = 2
      }
    }
  }
}
