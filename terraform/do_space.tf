resource "digitalocean_spaces_bucket" "space" {
  name   = var.do_space_name # Replace with your desired Space name
  region = var.region        # San Francisco data center, closest to Los Angeles

  acl = "private"
}