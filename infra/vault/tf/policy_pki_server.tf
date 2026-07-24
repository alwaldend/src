resource "vault_policy" "pki_server" {
  name   = "pki_server"
  policy = <<EOT
    path "${module.pki_ica_servers.backend}/issue/{{ identity.entity.name }}_pki_server" {
      capabilities = ["update"]
    }
EOT
}
