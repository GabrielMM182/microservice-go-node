variable "aws_region" {
  description = "AWS region where resources will be created"
  type        = string
  default     = "us-east-1"
}

variable "environment" {
  description = "Environment name (e.g., dev, prod)"
  type        = string
  default     = "dev"
}

variable "s3_report_bucket_name" {
  description = "Name of the S3 bucket for storing reports"
  type        = string
  validation {
    condition     = can(regex("^[a-z0-9][a-z0-9-]*[a-z0-9]$", var.s3_report_bucket_name))
    error_message = "Bucket name must be lowercase alphanumeric characters and hyphens only."
  }
}

variable "ses_sender_email" {
  description = "Email address to be verified for sending emails through SES"
  type        = string
  validation {
    condition     = can(regex("^[a-zA-Z0-9_%+-]+@[a-zA-Z0-9-]+\\.[a-zA-Z]{2,}$", var.ses_sender_email))
    error_message = "Please provide a valid email address."
  }
}