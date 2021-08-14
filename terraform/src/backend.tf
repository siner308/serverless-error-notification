terraform {
  backend "s3" {
    bucket = "error-helper-tfstate"
    key = "src/terraform.tfstate"
    region = "ap-northeast-2"
    encrypt = true
    dynamodb_table = "error-helper-tflock"
  }
}