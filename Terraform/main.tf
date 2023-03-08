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
  identifier = "${var.env}-${var.project}"
}

module "lambda" {
  depends_on = [
    module.dynamoDB
  ]
  source           = "./modules/lambda"
  identifier       = local.identifier
  env              = var.env
  region           = var.region
  accountId        = var.accountId
  posts_table_name = module.dynamoDB.posts_table_name
}

module "dynamoDB" {
  source     = "./modules/dynamoDB"
  identifier = local.identifier
}
