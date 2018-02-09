variable "cloudfront_origin_id" {
  type    = "string"
  default = "S3-serverless-image-organizer"
}

resource "aws_cloudfront_distribution" "s3_distribution" {
  origin {
    domain_name = "${module.data-storage.photos_bucket}.s3.amazonaws.com"
    origin_id   = "${var.cloudfront_origin_id}"
  }

  enabled = true

  default_cache_behavior {
    allowed_methods  = ["GET", "HEAD", "OPTIONS"]
    cached_methods   = ["GET", "HEAD", "OPTIONS"]
    target_origin_id = "${var.cloudfront_origin_id}"

    forwarded_values {
      query_string = false

      cookies {
        forward = "none"
      }
    }

    viewer_protocol_policy = "allow-all"
    min_ttl                = 0
    default_ttl            = 3600
    max_ttl                = 86400
  }

  is_ipv6_enabled = "true"

  price_class = "PriceClass_100"

  tags {
    Environment = "production"
  }

  restrictions {
    geo_restriction {
      restriction_type = "none"
    }
  }

  viewer_certificate {
    cloudfront_default_certificate = true
  }
}
