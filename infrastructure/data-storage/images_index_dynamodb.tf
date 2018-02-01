resource "aws_dynamodb_table" "image-db" {
  name           = "image-db"
  read_capacity  = 5
  write_capacity = 5
  hash_key       = "Gallery"
  range_key      = "FilePath"

  attribute {
    name = "Gallery"
    type = "S"
  }

  attribute {
    name = "FilePath"
    type = "S"
  }

  tags {
    Name        = "serverless-image-organizer"
    Environment = "dev"
  }
}
