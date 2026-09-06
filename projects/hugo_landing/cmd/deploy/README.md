---
title: Landing deployment
description: Publish rendered project sites to GitHub Pages repositories
---

`//projects/hugo_landing/cmd/deploy:deploy` publishes rendered sites to
`git@github.com:<owner>/<project-slug>-landing.git` on the `pages` branch.
Project underscores become hyphens in repository names and custom domains.
The owner defaults to `alwaldend`. Terraform owns repository creation and
GitHub Pages settings; this command requires the repositories to exist.

Supply repeated `--site project=rendered/site` arguments, or `--manifest`
with a JSON object mapping project names to rendered site directories.
Site paths may be absolute or relative to the process working directory,
including Bazel runfiles paths. A symlink at the site root is resolved;
symlinks inside the site remain unsupported. Paths in a manifest are relative
to the working directory, rather than the manifest's parent directory.
`--scratch` must be an absolute task-owned directory.
Use `--project project_name` to select one entry from the supplied sites;
an unknown project fails before any deployment.

The command validates every site before publishing, creates a fresh clone
under scratch for each project, replaces the tracked site files, and writes
`CNAME` and `.nojekyll`. It bootstraps an absent `pages` branch without changing
other branches. Unchanged sites produce no commit. Updates preserve history
and use normal pushes, so concurrent branch changes cause a push rejection.
Successful and unchanged receipts include the exact commit revision.
Temporary clones are removed after success or failure.

`--dry-run` clones and stages the rendered sites and reports whether each
would change, without committing or pushing. It still reads the remote.

Deployment runs outside build actions and intentionally requires network
access, a host Git executable (`--git`, default `git`), and existing SSH GitHub
authentication. Credentials stay in the existing Git/SSH authentication flow;
subprocess output is excluded from errors to avoid exposing credentials.
Git operations have a two-minute timeout and interactive Git prompts are
disabled. Integration tests require local host Git and use only temporary
local repositories, without network access.
