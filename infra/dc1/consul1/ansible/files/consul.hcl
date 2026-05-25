datacenter = "dc1"
data_dir = "{{ consul_data_dir }}"
node_name = "{{ inventory_hostname | replace('.', '-') }}"
encrypt = "{{ lookup('env', 'AL_CONSUL_GOSSIP_KEY') | trim }}"
log_level = "INFO"
domain = "consul"
server = true
client_addr = "0.0.0.0"

retry_join = [
  {% for host in groups["consul1"] %}
  "{{ host }}",
  {% endfor %}
]

ui_config {
  enabled = true
}

acl {
  enabled = true
  default_policy = "deny"
  enable_token_persistence = true
}

ports {
  http = -1
  https = 8501
  grpc = -1
  grpc_tls = 8502
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
