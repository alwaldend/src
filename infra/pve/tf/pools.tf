resource "proxmox_pool" "approles" {
  for_each = local.approles
  poolid   = each.key
  comment  = "Pool for '${each.key}' approle, managed by terraform (//infra/dc1/pve1)"
}

resource "proxmox_pool" "templates" {
  poolid  = "templates"
  comment = "Template pool. Managed by terraform (//infra/dc1/pve1)"
}
