output "photos_bucket" {
  value = "${aws_s3_bucket.photos_bucket.id}"
}

output "photos_bucket_arn" {
  value = "${aws_s3_bucket.photos_bucket.arn}"
}
