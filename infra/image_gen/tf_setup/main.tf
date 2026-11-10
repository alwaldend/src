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

locals {
  configuration   = jsondecode(file("${path.module}/../config.json"))
  model           = local.configuration.model_profiles[local.configuration.model_profile]
  hardware        = local.configuration.hardware_profiles[local.configuration.hardware_profile]
  model_downloads = local.model.downloads
}

data "vault_kv_secret_v2" "huggingface" {
  count = local.configuration.model_profile == "flux" && length(local.model_downloads) > 0 ? 1 : 0
  mount = "secrets"
  name  = "alwaldend.com/vault1/approles/src_infra_image_gen/huggingface"
}

resource "random_password" "web" {
  length  = 48
  special = false
}

resource "vault_kv_secret_v2" "web" {
  mount = "secrets"
  name  = "alwaldend.com/vault1/approles/src_infra_image_gen/web"
  data_json = jsonencode({
    username = local.configuration.web_username
    password = random_password.web.result
  })
  delete_all_versions = true
}

module "imagegen" {
  source = "../../../projects/tf_modules/vastai_instance_imagegen_comfyui"

  # A failed Vault write must prevent creation of the paid instance.
  depends_on = [vault_kv_secret_v2.web]

  image                   = local.configuration.image
  model_downloads         = local.model_downloads
  checkpoint_directory    = try(local.model.checkpoint_directory, "models/diffusion_models")
  custom_nodes            = try(local.model.custom_nodes, {})
  python_requirements     = try(local.model.python_requirements, null) == null ? null : file("${path.module}/../${local.model.python_requirements}")
  custom_node_sources     = { prompt_sequence = file("${path.module}/../../../projects/comfyui_tools/nodes/prompt_sequence/prompt_sequence.py") }
  custom_node_web_scripts = { prompt_sequence = file("${path.module}/../../../projects/comfyui_tools/nodes/prompt_sequence/web/index.js") }
  hugging_face_token      = length(data.vault_kv_secret_v2.huggingface) == 0 ? null : data.vault_kv_secret_v2.huggingface[0].data["token"]
  workflows = { for name in local.model.workflows :
    name => file("${path.module}/../../../projects/comfyui_tools/workflows/${local.configuration.model_profile}/${name}.json")
  }
  access_mode      = "exposed"
  web_username     = local.configuration.web_username
  web_password     = random_password.web.result
  min_gpu_ram_gb   = local.hardware.min_gpu_ram_gb
  disk_gb          = local.hardware.disk_gb
  max_hourly_price = local.hardware.max_hourly_price
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
