def prepareTerraform(env) {
    sh "cp terraforms/configs/base-tags.tfvars terraforms/base-env.tfvars"
    sh "cp terraforms/configs/${env}-vars.tfvars terraforms/env.tfvars"
    sh "./scripts/prepare-iam"
}

def initTerraform(moduleName) {
  echo "Initializing ${moduleName}..."

  sh "cd terraforms/${moduleName} && \
      terraform init -backend-config=../base-env.tfvars -reconfigure"
}

def planTerraform(moduleName) {
    initTerraform(moduleName)
    echo "Plannig ${moduleName}..."

    sh "cd terraforms/${moduleName} && \
        terraform plan -var-file=../base-env.tfvars -var-file=../env.tfvars -out=../plan.out"
}

def applyTerraform(moduleName) {
  echo "Applying ${moduleName}..."

  sh "cd terraforms/${moduleName} && \
      terraform apply -auto-approve ../plan.out"
}

def destroyTerraform(moduleName) {
  echo "Destroying ${moduleName}..."

  sh "cd terraforms/${moduleName} && \
      terraform destroy -var-file=../base-env.tfvars -var-file=../env.tfvars -auto-approve"
}

pipeline {
    agent any

    environment {
        terraformEnvironment = "undefined"
        apply = true
    }

    stages {
        stage('Prepare deploy') {            
            steps {
                script {
                    def applyOrDestroy = input(
                        message: "Apply or Destroy the environment?",
                        parameters: [
                            [$class: 'ChoiceParameterDefinition',
                                choices: ["apply", "destroy"],
                                name: "input",
                                description: "Apply or Destroy the environment?"]
                        ]
                    )

                    echo "Branch is ${env.BRANCH_NAME}..."
                    if (env.BRANCH_NAME == "main") {
                        terraformEnvironment = "prod"
                    } else {
                        terraformEnvironment = "lab"
                    }

                    if (applyOrDestroy == "apply") {
                        apply = true
                    } else {
                        apply = false
                    }
                    
                    prepareTerraform(terraformEnvironment)
                }

                checkout scm
            }
        }

        stage('Build') {
            when { expression { return apply } }

            steps {
                script {
                    sh "./scripts/generate-lambda-zips"
                }
            }
        }

        stage('Plan IAM') {
            when { expression { return apply } }

            steps {
                script {
                    planTerraform("iam")
                }
            }
        }

        stage('Apply IAM') {
            when { expression { return apply } }

            steps {
                script {
                    applyTerraform("iam")
                }
            }
        }

        stage('Plan Lambdas') {
            when { expression { return apply } }

            steps {
                script {
                    planTerraform("lambdas")
                }
            }
        }

        stage('Apply Lambdas') {
            when { expression { return apply } }

            steps {
                script {
                    applyTerraform("lambdas")
                }
            }
        }

        stage('Plan API Gateway') {
            when { expression { return apply } }

            steps {
                script {
                    planTerraform("api-gateway-v2")
                }
            }
        }

        stage('Apply API Gateway') {
            when { expression { return apply } }

            steps {
                script {
                    applyTerraform("api-gateway-v2")
                }
            }
        }

        stage('Destroy') {
            when { expression { return !apply } }

            steps {
                script {
                    destroyTerraform("api-gateway-v2")
                    destroyTerraform("lambdas")
                    destroyTerraform("iam")
                }
            }
        }

        stage('Clear artifacts') {
            when { expression { return !apply } }

            steps {
                script {
                    sh "./scripts/clear-functions-zips"
                }
            }
        }
    }
}