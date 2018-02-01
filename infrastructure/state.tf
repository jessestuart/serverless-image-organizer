provider "aws" {
  region = "us-east-1"
}

# =====================================================
# Terraform state file setup:
# Create an S3 bucket in which to store the state file.
# =====================================================
resource "aws_s3_bucket" "state_storage" {
  bucket = "js-terraform"

  versioning {
    enabled = true
  }

  lifecycle {
    prevent_destroy = true
  }

  tags {
    Name = "S3 Remote Terraform State Store"
  }
}
