resource "digitalocean_app" "app" {
  spec {
    name   = var.app_name
    region = var.region

		env {
			key   = "DISCORD_BOT_TOKEN"
			value = var.discord_bot_token
		}

    service {
      name            = var.app_name
      build_command = var.build_command
    	run_command   = var.run_command
			github {
				repo = var.repo_name
				branch     = var.repo_branch
				deploy_on_push = true
			}

			health_check {
				http_path = "/health"
				port      = 8080
			}
    }
  }
}

data "digitalocean_project" "project" {
  name        = var.project_name
}

resource "digitalocean_project_resources" "project_resources" {
  project = data.digitalocean_project.project.id
  resources = [
     digitalocean_app.app.urn
  ]
}
