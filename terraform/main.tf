terraform {
  required_providers {
    digitalocean = {
      source  = "digitalocean/digitalocean"
      version = "~> 2.0"
    }
  }
}

# Configure the DigitalOcean Provider
provider "digitalocean" {
  token = var.do_token
  spaces_access_id = var.do_spaces_access_key
  spaces_secret_key = var.do_spaces_secret_key
}

terraform {
  backend "s3" {
    endpoint   = "https://sfo3.digitaloceanspaces.com"  # San Francisco endpoint, closest to Los Angeles
    key        = "tf/nicklesseos.com/terraform.tfstate"
    region     = "us-east-1"                    # Dummy region for AWS S3 compatibility

    skip_requesting_account_id = true
    skip_credentials_validation = true
    skip_metadata_api_check     = true
    force_path_style            = true
    skip_region_validation = true
    skip_s3_checksum = true
  }
}