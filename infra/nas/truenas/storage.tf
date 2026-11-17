resource "truenas_pool" "managed" {
  # Retain the address already recorded in Terraform state.
  for_each = toset(["default"])

  name          = each.key
  autotrim      = false
  force_destroy = false
  topology = {
    data = [{
      type  = "RAIDZ2"
      disks = ["Z0IPOVQF", "Z0IPOVQY", "Z0IM3TTS", "WD-WX22D245K3EU", "WD-WCA2F42V00F5"]
    }]
    cache = []
    log   = []
    spare = []
  }

  lifecycle {
    prevent_destroy = true
  }
}

resource "truenas_dataset" "managed" {
  for_each = {
    "default" = {
      compression        = "lz4"
      acltype            = "posix"
      aclmode            = "DISCARD"
      copies             = null
      inherit_encryption = false
      encryption         = false
    }
    "default/encrypted" = {
      atime              = null
      acltype            = null
      aclmode            = null
      inherit_encryption = false
      encryption         = true
    }
    "default/encrypted/apps"                    = {}
    "default/encrypted/apps/nextcloud"          = {}
    "default/encrypted/apps/nextcloud/pgdata"   = {}
    "default/encrypted/apps/nextcloud/appdata"  = {}
    "default/encrypted/apps/nextcloud/userdata" = {}
  }

  name        = each.key
  type        = "filesystem"
  compression = try(each.value.compression, null)
  # App datasets share these settings; explicit null overrides preserve inheritance.
  atime              = try(each.value.atime, "OFF")
  acltype            = try(each.value.acltype, "nfsv4")
  aclmode            = try(each.value.aclmode, "PASSTHROUGH")
  copies             = try(each.value.copies, 1)
  encryption         = try(each.value.encryption, null)
  inherit_encryption = try(each.value.inherit_encryption, true)
  force_destroy      = false

  lifecycle {
    prevent_destroy = true
  }

  depends_on = [truenas_pool.managed]
}
