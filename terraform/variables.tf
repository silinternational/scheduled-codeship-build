variable "aws_region" {
  default = "us-east-1"
}

variable "aws_access_key" {
}

variable "aws_secret_key" {
}

/*
 * App configuration
 */

variable "ssm_param_path" {
  description = "The SSM parameter path to prefix each parameter name"
  default = "scheduled-codeship-build"
}

variable "codeship_organization" {
  description = "The Codeship organization containing the configured projects"
  type = string
}

variable "codeship_projects" {
  description = "A json list of objects specifying the projects to build, e.g. [{\"uuid\":\"26e97136-8265-4172-867d-3392c7b3c322\",\"ref\":\"20.04\"}]"
  type = string
}

/*
 * AWS tag values
 */

variable "app_customer" {
  description = "customer name to use for the itse_app_customer tag"
  type = string
  default = "gtis"
}

variable "app_environment" {
  description = "environment name to use for the itse_app_environment tag, e.g. staging, production"
  type = string
  default = "production"
}

variable "app_name" {
  description = "app name to use for the itse_app_name tag"
  type = string
  default = "scheduled-codeship-build"
}
