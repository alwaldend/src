---
title: Repositories
description: Shared organization, repository, and access configuration
---

This project owns the organization and repository catalog consumed by
[GitHub](../github/tf/README.md), [GitLab](../gitlab/tf/README.md), and
[Forgejo](../forgejo/tf/README.md). Named organization administrators and
developers belong in this catalog. [Root configuration](config.json) owns
the schema version and defaults. Each organization owns
`orgs/<organization>/org.json` and one
`orgs/<organization>/repos/<repository>.json` file per repository. Vault
continues to own Forgejo login identities and service-specific access.

First-party means owned by us, regardless of forge. Our repositories retain
their chosen names across GitHub, GitLab, and Forgejo: `alwaldend/src` stays
`alwaldend/src`, including when copied or synchronized between those forges.

Forks and mirrors of external repositories must belong to an organization in
the catalog. Their names use the original upstream hostname in reverse order,
followed by its full owner and repository path, in lowercase with punctuation
replaced by underscores. For example, `https://gitlab.com/fdroid/fdroiddata`
becomes `alwaldend/com_gitlab_fdroid_fdroiddata`. Copies on another forge keep
that name; they do not add a prefix for our intermediate copy. This naming
rule also applies to external repositories imported once for later syncing.

The [Terraform module](tf/README.md) owns this derivation. A forge key opts a
repository into that consumer. GitLab `import_from` refers to the corresponding
catalog source; `fork_from` identifies a GitLab upstream. External origins
belong in the repository's `upstream_url`. Forge-specific overrides preserve
settings such as Pages sites and repository descriptions during adoption.

Consumers retain their own providers, state, and Vault authentication. They
adopt existing remote resources by import and preserve managed resource
addresses or declare explicit moves. Named developers can read repositories,
write feature branches, and open pull or merge requests; protected default
branches restrict pushes and merges to administrators and any separately
owned existing service grants. Catalog GitHub repositories permit merge
commits only: the shared defaults disable squash merges and rebase merges, and
a catalog precondition rejects any repository record that enables them. Merging
therefore keeps the reviewed feature commit as a parent of the default-branch
commit instead of replacing it. GitLab and Forgejo merge methods are not owned
here. Landing repositories use `master` as their protected default branch while
their Pages source remains `pages`, allowing the developer's existing Pages
deployment access to continue.

Catalog targets are repository-internal infrastructure inputs. They are not
production build dependencies or published artifacts.

GitLab receives one-time repository imports and the configured upstream fork.
The user will arrange ongoing synchronization separately; the decision and
adoption evidence live in [OpenSpec](openspec/changes/archive/2026-09-13-adopt-shared-repository-catalog/design.md).
