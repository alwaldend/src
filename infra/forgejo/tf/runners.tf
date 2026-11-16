locals {
  runner_registration_tool = abspath("${path.module}/../cmd/runner_registration/runner_registration")
  runner_registrations = {
    secure1 = {
      # Registration consumes an existing repository; use catalog metadata so
      # targeted bootstrap does not reconcile unrelated repository resources.
      repository = "${module.repositories.forgejo_repositories[local.src_repository_key].organization}/${module.repositories.forgejo_repositories[local.src_repository_key].name}"
      name       = "secure1.dev.alwaldend.com"
      mount      = "secrets"
      path       = "alwaldend.com/vault1/approles/src_infra_dc1_forgejo1/runners/secure1"
    }
  }
}
