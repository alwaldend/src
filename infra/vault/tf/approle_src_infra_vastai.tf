module "src_infra_vastai_approle" {
  source = "../../../projects/tf_modules/vault_approle"
  name   = "src_infra_vastai"
  member_entity_ids = [
    vault_identity_entity.simeonwarren.id,
  ]
  secrets                  = vault_mount.secrets.path
  backend                  = vault_auth_backend.approle.path
  backend_accessor         = vault_auth_backend.approle.accessor
  disable_yc_folder_policy = true
  policies                 = [vault_policy.vastai_image_gen_delivery.name]
}

resource "vault_policy" "vastai_image_gen_delivery" {
  name   = "src_infra_vastai_image_gen_delivery"
  policy = <<-EOT
    path "${vault_mount.secrets.path}/data/alwaldend.com/vault1/approles/src_infra_image_gen/vastai" {
      capabilities = ["create", "read", "update", "delete"]
    }
    path "${vault_mount.secrets.path}/metadata/alwaldend.com/vault1/approles/src_infra_image_gen/vastai" {
      capabilities = ["read", "delete"]
    }
  EOT
}
