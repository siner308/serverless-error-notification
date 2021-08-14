variable "aws_region" {
  description = "AWS region for all resources."

  type = string
  default = "ap-northeast-2"
}

variable "archived_function" {
  description = "Archived Go Function"

  type = string
  default = "../../function.zip"
}