ui = true
disable_mlock = false
log_level = "info"
log_file = "/opt/vault/logs/vault.log"

listener "tcp" {
  address            = "[::]:8200"
  tls_cert_file      = "/opt/vault/tls/tls_cert_file.pem"
  tls_key_file       = "/opt/vault/tls/tls_key_file.pem"
  # tls_client_ca_file = "/opt/vault/tls/tls_client_ca_file.pem"
}

user_lockout "all" {
  lockout_duration = "10m"
  lockout_threshold = "25"
  lockout_counter_reset = "10m"
}

storage "file" {
  path = "/opt/vault/data"
}
