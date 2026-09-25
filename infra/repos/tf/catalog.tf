locals {
  catalog = jsondecode(file("${path.module}/../config.json"))
  organizations = {
    for org_file in fileset("${path.module}/../orgs", "*/org.json") : split("/", org_file)[0] => merge(
      jsondecode(file("${path.module}/../orgs/${org_file}")),
      {
        repositories = {
          for repo_file in fileset("${path.module}/../orgs/${dirname(org_file)}/repos", "*.json") : trimsuffix(repo_file, ".json") => jsondecode(file("${path.module}/../orgs/${dirname(org_file)}/repos/${repo_file}"))
        }
      },
    )
  }

  # A forge key opts a repository into that consumer. Provider settings stay
  # with the forge key; the common identity and access lists have one owner.
  repositories = flatten([
    for organization, org in local.organizations : [
      for name, repository in org.repositories : [
        for forge in ["github", "gitlab", "forgejo"] : {
          organization   = organization
          catalog_key    = name
          forge          = forge
          description    = try(repository[forge].description, repository.description)
          homepage_url   = try(repository[forge].homepage_url, repository.homepage_url)
          default_branch = repository.default_branch
          visibility     = repository.visibility
          config         = merge(try(local.catalog.defaults[forge], {}), repository[forge])
          upstream_url = try(repository.upstream_url, try(repository.gitlab.fork_from, null) != null ? (
            "https://gitlab.com/${repository.gitlab.fork_from}"
          ) : null)
        } if can(repository[forge])
      ]
    ]
  ])

  # Reverse DNS followed by the complete upstream path, lowercase ASCII with
  # punctuation normalized to underscores. Ownership determines first-party
  # status: copies of our repositories retain their names on every forge.
  upstream_parts = {
    for repository in local.repositories : "${repository.forge}/${repository.organization}/${repository.catalog_key}" => (
      repository.upstream_url == null ? [] : split("/", trimprefix(repository.upstream_url, "https://"))
    )
  }
  projected = {
    for repository in local.repositories : "${repository.forge}/${repository.organization}/${repository.catalog_key}" => merge(repository, {
      name = repository.upstream_url == null ? repository.catalog_key : replace(lower(join("_", concat(
        reverse(split(".", local.upstream_parts["${repository.forge}/${repository.organization}/${repository.catalog_key}"][0])),
        slice(local.upstream_parts["${repository.forge}/${repository.organization}/${repository.catalog_key}"], 1, length(local.upstream_parts["${repository.forge}/${repository.organization}/${repository.catalog_key}"])),
      ))), "/[^a-z0-9_]+/", "_")
      config = merge(repository.config, repository.forge == "gitlab" ? {
        import_url = try(repository.config.import_from, null) == "github" ? "https://github.com/${repository.organization}/${repository.catalog_key}.git" : null
        fork_from  = try(repository.config.fork_from, null)
        } : repository.forge == "forgejo" ? {
        clone_addr = try(repository.config.clone_from, null) == "upstream" ? repository.upstream_url : (
          try(repository.config.clone_from, null) == "github" ? "https://github.com/${repository.organization}/${repository.catalog_key}.git" : null
        )
      } : {})
    })
  }
}

output "organizations" {
  value = local.organizations

  precondition {
    condition     = local.catalog.version == 1
    error_message = "Unsupported repository catalog version."
  }
  precondition {
    condition = alltrue([
      for org in values(local.organizations) : length(setintersection(toset(org.admins), toset(org.users))) == 0
    ])
    error_message = "Organization administrators and developers must be disjoint."
  }
  precondition {
    condition = alltrue([
      for repository in local.repositories : repository.upstream_url == null ? true : can(regex("^https://[a-zA-Z0-9.-]+/[^?#@]+$", repository.upstream_url))
    ])
    error_message = "Fork and mirror origins must be public HTTPS repository URLs without credentials, queries, or fragments."
  }
  precondition {
    condition = alltrue([
      for repository in values(local.projected) : repository.forge != "forgejo" ? true : (
        (try(repository.config.clone_from, null) == null ? true : contains(["github", "upstream"], repository.config.clone_from)) &&
        (try(repository.config.clone_from, null) != "upstream" || repository.upstream_url != null) &&
        (!try(repository.config.mirror, false) || repository.config.clone_addr != null)
      )
    ])
    error_message = "Forgejo clone_from must select github or upstream; upstream requires upstream_url and pull mirrors require a clone source."
  }
  precondition {
    condition = length(distinct([
      for repository in values(local.projected) : "${repository.forge}/${repository.organization}/${repository.name}"
    ])) == length(local.projected)
    error_message = "Normalized repository names must be unique within each organization and forge."
  }
  precondition {
    condition = alltrue([
      for repository in local.repositories : repository.forge != "github" ? true : (
        try(repository.config.allow_merge_commit, false) &&
        !try(repository.config.allow_squash_merge, false) &&
        !try(repository.config.allow_rebase_merge, false)
      )
    ])
    error_message = format(
      "GitHub repositories must permit merge commits only; squash and rebase merges are disabled. Offending records: %s",
      join(", ", [
        for repository in local.repositories : "${repository.organization}/${repository.catalog_key}"
        if repository.forge == "github" && (
          !try(repository.config.allow_merge_commit, false) ||
          try(repository.config.allow_squash_merge, false) ||
          try(repository.config.allow_rebase_merge, false)
        )
      ])
    )
  }
}

output "github_repositories" {
  value = {
    for repository in values(local.projected) : "${repository.organization}/${repository.catalog_key}" => repository
    if repository.forge == "github"
  }
}

output "gitlab_repositories" {
  value = {
    for repository in values(local.projected) : "${repository.organization}/${repository.catalog_key}" => repository
    if repository.forge == "gitlab"
  }
}

output "forgejo_repositories" {
  value = {
    for repository in values(local.projected) : "${repository.organization}/${repository.catalog_key}" => repository
    if repository.forge == "forgejo"
  }
}
