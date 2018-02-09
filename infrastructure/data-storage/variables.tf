variable "app_name" {
  type    = "string"
  default = "mps_upload"
}

variable "aws_region" {
  type    = "string"
  default = "us-east-1"
}

variable "bucket_name" {
  type    = "string"
  default = "mps_upload_bucket"
}

variable "target_photos_bucket_name" {
  type    = "string"
  default = "js-photos"
}
