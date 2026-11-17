resource "terraform_data" "commit_signing" {
  provisioner "local-exec" {
    interpreter = [local.signing_key_tool]
    command     = "ensure"
    environment = {
      VAULT_SSH_KEY_MOUNT = local.signing_key_mount
      VAULT_SSH_KEY_PATH  = local.signing_key_path
    }
  }
}
