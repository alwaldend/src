ui = true
disable_mlock = false
log_level = "info"
log_file = "/opt/vault/logs/vault.log"
api_addr = "https://vault.dc1.alwaldend.com:8200"
cluster_addr = "https://vault.dc1.alwaldend.com:8201"

listener "tcp" {
  address            = "0.0.0.0:8200"
  tls_cert_file      = "/opt/vault/tls/vault_cert.pem"
  tls_key_file       = "/opt/vault/tls/vault_key.pem"
  tls_client_ca_file = "/opt/vault/tls/vault_ca.pem"
}

user_lockout "all" {
  lockout_duration = "10m"
  lockout_threshold = "25"
  lockout_counter_reset = "10m"
}

storage "raft" {
  path    = "/opt/vault/data"
  node_id = "bm3.dc1.alwaldend.com"
}
