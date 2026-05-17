resource "proxmox_vm_qemu" "cloudinti_test" {
  vmid        = 100
  name        = "cloudinit-test.vm.pve1.dc1.alwaldend.com"
  target_node = "host1"
  cpu {
    cores = 1
  }
  agent            = 1
  memory           = 1024
  tags             = join(",", ["src_infra_dc1_pve1", "test"])
  clone            = "fedora-cloud.templates.pve1.dc1.alwaldend.com"
  scsihw           = "virtio-scsi-single"
  vm_state         = "running"
  automatic_reboot = true
  cicustom         = "user=local:snippets/cloud_init.yaml"
  ipconfig0        = "ip=192.168.10.10/24,gw=192.168.1.1"
  disks {
    scsi {
      scsi0 {
        disk {
          storage = "local-lvm"
          size    = "5G"
        }
      }
    }
    ide {
      ide1 {
        cloudinit {
          storage = "local-lvm"
        }
      }
    }
  }
  network {
    id     = 0
    bridge = "vmbr0"
    model  = "virtio"
  }
}
