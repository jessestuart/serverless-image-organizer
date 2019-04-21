resource "aws_s3_bucket" "photos_bucket" {
  bucket = "${var.bucket_name}"
  acl    = "public-read"

  acceleration_status = "Enabled"

  cors_rule {
    allowed_headers = ["*"]
    allowed_methods = ["PUT", "POST"]
    allowed_origins = ["*"]
    expose_headers  = ["ETag"]
    max_age_seconds = 3000
  }
}
