resource "forgejo_organization" "alwaldend" {
  name       = local.organization.name
  visibility = local.organization.forgejo.visibility

  lifecycle {
    prevent_destroy = true
  }
}
