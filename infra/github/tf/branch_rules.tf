locals {
  src_protected_refs = ["~DEFAULT_BRANCH", "refs/heads/releases/**/*"]
}

# Preserve the two existing src rulesets independently, including their
# additional release-branch coverage and previously configured bypass actors.
resource "github_repository_ruleset" "src_main" {
  name        = "main"
  repository  = github_repository.other["src"].name
  target      = "branch"
  enforcement = "active"

  conditions {
    ref_name {
      include = local.src_protected_refs
      exclude = []
    }
  }

  bypass_actors {
    actor_id    = 5
    actor_type  = "RepositoryRole"
    bypass_mode = "always"
  }

  rules {
    deletion         = true
    non_fast_forward = true
  }

  lifecycle {
    prevent_destroy = true
  }
}

resource "github_repository_ruleset" "src_protected" {
  name        = "protected"
  repository  = github_repository.other["src"].name
  target      = "branch"
  enforcement = "active"

  conditions {
    ref_name {
      include = local.src_protected_refs
      exclude = []
    }
  }

  bypass_actors {
    # GitHub returns null for this actor; the provider represents that as zero.
    actor_id    = 0
    actor_type  = "OrganizationAdmin"
    bypass_mode = "always"
  }

  dynamic "bypass_actors" {
    for_each = [2, 5]
    content {
      actor_id    = bypass_actors.value
      actor_type  = "RepositoryRole"
      bypass_mode = "always"
    }
  }

  rules {
    creation            = true
    update              = true
    deletion            = true
    non_fast_forward    = true
    required_signatures = true
  }

  lifecycle {
    prevent_destroy = true
  }
}

# Repository rulesets apply to public repositories on GitHub Free.
# Wait for default-branch changes so landing publication branches stay writable.
resource "github_repository_ruleset" "default_branch" {
  for_each = local.default_branch_repositories

  name        = "default_branch_admins"
  repository  = local.managed_repository_names[each.key]
  target      = "branch"
  enforcement = "active"

  conditions {
    ref_name {
      include = ["~DEFAULT_BRANCH"]
      exclude = []
    }
  }

  bypass_actors {
    # Zero serializes as GitHub's null OrganizationAdmin actor ID.
    actor_id    = 0
    actor_type  = "OrganizationAdmin"
    bypass_mode = "always"
  }

  rules {
    creation                      = true
    update                        = true
    update_allows_fetch_and_merge = false
    deletion                      = true
    non_fast_forward              = true
  }

  lifecycle {
    prevent_destroy = true
  }

  depends_on = [github_branch_default.repository]
}
