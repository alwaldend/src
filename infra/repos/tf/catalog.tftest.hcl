run "catalog_contract" {
  command = plan

  assert {
    condition     = output.github_repositories["alwaldend/src"].name == "src"
    error_message = "First-party source names must remain unchanged."
  }
  assert {
    condition     = output.gitlab_repositories["alwaldend/src"].name == "src"
    error_message = "First-party ownership is independent of forge: src must remain src on GitLab."
  }
  assert {
    condition     = output.gitlab_repositories["alwaldend/com_gitlab_fdroid_fdroiddata"].name == "com_gitlab_fdroid_fdroiddata"
    error_message = "The F-Droid fork must use its upstream reverse hostname and path."
  }
  assert {
    condition     = output.github_repositories["alwaldend/com_github_infinitimeorg_infinitime"].name == "com_github_infinitimeorg_infinitime"
    error_message = "Upstream capitalization must normalize without renaming existing forks."
  }
  assert {
    condition     = output.gitlab_repositories["alwaldend/com_github_infinitimeorg_infinitime"].name == output.github_repositories["alwaldend/com_github_infinitimeorg_infinitime"].name
    error_message = "Copies of external forks must retain the original upstream name across forges."
  }
  assert {
    condition = alltrue([
      for key, repository in output.github_repositories : output.gitlab_repositories[key].config.import_url == "https://github.com/${repository.organization}/${repository.name}.git"
    ])
    error_message = "Every GitHub repository must retain a matching one-time GitLab import."
  }
  assert {
    condition = alltrue([
      for repository in output.github_repositories : (
        repository.config.allow_merge_commit &&
        !repository.config.allow_squash_merge &&
        !repository.config.allow_rebase_merge
      )
    ])
    error_message = "Every GitHub repository must permit merge commits only."
  }
  assert {
    condition     = !output.github_repositories["alwaldend/src"].config.allow_squash_merge && !output.github_repositories["alwaldend/src"].config.allow_rebase_merge
    error_message = "src must inherit the shared merge-commit-only policy rather than override it."
  }
  assert {
    condition = (
      output.forgejo_repositories["alwaldend/com_github_actions_checkout"].name == "com_github_actions_checkout" &&
      output.forgejo_repositories["alwaldend/com_github_actions_checkout"].config.clone_addr == "https://github.com/actions/checkout" &&
      output.forgejo_repositories["alwaldend/com_github_actions_checkout"].config.mirror &&
      output.forgejo_repositories["alwaldend/com_github_actions_checkout"].config.mirror_interval == "12h0m0s" &&
      !output.forgejo_repositories["alwaldend/com_github_actions_checkout"].config.has_actions
    )
    error_message = "Checkout must be a Forgejo pull mirror of its original upstream without running upstream workflows."
  }
  assert {
    condition = (
      !contains(keys(output.github_repositories), "alwaldend/com_github_actions_checkout") &&
      !contains(keys(output.gitlab_repositories), "alwaldend/com_github_actions_checkout") &&
      output.forgejo_repositories["alwaldend/src"].config.clone_addr == "https://github.com/alwaldend/src.git"
    )
    error_message = "The checkout mirror must not create other forge copies or change the source repository clone origin."
  }
}
