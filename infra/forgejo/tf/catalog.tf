module "repositories" {
  source = "../../repos/tf"
}

locals {
  organization = one([
    for name, organization in module.repositories.organizations : merge(organization, { name = name })
    if can(organization.forgejo)
  ])
  src_repository_key = "${local.organization.name}/src"
}
