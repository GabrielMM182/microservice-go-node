# Configure the AWS Provider
provider "aws" {
  region = var.aws_region
}

# S3 bucket for storing CSV reports
resource "aws_s3_bucket" "reports_bucket" {
  bucket = var.s3_report_bucket_name

  tags = {
    Name        = var.s3_report_bucket_name
    Environment = var.environment
    Project     = "Todo Reports"
  }
}

# Block public access to the S3 bucket
resource "aws_s3_bucket_public_access_block" "reports_bucket_public_access_block" {
  bucket = aws_s3_bucket.reports_bucket.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

# S3 bucket versioning
resource "aws_s3_bucket_versioning" "reports_bucket_versioning" {
  bucket = aws_s3_bucket.reports_bucket.id
  versioning_configuration {
    status = "Enabled"
  }
}

# SES email identity
resource "aws_ses_email_identity" "sender" {
  email = var.ses_sender_email
}

# IAM policy for S3 and SES access
resource "aws_iam_policy" "app_policy" {
  name        = "${var.environment}-todo-app-policy"
  description = "Policy for Todo Reports application"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "s3:PutObject",
          "s3:GetObject"
        ]
        Resource = [
          "${aws_s3_bucket.reports_bucket.arn}/*"
        ]
      },
      {
        Effect = "Allow"
        Action = [
          "ses:SendRawEmail",
          "ses:SendEmail"
        ]
        Resource = [
          aws_ses_email_identity.sender.arn
        ]
      }
    ]
  })
}

# IAM user for application access
resource "aws_iam_user" "app_user" {
  name = "${var.environment}-todo-app-user"

  tags = {
    Environment = var.environment
    Project     = "Todo Reports"
  }
}

# Attach policy to user
resource "aws_iam_user_policy_attachment" "app_user_policy" {
  user       = aws_iam_user.app_user.name
  policy_arn = aws_iam_policy.app_policy.arn
}

# Generate access keys for the IAM user
resource "aws_iam_access_key" "app_user_key" {
  user = aws_iam_user.app_user.name
}