locals {
  dns      = jsondecode(file("${path.module}/../dnsconfig.json")).records
  hostname = "${local.dns.host1.A.name}.alwaldend.com"
  address  = local.dns.host1.A.address
  mac      = join(":", concat(["02"], [for i in range(5) : substr(sha256(local.hostname), i * 2, 2)]))
}

data "xenorchestra_resource_set" "forgejo" {
  name = "src_infra_dc1_forgejo1"
}

resource "xenorchestra_vm" "forgejo" {
  name_label       = local.hostname
  name_description = "Forgejo managed by infra/forgejo/tf_setup"
  template         = data.xenorchestra_template.forgejo.id
  resource_set     = data.xenorchestra_resource_set.forgejo.id
  cpus             = 2
  memory_max       = 4 * 1024 * 1024 * 1024
  memory_min       = 4 * 1024 * 1024 * 1024
  auto_poweron     = true
  power_state      = "Running"
  tags             = ["forgejo", "src_infra_dc1_forgejo1"]

  cloud_config = "#cloud-config\n${yamlencode(merge(
    jsondecode(file("${path.module}/../../cloud_init/xen_linux.json")),
    { hostname = "host1", fqdn = local.hostname },
  ))}"

  cloud_network_config = yamlencode({
    version = 1
    config = [{
      type        = "physical"
      name        = var.xoa_guest_interface
      mac_address = local.mac
      subnets = [{
        type            = "static"
        address         = "${local.address}/24"
        gateway         = var.gateway
        dns_nameservers = var.dns_servers
      }]
    }]
  })

  network {
    network_id  = data.xenorchestra_network.forgejo.id
    mac_address = local.mac
  }

  # A single-disk cloud template maps its boot disk to xvda; additional
  # disks are xvdb (Forgejo) and xvdc (Traefik) in the Xen Linux guest.
  disk {
    sr_id      = data.xenorchestra_sr.forgejo.id
    name_label = "${local.hostname} boot"
    size       = 20 * 1024 * 1024 * 1024
  }

  disk {
    sr_id      = data.xenorchestra_sr.forgejo.id
    name_label = "${local.hostname} forgejo"
    size       = 40 * 1024 * 1024 * 1024
  }

  disk {
    sr_id      = data.xenorchestra_sr.forgejo.id
    name_label = "${local.hostname} traefik"
    size       = 5 * 1024 * 1024 * 1024
  }

  timeouts {
    create = "20m"
  }
}
