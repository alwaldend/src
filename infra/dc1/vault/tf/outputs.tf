resource "local_file" "pki_ca_clients" {
  content  = module.pki_ica_clients.certificate
  filename = "${path.module}/output/pki_ca_clients.crt"
}

resource "local_file" "pki_ca_servers" {
  content  = module.pki_ica_servers.certificate
  filename = "${path.module}/output/pki_ca_servers.crt"
}

resource "local_file" "ssh_ca_clients" {
  content  = vault_ssh_secret_backend_ca.clients.public_key
  filename = "${path.module}/output/ssh_ca_clients.crt"
}

resource "local_file" "ssh_ca_servers" {
  content  = vault_ssh_secret_backend_ca.servers.public_key
  filename = "${path.module}/output/ssh_ca_servers.crt"
}

resource "local_file" "approles" {
  content  = jsonencode({ approles = local.approles })
  filename = "${path.module}/output/approles.json"
}
