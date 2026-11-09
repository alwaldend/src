module "src_infra_vastai_approle" {
  source = "../../../projects/tf_modules/vault_approle"
  name   = "src_infra_vastai"
  policies = [
    vault_policy.vastai_imagegen_delivery.name,
  ]
  member_entity_ids = [
    vault_identity_entity.simeonwarren.id,
  ]
  secrets                  = vault_mount.secrets.path
  backend                  = vault_auth_backend.approle.path
  backend_accessor         = vault_auth_backend.approle.accessor
  disable_yc_folder_policy = true
}

resource "vault_policy" "vastai_imagegen_delivery" {
  name   = "src_infra_vastai_imagegen_delivery"
  policy = <<-EOT
    path "${vault_mount.secrets.path}/data/alwaldend.com/vault1/approles/user_simeonwarren/imagegen/vastai" {
      capabilities = ["create", "read", "update", "delete"]
    }
    path "${vault_mount.secrets.path}/metadata/alwaldend.com/vault1/approles/user_simeonwarren/imagegen/vastai" {
      capabilities = ["read"]
    }
  EOT
}
