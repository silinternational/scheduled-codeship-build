
/*
 * Create IAM user for Serverless framework to use to deploy the lambda function
 */
module "serverless-user" {
  source  = "silinternational/serverless-user/aws"
  version = "0.1.3"

  app_name   = "scheduled-codeship-build"
  aws_region = var.aws_region
}

output "serverless-access-key-id" {
  value = module.serverless-user.aws_access_key_id
}

output "serverless-secret-access-key" {
  value = module.serverless-user.aws_secret_access_key
}

resource "aws_ssm_parameter" "username" {
  name = "${var.ssm_param_path}/username"
  type = "String"
  insecure_value = var.codeship_username
}

resource "aws_ssm_parameter" "password" {
  name = "${var.ssm_param_path}/password"
  type = "SecureString"
  value = var.codeship_password
}

resource "aws_ssm_parameter" "organization" {
  name = "${var.ssm_param_path}/organization"
  type = "String"
  insecure_value = var.codeship_organization
}

resource "aws_ssm_parameter" "projects" {
  name = "${var.ssm_param_path}/projects"
  type = "String"
  insecure_value = var.codeship_projects
}
