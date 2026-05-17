terraform {
  required_version = ">= 1.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
    grafana = {
      source  = "grafana/grafana"
      version = "~> 3.0"
    }
    stripe = {
      source  = "lukasaron/stripe"
      version = "~> 1.0"
    }
  }

  # Uncomment to store state remotely
  # backend "s3" {
  #   bucket = "rabi-food-core-terraform-state"
  #   key    = "rabi-food-core/terraform.tfstate"
  #   region = "us-east-1"
  # }
}

provider "aws" {
  region = var.aws_region
}

provider "stripe" {
  api_key = var.stripe_secret_key
}

provider "grafana" {
  url           = var.grafana_url
  auth          = var.grafana_admin_token
  cloud_api_key = var.grafana_cloud_api_key
}
