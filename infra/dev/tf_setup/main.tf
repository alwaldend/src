terraform {
  required_providers {
    vault = {
      source  = "hashicorp/vault"
      version = "5.8.0"
    }
    xenorchestra = {
      source  = "vatesfr/xenorchestra"
      version = "0.41.0"
    }
  }
}

provider "vault" {}

provider "xenorchestra" {
  url = "wss://xoa.xcp-ng.alwaldend.com"
}

variable "ssh_private_key_file" {
  description = "Ephemeral private key file supplied by the Vault-backed injector."
  type        = string
  sensitive   = true
  default     = ""
}

variable "ssh_public_key_file" {
  description = "Matching public key file supplied by the Vault-backed injector."
  type        = string
  default     = ""
}

locals {
  secure1  = jsondecode(file("${path.module}/../dnsconfig.json")).records.secure1.A
  hostname = "${local.secure1.name}.alwaldend.com"
  mac      = join(":", concat(["02"], [for i in range(5) : substr(sha256(local.hostname), i * 2, 2)]))
}

data "xenorchestra_resource_set" "dev" { name = "src_infra_dev" }
data "xenorchestra_template" "secure1" { name_label = "fedora44-cloud.templates.xcp-ng.alwaldend.com" }
data "xenorchestra_sr" "secure1" { name_label = "Local storage" }
data "xenorchestra_sr" "runner_cache" { name_label = "Runner cache" }
data "xenorchestra_network" "secure1" { name_label = "Pool-wide network 1" }

resource "xenorchestra_vm" "secure1" {
  name_label       = local.hostname
  name_description = "Secure host managed by infra/dev"
  template         = data.xenorchestra_template.secure1.id
  resource_set     = data.xenorchestra_resource_set.dev.id
  cpus             = 8
  memory_max       = 16 * 1024 * 1024 * 1024
  memory_min       = 16 * 1024 * 1024 * 1024
  auto_poweron     = true
  power_state      = "Running"
  tags             = ["secure1", "src_infra_dev"]

  cloud_config = "#cloud-config\n${yamlencode(merge(
    jsondecode(file("${path.module}/../../cloud_init/xen_linux.json")),
    {
      hostname = "secure1", fqdn = local.hostname
      ssh_keys = {
        ed25519_private = sensitive(file(var.ssh_private_key_file))
        ed25519_public  = file(var.ssh_public_key_file)
      }
    },
  ))}"
  cloud_network_config = yamlencode({
    version = 1
    config = [{
      type = "physical", name = "enX0", mac_address = local.mac
      subnets = [{
        type    = "static", address = "${local.secure1.address}/24"
        gateway = "192.168.10.1", dns_nameservers = ["192.168.10.1"]
      }]
    }]
  })
  network {
    network_id  = data.xenorchestra_network.secure1.id
    mac_address = local.mac
  }
  disk {
    sr_id      = data.xenorchestra_sr.secure1.id
    name_label = "${local.hostname} boot"
    size       = 500 * 1024 * 1024 * 1024
  }
  disk {
    sr_id      = data.xenorchestra_sr.runner_cache.id
    name_label = "${local.hostname} cache"
    size       = 800 * 1024 * 1024 * 1024
  }
  lifecycle {
    # Cloud-init and its initial host key are consumed only at creation.
    ignore_changes = [cloud_config]
  }
  timeouts { create = "20m" }
}

output "id" { value = xenorchestra_vm.secure1.id }
output "address" { value = local.secure1.address }
output "hostname" { value = local.hostname }

output "ssh_host_public_key" {
  description = "Public bootstrap identity from the created VM's persisted cloud-init configuration."
  value       = nonsensitive(yamldecode(xenorchestra_vm.secure1.cloud_config).ssh_keys.ed25519_public)
}
