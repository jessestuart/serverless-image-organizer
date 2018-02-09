variable "aws_region" {}
variable "terraform_state_backend_bucket_name" {}

provider "aws" {
  region = "${var.aws_region}"
}

# =====================================================
# Terraform state file setup:
# Create an S3 bucket in which to store the state file.
# =====================================================
resource "aws_s3_bucket" "state_storage" {
  bucket = "${var.terraform_state_backend_bucket_name}"

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
