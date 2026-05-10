output "ica_clients_certificate" {
  description = "Client certificate CA"
  value       = module.pki_ica_clients.certificate
}
