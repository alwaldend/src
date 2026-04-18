resource "vault_auth_backend" "cert" {
    path = "cert"
    type = "cert"
}

/*
resource "vault_cert_auth_backend_role" "cert" {
    name           = "foo"
    certificate    = file("/path/to/certs/ca-cert.pem")
    backend        = vault_auth_backend.cert.path
    allowed_names  = ["foo.example.org", "baz.example.org"]
    token_ttl      = 300
    token_max_ttl  = 600
    token_policies = ["foo"]
}
*/
