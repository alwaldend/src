locals {
  runner_registration_tool = abspath("${path.module}/../cmd/runner_registration/runner_registration")
  runner_registrations = {
    secure1 = {
      repository = "${forgejo_repository.repositories[local.src_repository_key].owner}/${forgejo_repository.repositories[local.src_repository_key].name}"
      name       = "secure1.dev.alwaldend.com"
      mount      = "secrets"
      path       = "alwaldend.com/vault1/approles/src_infra_dc1_forgejo1/runners/secure1"
    }
  }
}
