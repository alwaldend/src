terraform {
  required_providers {
    proxmox = {
      source  = "Telmate/proxmox"
      version = ">= 3"
    }
  }
}

variable "name" {
  type        = string
  description = "VM name"
}

variable "vmid" {
  type        = number
  description = "VM id"
}

variable "ip" {
  type        = string
  description = "IP"
}

variable "gw" {
  type        = string
  description = "Gateway"
  default     = "192.168.10.2"
}

variable "ipconfig0" {
  type    = string
  default = ""
}

variable "cores" {
  type        = number
  description = "Cpu cores"
  default     = 1
}

variable "memory" {
  type        = number
  description = "Memory"
  default     = 1024
}

variable "pool" {
  type        = string
  description = "Pool id"
}

variable "target_nodes" {
  type        = list(string)
  description = "Target node"
  default     = ["host1"]
}

variable "tags" {
  type        = list(string)
  description = "Tags"
  default     = []
}

variable "clone" {
  type        = string
  description = "VM name to clone"
  default     = "fedora-cloud.templates.pve1.dc1.alwaldend.com"
}

variable "scsihw" {
  type    = string
  default = "virtio-scsi-single"
}

variable "vm_state" {
  type    = string
  default = "running"
}

variable "cicustom" {
  type    = string
  default = "user=local:snippets/cloud_init.yaml"
}

variable "storage" {
  type    = string
  default = "local-lvm"
}

variable "storage_size" {
  type    = string
  default = "10G"
}

variable "scsi1" {
  type = object({
    size       = string
    storage    = optional(string, "ceph-repl")
    asyncio    = optional(string, "native")
    iothread   = optional(bool, true)
    emulatessd = optional(bool, true)
    discard    = optional(bool, true)
  })
  default = null
}

variable "network_bridge" {
  type    = string
  default = "vmbr0"
}

variable "balloon" {
  type    = number
  default = 1024
}

variable "network_model" {
  type    = string
  default = "virtio"
}

resource "proxmox_vm_qemu" "vm" {
  vmid         = var.vmid
  name         = "${var.name}.alwaldend.com"
  target_nodes = var.target_nodes
  pool         = var.pool
  cpu {
    cores = var.cores
  }
  agent              = 1
  memory             = var.memory
  balloon            = var.balloon
  tags               = join(",", concat(var.tags, ["src_projects_tf_modules_pve_vm_qemu"]))
  clone              = var.clone
  scsihw             = var.scsihw
  vm_state           = var.vm_state
  automatic_reboot   = true
  start_at_node_boot = true
  cicustom           = var.cicustom
  ipconfig0          = var.ipconfig0 == "" ? "ip=${var.ip},gw=${var.gw}" : var.ipconfig0
  startup_shutdown {
    order            = -1
    shutdown_timeout = -1
    startup_delay    = -1
  }
  disks {
    scsi {
      scsi0 {
        disk {
          storage    = var.storage
          asyncio    = "native"
          iothread   = true
          emulatessd = true
          discard    = true
          size       = var.storage_size
        }
      }
      dynamic "scsi1" {
        for_each = var.scsi1 != null ? { disk = var.scsi1 } : {}
        content {
          disk {
            storage    = scsi1.value.storage
            asyncio    = scsi1.value.asyncio
            iothread   = scsi1.value.iothread
            emulatessd = scsi1.value.emulatessd
            discard    = scsi1.value.discard
            size       = scsi1.value.size
          }
        }
      }
    }
    ide {
      ide1 {
        cloudinit {
          storage = var.storage
        }
      }
    }
  }
  network {
    id     = 0
    bridge = var.network_bridge
    model  = var.network_model
  }
}
