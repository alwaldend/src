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
  description = "Optional reviewed Vast ComfyUI template hash. Omit to use the configured image tag without inherited template options."
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
  default = [
    {
      url      = "https://huggingface.co/lodestones/Chroma1-Flash/resolve/093c4f63b60507234aa53a073f1f9f565b3f4336/Chroma1-HD-Flash.safetensors"
      filename = "Chroma1-HD-Flash.safetensors"
      kind     = "checkpoint"
    },
    {
      url      = "https://huggingface.co/comfyanonymous/flux_text_encoders/resolve/6af2a98e3f615bdfa612fbd85da93d1ed5f69ef5/t5xxl_fp16.safetensors"
      filename = "t5xxl_fp16.safetensors"
      kind     = "text_encoder"
    },
    {
      url      = "https://huggingface.co/lodestones/Chroma/resolve/2f3b2730d7b5edbd02cdff12f72610af5787300b/ae.safetensors"
      filename = "ae.safetensors"
      kind     = "vae"
    }
  ]
  nullable    = false
  description = "Public model files downloaded before ComfyUI starts; defaults to Chroma1-HD-Flash with its T5 encoder and VAE."
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
  source = "../../../../projects/tf_modules/vastai_instance_imagegen_comfyui"

  # A failed Vault write must prevent creation of the paid instance.
  depends_on = [vault_kv_secret_v2.web]

  image            = "vastai/comfy:v0.44.0-cuda-12.9-py312"
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
