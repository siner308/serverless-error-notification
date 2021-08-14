terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
      version = "~> 3.48.0"
    }
    random = {
      source = "hashicorp/random"
      version = "~> 3.1.0"
    }
    archive = {
      source = "hashicorp/archive"
      version = "~> 2.2.0"
    }
  }

  required_version = "~> 1.0"
}

provider "aws" {
  region = var.aws_region
}

resource "random_pet" "lambda_bucket_name" {
  prefix = "error-helper"
  length = 4
}

resource "aws_s3_bucket" "lambda_bucket" {
  bucket = random_pet.lambda_bucket_name.id

  acl = "private"
  force_destroy = true
}

data "local_file" "dotenv" {
  filename = "${path.root}/../../.env"
}

data "local_file" "main" {
  filename = "${path.module}/../../main"
}


data "archive_file" "zip" {
  type = "zip"
  output_path = "main.zip"

  source {
    content = data.local_file.main.content
    filename = "main"
  }

  source {
    content = data.local_file.dotenv.content
    filename = ".env"
  }
}

resource "aws_s3_bucket_object" "lambda" {
  bucket = aws_s3_bucket.lambda_bucket.id

  key = "main.zip"
  source = data.archive_file.zip.output_path

  etag = filemd5(data.archive_file.zip.output_path)
}

resource "aws_lambda_function" "error-helper" {
  function_name = "error-helper"

  s3_bucket = aws_s3_bucket.lambda_bucket.id
  s3_key = aws_s3_bucket_object.lambda.key

  runtime = "go1.x"
  handler = "main"
  memory_size = 128

  source_code_hash = data.archive_file.zip.output_base64sha256

  role = aws_iam_role.lambda_exec.arn
}

resource "aws_cloudwatch_log_group" "error-helper" {
  name = "/aws/lambda/${aws_lambda_function.error-helper.function_name}"

  retention_in_days = 30
}

resource "aws_iam_role" "lambda_exec" {
  name = "serverless_lambda"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Sid = ""
        Principal = {
          Service = "lambda.amazonaws.com"
        }
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "lambda_policy" {
  role = aws_iam_role.lambda_exec.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}