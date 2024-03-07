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
      build_command = "go build -o main ./cmd/bromigo/main.go"
    	run_command   = "./main"
			github {
				repo = "bromigos-org/bromigo"
				branch     = "main"
				deploy_on_push = true
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