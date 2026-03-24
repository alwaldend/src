ui = true
disable_mlock = false
log_level = "info"
log_file = "/opt/vault/logs/vault.log"

listener "tcp" {
  address            = "[::]:8200"
  tls_disable = true
  # tls_cert_file      = "/opt/vault/tls/vault_cert.pem"
  # tls_key_file       = "/opt/vault/tls/vault_key.pem"
  # tls_client_ca_file = "/opt/vault/tls/vault_ca.pem"
}

user_lockout "all" {
  lockout_duration = "10m"
  lockout_threshold = "25"
  lockout_counter_reset = "10m"
}

storage "file" {
  path = "/opt/vault/data"
}
