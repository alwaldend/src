terraform {
  required_providers {
    vastai = {
      source = "alwaldend/vastai"
    }
  }
}

variable "name" {
  type        = string
  description = "Name of the zero-credit team."
}

variable "api_keys" {
  type = map(object({
    name        = string
    permissions = string
  }))
  default     = {}
  description = "Project keys with explicit JSON permission objects."
}

resource "vastai_team" "team" {
  name = var.name
}

resource "vastai_api_key" "project" {
  for_each     = var.api_keys
  name         = each.value.name
  permissions  = each.value.permissions
  team_api_key = vastai_team.team.api_key
}

output "team_id" {
  value       = vastai_team.team.team_id
  description = "Created team account ID."
}

output "api_keys" {
  value       = { for name, key in vastai_api_key.project : name => key.key }
  sensitive   = true
  description = "Generated project tokens for direct delivery to secret resources."
}
