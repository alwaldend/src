datacenter = "dc1"
data_dir = "{{ consul_data_dir }}"
node_name = "{{ inventory_hostname }}"
log_level = "INFO"
server = true
ui_config {
  enabled = true
}
addresses {
  http = "0.0.0.0"
}
tls {
  defaults {
    verify_incoming = true
    verify_outgoing = true
    verify_server_hostname = true
    ca_file = "{{ consul_tls_dir }}/consul-agent-ca.pem"
    cert_file = "{{ consul_tls_dir }}/dc1-server-consul-0.pem"
    key_file = "{{ consul_tls_dir }}/dc1-server-consul-0-key.pem"
  }
}
