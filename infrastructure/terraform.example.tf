# Unfortunately Terraform doesn't support any variable
# interpolation in the `terraform` config block, so we
# have to manually ensure these values match those
# defined in `state.tf`.
terraform {
  backend "s3" {
    encrypt = true
    bucket  = "js-terraform"
    region  = "us-east-1"
    key     = "js-photo-upload/terraform.tfstate"
  }
}
