variable "context" {
  type        = string
  description = "The project context"
}

variable "project_name" {
  type        = string
  description = "The name of the project"
}

variable "repo_name" {
  type        = string
  description = "The name of the repository. Ex: https://github.com/org/{{ repository_name }}"
}

variable "region" {
  type        = string
  description = "The AWS region"
}

variable "route53_zone_name" {
  type        = string
  description = "The name of the Route53 zone"
}
