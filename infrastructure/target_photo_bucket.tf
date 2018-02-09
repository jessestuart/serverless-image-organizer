resource "aws_s3_bucket" "target_photos_bucket" {
  bucket = "${var.target_photos_bucket_name}"
  acl    = "public-read"

  cors_rule {
    allowed_headers = ["*"]
    allowed_methods = ["GET", "HEAD"]
    allowed_origins = ["*"]
    expose_headers  = ["ETag"]
  }
}

resource "aws_s3_bucket_policy" "public_read_policy" {
  bucket = "${aws_s3_bucket.target_photos_bucket.id}"

  policy = <<POLICY
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "PublicReadGetObject",
      "Action": "s3:GetObject",
      "Effect": "Allow",
      "Resource": "${aws_s3_bucket.target_photos_bucket.arn}/*",
      "Principal": "*"
    }
  ]
}
POLICY
}
