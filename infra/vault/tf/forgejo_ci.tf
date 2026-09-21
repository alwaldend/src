module "forgejo_ci_repositories" {
  source = "../../repos/tf"
}

locals {
  forgejo_ci_repository = module.forgejo_ci_repositories.forgejo_repositories["alwaldend/src"]
  forgejo_ci_owner_repo = "${local.forgejo_ci_repository.organization}/${local.forgejo_ci_repository.name}"
  forgejo_ci_issuer     = "${var.forgejo_url}/api/actions"
  forgejo_ci_refs = [
    "refs/heads/${local.forgejo_ci_repository.default_branch}",
  ]
}

resource "vault_jwt_auth_backend" "forgejo_ci" {
  path                  = "forgejo_ci"
  description           = "Short-lived Forgejo Actions authentication"
  oidc_discovery_url    = local.forgejo_ci_issuer
  oidc_discovery_ca_pem = file("${path.module}/../../../data/ssl/alwaldend.com/ica1.crt")
  bound_issuer          = local.forgejo_ci_issuer
  jwt_supported_algs    = ["RS256"]
}

resource "vault_policy" "forgejo_ci_revoke_self" {
  name   = "forgejo_ci_revoke_self"
  policy = <<EOT
    path "auth/token/revoke-self" {
      capabilities = ["update"]
    }
EOT
}

resource "vault_jwt_auth_backend_role" "forgejo_ci_secure" {
  backend           = vault_jwt_auth_backend.forgejo_ci.path
  role_name         = "secure"
  role_type         = "jwt"
  user_claim        = "sub"
  bound_audiences   = [var.vault_url]
  bound_claims_type = "string"

  # Provider 5.8 splits comma-separated alternatives into Vault JSON arrays.
  disable_bound_claims_parsing = false

  # Forgejo 15.0.3 hardcodes ref_protected=false in services/actions/context.go.
  # Bind the default branch, covered by protection in
  # infra/forgejo/tf. This restricts Vault authentication, not runner scheduling.
  bound_claims = {
    repository = local.forgejo_ci_owner_repo
    ref        = join(",", local.forgejo_ci_refs)
    ref_type   = "branch"
    event_name = "push,workflow_dispatch"
    sub = join(",", [
      for ref in local.forgejo_ci_refs : "repo:${local.forgejo_ci_owner_repo}:ref:${ref}"
    ])
    workflow_ref = join(",", [
      for ref in local.forgejo_ci_refs : "${local.forgejo_ci_owner_repo}/.forgejo/workflows/secure.yaml@${ref}"
    ])
  }

  claim_mappings = {
    repository = "repository"
    ref        = "ref"
    run_id     = "run_id"
    sha        = "sha"
  }

  token_policies = [
    vault_policy.auth_token_lookup_self.name,
    vault_policy.forgejo_ci_revoke_self.name,
  ]
  token_no_default_policy = true
  token_ttl               = 300
  token_max_ttl           = 300
  token_explicit_max_ttl  = 300
}
