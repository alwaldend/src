resource "truenas_app" "managed" {
  for_each = toset(["nextcloud"])

  name        = each.key
  catalog_app = "nextcloud"
  train       = "stable"
  version     = "2.3.60"
  running     = true
  custom_app  = false

  # The appliance retains its existing application configuration and secrets.
  lifecycle {
    prevent_destroy = true
  }

  depends_on = [truenas_dataset.managed]
}
