resource "proxmox_vm_qemu" "cloudinti_test" {
  vmid        = 100
  name        = "cloudinit-test.vm.pve1.dc1.alwaldend.com"
  target_node = "host1"
  agent       = 1
  cpu {
    cores = 2
  }
  memory           = 1024
  clone            = "fedora-cloud.templates.pve1.dc1.alwaldend.com"
  boot             = "order=scsi0"
  scsihw           = "virtio-scsi-single"
  vm_state         = "running"
  automatic_reboot = true
  cicustom         = "vendor=local:snippets/cloud_init.yaml"
  ipconfig0        = "ip=192.168.1.10/24,gw=192.168.1.1"
  disks {
    scsi {
      scsi0 {
        disk {
          storage = "local-lvm"
          size    = "20G"
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
}
