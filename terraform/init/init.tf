provider "aws" {
  region = var.aws_region
}

resource "aws_s3_bucket" "tfstate" {
  bucket = "error-helper-tfstate"

  versioning {
    enabled = true
  }
}

resource "aws_dynamodb_table" "terraform_state_lock" {
  name = "error-helper-tflock"
  hash_key = "LockID"
  billing_mode = "PAY_PER_REQUEST"

  attribute {
    name = "LockID"
    type = "S"
  }
}
