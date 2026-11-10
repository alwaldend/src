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
      url      = "https://huggingface.co/Comfy-Org/Qwen-Image-Edit/resolve/f68ace85e60b4a02a323e394253731947657b7d2/split_files/diffusion_models/qwen_image_edit_2511_fp8mixed.safetensors"
      filename = "qwen_image_edit_2511_fp8mixed.safetensors"
      kind     = "checkpoint"
    },
    {
      url      = "https://huggingface.co/Comfy-Org/Qwen-Image/resolve/1f12b17be14c89b026c51a91d67c32f84bb047bc/split_files/text_encoders/qwen_2.5_vl_7b_fp8_scaled.safetensors"
      filename = "qwen_2.5_vl_7b_fp8_scaled.safetensors"
      kind     = "text_encoder"
    },
    {
      url      = "https://huggingface.co/Comfy-Org/Qwen-Image/resolve/1f12b17be14c89b026c51a91d67c32f84bb047bc/split_files/vae/qwen_image_vae.safetensors"
      filename = "qwen_image_vae.safetensors"
      kind     = "vae"
    }
  ]
  nullable    = false
  description = "Public model files downloaded before ComfyUI starts; defaults to the native Qwen-Image-Edit-2511 bundle."
}

resource "random_password" "web" {
  length  = 48
  special = false
}

resource "vault_kv_secret_v2" "web" {
  mount = "secrets"
  name  = "alwaldend.com/vault1/approles/src_infra_image_gen/web"
  data_json = jsonencode({
    username = var.web_username
    password = random_password.web.result
  })
  delete_all_versions = true
}

module "imagegen" {
  source = "../../../projects/tf_modules/vastai_instance_imagegen_comfyui"

  # A failed Vault write must prevent creation of the paid instance.
  depends_on = [vault_kv_secret_v2.web]

  image           = "vastai/comfy:v0.44.0-cuda-12.9-py312"
  model_downloads = var.model_downloads
  workflows = { for name in ["qwen-image-edit-single-reference", "qwen-image-edit-character-style"] :
    name => file("${path.module}/../workflows/${name}.json")
  }
  template_hash_id = var.template_hash_id
  access_mode      = "exposed"
  web_username     = var.web_username
  web_password     = random_password.web.result
  min_gpu_ram_gb   = 32
  disk_gb          = 100
  max_hourly_price = var.max_hourly_price
  offer_id         = var.offer_id
  label            = "src_infra_image_gen"
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
