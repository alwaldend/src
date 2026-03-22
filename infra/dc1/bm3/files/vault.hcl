ui = true
disable_mlock = false
log_level = "info"
log_file = "/opt/vault/logs/vault.log"
api_addr = "http://127.0.0.1:8200"
cluster_addr = "http://127.0.0.1:8201"

listener "tcp" {
  address            = "127.0.0.1:8200"
  tls_disable = true
  # tls_cert_file      = "/opt/vault/tls/vault-cert.pem"
  # tls_key_file       = "/opt/vault/tls/vault-key.pem"
  # tls_client_ca_file = "/opt/vault/tls/vault-ca.pem"
}

user_lockout "all" {
  lockout_duration = "10m"
  lockout_threshold = "25"
  lockout_counter_reset = "10m"
}

storage "raft" {
  path    = "/opt/vault/data"
  node_id = "vault.bm3.dc1.alwaldend.com"
}
