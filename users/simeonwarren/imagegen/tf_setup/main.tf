terraform {
  required_providers {
    vastai = {
      source  = "alwaldend/vastai"
      version = "0.0.0-dev"
    }
    vault = {
      source  = "hashicorp/vault"
      version = ">= 5"
    }
    external = {
      source = "hashicorp/external"
    }
    random = {
      source  = "hashicorp/random"
      version = "3.9.1"
    }
  }
}

provider "vastai" {}
provider "vault" {
  skip_child_token = true
}
provider "random" {}

variable "template_hash_id" {
  type        = string
  default     = null
  description = "Optional reviewed Vast Forge template hash. Omit to use the newest selected runtime without inherited template options."
}

variable "web_username" {
  type        = string
  default     = "imagegen"
  nullable    = false
  description = "HTTPS login name stored alongside the generated password in Vault."
}

variable "max_hourly_price" {
  type        = number
  default     = 1
  nullable    = false
  description = "Maximum USD/hour including requested disk; not a total spending limit."
}

variable "offer_id" {
  type        = number
  default     = null
  description = "Optional reviewed marketplace offer ID; all module requirements still apply."
}

variable "model_downloads" {
  type = list(object({
    url      = string
    filename = string
    kind     = optional(string, "checkpoint")
  }))
  default = [{
    url      = "https://huggingface.co/RunDiffusion/Juggernaut-XL-v9/resolve/cf419233522daa0b9ea36c3aff98fa2cab1fb0fb/Juggernaut-XL_v9_RunDiffusionPhoto_v2.safetensors"
    filename = "Juggernaut-XL_v9_RunDiffusionPhoto_v2.safetensors"
    kind     = "checkpoint"
  }]
  nullable    = false
  description = "Public model files downloaded before Forge starts; defaults to the approved Juggernaut XL v9 checkpoint."
}

data "external" "forge_image" {
  program = [abspath("${path.module}/../../../../tools/vastai/cmd/latest_image/latest_image_/latest_image")]
  query   = {}
}

resource "random_password" "web" {
  length  = 48
  special = false
}

resource "vault_kv_secret_v2" "web" {
  mount = "secrets"
  name  = "alwaldend.com/vault1/approles/user_simeonwarren/imagegen/web"
  data_json = jsonencode({
    username = var.web_username
    password = random_password.web.result
  })
  delete_all_versions = false
}

module "imagegen" {
  source = "../../../../projects/tf_modules/vastai_instance_imagegen"

  # A failed Vault write must prevent creation of the paid instance.
  depends_on = [vault_kv_secret_v2.web]

  image            = data.external.forge_image.result.image
  model_downloads  = var.model_downloads
  template_hash_id = var.template_hash_id
  access_mode      = "exposed"
  web_username     = var.web_username
  web_password     = random_password.web.result
  min_gpu_ram_gb   = 32
  disk_gb          = 100
  max_hourly_price = var.max_hourly_price
  offer_id         = var.offer_id
  label            = "simeonwarren_imagegen"
}

output "instance_id" {
  value = module.imagegen.instance_id
}

output "https_url" {
  value = module.imagegen.https_url
}

output "ssh_forward_command" {
  value = module.imagegen.ssh_forward_command
}

output "credentials_path" {
  value = "secrets/${vault_kv_secret_v2.web.name}"
}
