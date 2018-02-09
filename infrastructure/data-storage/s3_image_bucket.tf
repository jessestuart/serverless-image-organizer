resource "aws_s3_bucket" "photos_bucket" {
  bucket = "${var.bucket_name}"
  acl    = "private"

  cors_rule {
    allowed_headers = ["*"]
    allowed_methods = ["PUT", "POST"]
    allowed_origins = ["*"]
    expose_headers  = ["ETag"]
    max_age_seconds = 3000
  }
}

resource "aws_s3_bucket" "target_photos_bucket" {
  bucket = "${var.target_photos_bucket_name}"
  acl    = "public-read"

  cors_rule {
    allowed_headers = ["*"]
    allowed_methods = ["GET", "HEAD", "OPTIONS"]
    allowed_origins = ["*"]
    expose_headers  = ["*"]
  }
}
