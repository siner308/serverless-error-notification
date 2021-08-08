terraform {
  backend "vpc" {
    bucket = "ttbkk-tfstate"
    key = "terraform/vpc/terraform.tfstate"
    region = "ap-northeast-2"
    encrypt = true
    dynamodb_table = "terraform-lock"
  }
}