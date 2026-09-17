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
      for repository in output.github_repositories : repository.default_branch == "master" && repository.config.pages.source.branch == "pages"
      if can(repository.config.landing_project)
    ])
    error_message = "Landing defaults must be master while publication remains on pages."
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
}
