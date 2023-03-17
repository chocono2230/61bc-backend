terraform {
  required_version = "1.3.9"

  backend "s3" {
  }

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~>4.56.0"
    }
    null = {
      source  = "hashicorp/null"
      version = ">= 3.1.1"
    }
  }
}

provider "aws" {
  region = "ap-northeast-1"
  default_tags {
    tags = {
      from    = "terraform"
      project = var.project
      env     = var.env
    }
  }
}

locals {
  identifier        = "${var.env}-${var.project}"
  gsi_name_all      = "${local.identifier}-ddb-posts-gsi-alltime"
  gsi_name_usr      = "${local.identifier}-ddb-posts-gsi-usrtime"
  gsi_name_identity = "${local.identifier}-ddb-users-gsi-identity"
}

module "lambda" {
  depends_on = [
    module.dynamoDB,
    module.s3,
  ]
  source                        = "./modules/lambda"
  identifier                    = local.identifier
  env                           = var.env
  region                        = var.region
  accountId                     = var.accountId
  posts_table_name              = module.dynamoDB.posts_table_name
  posts_table_gsi_name_all      = local.gsi_name_all
  posts_table_gsi_name_usr      = local.gsi_name_usr
  users_table_name              = module.dynamoDB.users_table_name
  users_table_gsi_name_identity = local.gsi_name_identity
  image_bucket_name             = module.s3.bucket_name
}

module "dynamoDB" {
  source            = "./modules/dynamoDB"
  identifier        = local.identifier
  gsi_name_all      = local.gsi_name_all
  gsi_name_usr      = local.gsi_name_usr
  gsi_name_identity = local.gsi_name_identity
}

module "s3" {
  source     = "./modules/s3"
  identifier = "${local.identifier}-image"
  env        = var.env
}
