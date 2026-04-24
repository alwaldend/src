resource "vault_auth_backend" "userpass" {
  path = "userpass"
  type = "userpass"
  description = "Userpass auth"

  tune {
    default_lease_ttl = "24h"
    max_lease_ttl      = "24h"
    listing_visibility = "unauth"
  }
}
resource "vault_auth_backend" "approle" {
  path = "approle"
  type = "approle"
  description = "Approle auth"
}
