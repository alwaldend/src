resource "truenas_smb_config" "transport" {
  encryption = "REQUIRED"
}

resource "truenas_service" "smb" {
  name    = "cifs"
  enabled = true
  running = true

  depends_on = [truenas_smb_config.transport]
}
