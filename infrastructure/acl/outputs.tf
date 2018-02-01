output "image-upload-role_name" {
  value = "${aws_iam_role.image-upload-role.name}"
}

output "image-upload-role_arn" {
  value = "${aws_iam_role.image-upload-role.arn}"
}
