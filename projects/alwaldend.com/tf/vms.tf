resource "proxmox_vm_qemu" "vms" {
  for_each = {
    host1 = { vmid = 900 },
  }
  name   = "com.alwaldend.www.pve.${each.key}"
  vmid   = each.value.vmid
  target_nodes = ["host1"]
  pool   = "src_projects_alwaldend_com"
  cpu {
    cores = 1
  }
  memory = 1024
  balloon = 512
  tags = "src_projects_alwaldend_com,unikernel"
  scsihw             = "virtio-scsi-single"
  vm_state           = "running"
  automatic_reboot   = true
  start_at_node_boot = true
  boot = "order=scsi0"
  serial {
    id = 0
    type = "socket"
  }
  vga {
    type = "serial0"
  }
  agent = 0
  startup_shutdown {
    order            = -1
    shutdown_timeout = -1
    startup_delay    = -1
  }
  network {
    id     = 0
    bridge = "vmbr0"
    model  = "virtio"
  }
  disks {
    scsi {
      scsi0 {
          disk {
              backup               = true
              discard              = true
              emulatessd           = true
              format               = "qcow2"
              size                 = "1G"
              storage              = "local-lvm"
            }
        }
    }
  }
}
