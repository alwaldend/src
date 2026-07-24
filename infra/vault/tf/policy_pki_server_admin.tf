resource "vault_policy" "pki_server_admin" {
  name   = "pki_server_admin"
  policy = <<EOT
    path "${module.pki_ica_servers.backend}/issue/{{ identity.entity.name }}_pki_server" {
      capabilities = ["create", "read", "update", "delete", "list", "sudo"]
    }
EOT
}
