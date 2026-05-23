resource "proxmox_pool" "approles" {
  for_each = local.approles
  poolid   = each.key
  comment  = "Managed by terraform (//infra/dc1/pve1)"
}
