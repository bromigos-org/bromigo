resource "digitalocean_app" "app" {
  spec {
    name   = var.project_name
    region = var.region

    env {
      key   = "AWS_ACCESS_KEY_ID"
      value = var.do_spaces_access_key
    }

    env {
      key   = "AWS_SECRET_ACCESS_KEY"
      value = var.do_spaces_secret_key
    }

    env {
      key   = "DO_SPACE_NAME"
      value = var.do_space_name
    }

		env {
			key   = "DISCORD_BOT_TOKEN"
			value = var.discord_bot_token
		}

    service {
      name            = var.project_name
      build_command = "go build -o main ./cmd/bromigo/main.go"
    	run_command   = "./main"

      github {
        repo           = "blackflame007/nicklesseos.com"
        branch         = "main"
        deploy_on_push = true
      }
    }
  }
}