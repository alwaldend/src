output "ica_clients_certificate" {
  description = "Client certificate CA"
  value       = module.pki_ica_clients.certificate
}


output "ssh_ca_clients_certificate" {
  description = "Public key for the ssh CA for clients"
  value       = vault_ssh_secret_backend_ca.clients.public_key
}

output "ssh_ca_servers_certificate" {
  description = "Public key for the ssh CA for servers"
  value       = vault_ssh_secret_backend_ca.servers.public_key
}
