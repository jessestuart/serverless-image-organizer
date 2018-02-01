variable "apex_function_image-upload-handler" {}

resource "aws_lambda_permission" "allow_bucket" {
  statement_id  = "AllowExecutionFromS3Bucket"
  action        = "lambda:InvokeFunction"
  principal     = "s3.amazonaws.com"
  function_name = "${var.apex_function_image-upload-handler}"
  source_arn    = "${module.data-storage.photos_bucket_arn}"
}

resource "aws_s3_bucket_notification" "bucket_notification" {
  bucket = "${module.data-storage.photos_bucket}"

  lambda_function {
    filter_suffix       = "jpg"
    lambda_function_arn = "${var.apex_function_image-upload-handler}"
    events              = ["s3:ObjectCreated:*"]
  }
}
