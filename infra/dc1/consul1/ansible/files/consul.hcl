datacenter = "dc1"
data_dir = "{{ consul_data_dir }}"
node_name = "{{ inventory_hostname }}"
encrypt = "{{ lookup('env', 'AL_CONSUL_GOSSIP_KEY') }}"
log_level = "INFO"
domain = "consul"
server = true
bootstrap_expect = 1

ui_config {
  enabled = true
}

tls {
  internal_rpc {
    verify_server_hostname = true
  }
  defaults {
    verify_incoming = true
    verify_outgoing = true
    verify_server_hostname = true
    ca_file = "{{ consul_tls_ca_file_path }}"
    cert_file = "{{ consul_tls_cert_file_path }}"
    key_file = "{{ consul_tls_key_file_path }}"
  }
}

auto_encrypt {
  allow_tls = true
}
