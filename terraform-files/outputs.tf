output "s3_bucket_name" {
  description = "Name of the created S3 bucket"
  value       = aws_s3_bucket.reports_bucket.id
}

output "s3_bucket_arn" {
  description = "ARN of the created S3 bucket"
  value       = aws_s3_bucket.reports_bucket.arn
}

output "ses_identity_arn" {
  description = "ARN of the SES email identity"
  value       = aws_ses_email_identity.sender.arn
}

output "iam_user_name" {
  description = "Name of the created IAM user"
  value       = aws_iam_user.app_user.name
}

output "iam_access_key_id" {
  description = "Access key ID for the IAM user"
  value       = aws_iam_access_key.app_user_key.id
  # sensitive   = true
}

output "iam_access_key_secret" {
  description = "Secret access key for the IAM user"
  value       = aws_iam_access_key.app_user_key.secret
  sensitive   = true
}