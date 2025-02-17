
variable "aws_secret_access_key" {
  type = string
}
variable "aws_secret_key_id" {
  type = string
}
variable "aws_default_region" {
  type = string
}
variable "s3_bucket" {
  type = string
}

variable "sqs_queue" {
  type = string
}
variable "s3_localstack_endpoint" {
  type = string
}
provider "aws" {

  access_key                  = "mock_access_key"
  secret_key                  = "mock_secret_key"
  region                      = "us-east-1"

  s3_use_path_style           = true
  skip_credentials_validation = true
  skip_metadata_api_check     = true
  skip_requesting_account_id  = true

  endpoints {
    s3             = "http://s3.localhost.localstack.cloud:4566"
    sqs            = "http://localhost:4566"
  }
}

resource "aws_s3_bucket" "Reports-s3-bucket" {
  bucket = var.s3_bucket
}

resource "aws_sqs_queue" "terraform-sqs-queue" {
  name                      = var.sqs_queue
  delay_seconds             = 5
  max_message_size          = 2048
  message_retention_seconds = 86400
  receive_wait_time_seconds = 16

}