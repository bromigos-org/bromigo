variable "do_token" {}

variable "region" {
  description = "The region to deploy to"
  default     = "sfo3"

}

variable "do_space_name" {
  description = "The name of your DigitalOcean Space"
  default     = "bromigo-space"
}


variable "do_spaces_access_key" {
  description = "DigitalOcean Spaces Access Key"
  type        = string
}

variable "do_spaces_secret_key" {
  description = "DigitalOcean Spaces Secret Key"
  type        = string
}

variable "repo_name" {
  description = "The name of the GitHub repository"
  default     = "bromigos-org/bromigo"
  
}

variable "project_name" {
  description = "The name of the DigitalOcean App"
  default     = "bromigo-app"
  
}

variable "discord_bot_token" {
  description = "The token for the Discord bot"
  type        = string
}