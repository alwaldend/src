# Public snapshots are refreshed by the data-owned Terraform roots.
# Relinquish the former exports without deleting files on the next apply.
removed {
  from = local_file.pki_ca_clients
  lifecycle {
    destroy = false
  }
}

removed {
  from = local_file.pki_ca_servers
  lifecycle {
    destroy = false
  }
}

removed {
  from = local_file.ssh_ca_clients
  lifecycle {
    destroy = false
  }
}

removed {
  from = local_file.ssh_ca_servers
  lifecycle {
    destroy = false
  }
}
