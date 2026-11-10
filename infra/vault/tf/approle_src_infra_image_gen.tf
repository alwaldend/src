module "src_infra_image_gen_approle" {
  source = "../../../projects/tf_modules/vault_approle"
  name   = "src_infra_image_gen"
  member_entity_ids = [
    vault_identity_entity.simeonwarren.id,
  ]
  secrets                    = vault_mount.secrets.path
  backend                    = vault_auth_backend.approle.path
  backend_accessor           = vault_auth_backend.approle.accessor
  disable_yc_folder_policy   = true
  allow_member_secret_access = true
  policies                   = [vault_policy.image_gen_web_lifecycle.name]
}

# The role's standard namespace policy permits soft deletion only. This root
# owns the web login and needs metadata deletion for delete_all_versions=true.
# Exact-path capabilities replace the wildcard grant, so retain metadata reads.
resource "vault_policy" "image_gen_web_lifecycle" {
  name   = "src_infra_image_gen_web_lifecycle"
  policy = <<-EOT
    path "${vault_mount.secrets.path}/metadata/alwaldend.com/vault1/approles/src_infra_image_gen/web" {
      capabilities = ["read", "delete"]
    }
  EOT
}
