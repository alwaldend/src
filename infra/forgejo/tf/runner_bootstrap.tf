resource "terraform_data" "runners" {
  for_each = local.runner_registrations

  provisioner "local-exec" {
    interpreter = [local.runner_registration_tool]
    command     = "ensure"
    environment = {
      FORGEJO_RUNNER_REPOSITORY  = each.value.repository
      FORGEJO_RUNNER_NAME        = each.value.name
      FORGEJO_RUNNER_VAULT_MOUNT = each.value.mount
      FORGEJO_RUNNER_VAULT_PATH  = each.value.path
    }
  }
}
