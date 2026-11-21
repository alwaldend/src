terraform {
  backend "http" {}
}

module "repositories" {
  source = "../../repos/tf"
}

locals {
  push_mirror_tool = abspath("${path.module}/../cmd/push_mirror/push_mirror")
  destination_repositories = {
    github = module.repositories.github_repositories
    gitlab = module.repositories.gitlab_repositories
  }
  mirror_vault_prefix = "alwaldend.com/vault1/approles/src_infra_dc1_forgejo1/push_mirrors"
}

resource "terraform_data" "push_mirror" {
  for_each = module.repositories.forgejo_push_mirrors

  triggers_replace = "${try(module.repositories.forgejo_repositories[each.value.repository].default_branch, "")},releases/*"

  lifecycle {
    precondition {
      condition     = can(module.repositories.forgejo_repositories[each.value.repository]) && can(local.destination_repositories[each.value.forge][each.value.repository])
      error_message = "The catalog must declare a corresponding GitHub or GitLab destination for ${each.key}."
    }

    precondition {
      condition     = each.key == "${each.value.repository}:${each.value.forge}" && trimspace(each.value.username) != ""
      error_message = "Use the repository:forge key and specify a mirror username for ${each.key}."
    }
  }

  provisioner "local-exec" {
    interpreter = [local.push_mirror_tool]
    command     = "ensure"
    environment = {
      FORGEJO_MIRROR_REPOSITORY    = "${module.repositories.forgejo_repositories[each.value.repository].organization}/${module.repositories.forgejo_repositories[each.value.repository].name}"
      FORGEJO_MIRROR_DESTINATION   = each.value.forge
      FORGEJO_MIRROR_TARGET        = "${local.destination_repositories[each.value.forge][each.value.repository].organization}/${local.destination_repositories[each.value.forge][each.value.repository].name}"
      FORGEJO_MIRROR_BRANCH_FILTER = self.triggers_replace
      FORGEJO_MIRROR_USERNAME      = each.value.username
      FORGEJO_MIRROR_VAULT_PATH    = "${local.mirror_vault_prefix}/${each.value.vault_key}"
    }
  }
}
