terraform {
  required_providers {
    vault = {
      source  = "hashicorp/vault"
      version = ">= 5"
    }
  }
}

variable "name" {
  type        = string
  description = "AppRole receiving the DNS credential read policy"
}

variable "secrets" {
  type        = string
  description = "KV v2 secrets mount"
}

variable "views" {
  type        = set(string)
  description = "DNS provider views whose credentials the owner needs"

  validation {
    condition     = length(var.views) > 0 && length(setsubtract(var.views, ["global", "dc1"])) == 0
    error_message = "DNS access requires at least one of the global or dc1 views."
  }
}

locals {
  secret_paths = {
    global = "cloudflare.com/dns_token"
    dc1    = "alwaldend.com/vault1/approles/src_infra_dns/mikrotik"
  }
}

resource "vault_policy" "dns" {
  name = "${var.name}_dns"
  policy = join("\n", [
    for view in sort(tolist(var.views)) : <<-EOT
      path "${var.secrets}/data/${local.secret_paths[view]}" {
        capabilities = ["read"]
      }
    EOT
  ])
}

output "policy_name" {
  description = "Policy granting reads of the required provider credentials"
  value       = vault_policy.dns.name
}
