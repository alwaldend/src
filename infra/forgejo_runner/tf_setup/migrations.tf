# The operator confirmed this Proxmox VM no longer exists. Forget its stale
# state without contacting or attempting to destroy the retired host.
removed {
  from = module.vm
  lifecycle {
    destroy = false
  }
}
