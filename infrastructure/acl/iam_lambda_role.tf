resource "aws_iam_role" "image-upload-role" {
  name = "image-upload-role"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow"
    }
  ]
}
EOF
}

# Attach role to Managed Policy
resource "aws_iam_policy_attachment" "lambda_basic_policy_attachment" {
  name       = "AWSLambdaBasicExecutionRole"
  roles      = ["${aws_iam_role.image-upload-role.id}"]
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_iam_role_policy" "lambda_dynamodb_policy" {
  name = "DynamoDBWriteAccess"
  role = "${aws_iam_role.image-upload-role.id}"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
      {
          "Sid": "",
          "Effect": "Allow",
          "Action": "dynamodb:PutItem",
          "Resource": "arn:aws:dynamodb:*:*:table/*"
      }
  ]
}
EOF
}

resource "aws_iam_role_policy" "lambda_logging_policy" {
  name = "LoggingRWAcess"
  role = "${aws_iam_role.image-upload-role.id}"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "logs:CreateLogGroup",
        "logs:CreateLogStream",
        "logs:PutLogEvents"
      ],
      "Resource": "arn:aws:logs:*:*:*"
    }
  ]
}
EOF
}

resource "aws_iam_role_policy" "lambda_s3_access_policy" {
  name = "S3RWAccess"
  role = "${aws_iam_role.image-upload-role.id}"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "s3:GetObject",
        "s3:PutObject",
        "s3:DeleteObject"
      ],
      "Resource": "arn:aws:s3:::*"
    }
  ]
}
EOF
}
