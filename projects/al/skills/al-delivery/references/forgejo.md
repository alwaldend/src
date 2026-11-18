# Forgejo compatibility workflow

Except for its required read-only `provider` command, do not invoke
`al delivery` for a Forgejo remote until it has a Forgejo adapter. Use
`$al-rebase` for exact ref discovery, synchronization, rebasing, and
an exact-lease push. Use
the repository-pinned `fj` CLI for Forgejo operations. An already-authenticated
host installation is preferred; the configured-instance Vault wrapper remains
available as `bazel run //tools/fj -- ...`. An unavailable AL adapter does not
authorize publishing to a GitHub mirror.

## Authentication and command context

The repository's `tools/fj/README.md` links the installed command and pinned
executable. Use `fj --host https://forgejo.alwaldend.com whoami` to
check the native account. Native login and the wrapper's short-lived
Vault-derived API token are separate authentication paths; neither requires
changing the repository remote. The wrapper authenticates to its configured
instance: do not forward that instance's injected token to another host.

`fj --host <instance-url> auth login` starts browser OAuth when a client ID is
available for that instance. A custom instance can require a client-ID mapping;
the CLI reporting unsupported login does not establish that Forgejo lacks OAuth.
The [pinned authentication source](https://codeberg.org/forgejo-contrib/forgejo-cli/src/tag/v0.7.0/src/auth.rs)
reads the user configuration's `forgejo-cli/client_ids` on Linux (normally
`~/.config/forgejo-cli/client_ids`), then `/etc/fj/client_ids`:
each non-comment line is `host client_id`. Its public-client PKCE callback is
`http://127.0.0.1:26218/`, so the browser callback must reach the host running fj.
Configure a matching registered public client through the owning infrastructure
workflow before retrying OAuth; do not create or change server applications
merely to work around a failed login without authority.
If the user selects token login, `fj auth add-token` accepts stdin when its
token argument is omitted. Use a private credential handoff, never token-bearing
arguments or chat. Do not provision a password, token or native login merely
to inspect help.

Pass `-C "$PWD"` so repository discovery uses the intended worktree, including
through Bazel's runfiles. Select the remote explicitly at the `pr` command:

```sh
fj -C "$PWD" pr --remote origin search
fj -C "$PWD" pr --remote origin view 123
fj -C "$PWD" pr --remote origin status 123
fj -C "$PWD" pr --remote origin review 123 list --comments
```

Replace `origin` with the remote selected by the delivery skill. Inspect help
without Vault authentication through the pinned binary:

```sh
bazel run //third_party/org_codeberg_forgejo_contrib_forgejo_cli_bin -- pr --help
```

The repository pin in
`third_party/org_codeberg_forgejo_contrib_forgejo_cli_bin/binary_toolchain.json`
selects fj 0.7.0. Its help exposes PR search, create, view, status, metadata editing,
comments and review listing, including inline comments. Check subcommand help
before relying on other operations; help inspection is not live API acceptance.

## Preserve Git and PR ownership

The Forgejo fallback requires a completely clean worktree and index before
`$al-rebase`; it does not preserve dirty files or autostash them.
Preserve unrelated work by stopping until its owner moves or commits it. Stage
only task-owned paths or hunks, and maintain exactly one non-merge task-owned
feature commit. Keep its aggregate message and signing correct. A commit,
amend, rebase, or conflict resolution invalidates
affected earlier checks; rerun them against the resulting tree before the
lease-protected push. Complete the skill's local `review-agent` review before
each PR creation or update, repeating it when a rebase or fix changes the
reviewed diff. Follow `$al-rebase` through its final fetch,
explicit lease, push, and Git-level verification. Stop rather than rewrite
shared, stacked, human-owned, or ambiguous history.

Read the optional agent co-author with
`git config --global --no-includes --get-all al.deliveryCoAuthor`. When set,
include the single configured `Name <email>` value once as a `Co-authored-by`
trailer in the same block as the final LLM disclaimer. Preserve other
co-authors. Resolve malformed or multiple values before committing. Inspect
the committed trailers with `git show -s --format='%(trailers:only)' HEAD`.

Use `fj pr search` and `pr view` to inspect candidate pull requests before
creating one. Create the exact-head pull request only when its absence is
clear. For an agent-owned pull request, set the title from Git's `%s` and the
body from `%b`, replacing only its final commit-disclaimer line with the pull
request disclaimer. Use `pr edit` for those fields; preserve human-authored
metadata by stopping for direction. Treat `fj`'s human-oriented output as
observational context only: never scrape or parse it as structured proof, and
do not claim exact metadata or review postconditions that the command cannot
verify mechanically.

Use an ignored body file and explicit base/head for creation, for example:

```sh
fj -C "$PWD" pr --remote origin create \
  --base master --head feature-branch \
  --body-file "$PWD/out/task/pull-request.md" "Pull request title"
```

The example shows CLI syntax, not authorization to skip the checks, ownership
inspection or local review above. Preserve an existing PR's actual base.

The exact ownership markers are
`LLM-disclaimer: This commit was generated by an LLM.` for the commit and
`LLM-disclaimer: This pull request was generated by an LLM.` for the pull
request. Do not improvise their spelling or padding.

## Inspect review feedback within CLI capabilities

fj 0.7.0 can display and add top-level PR comments, list reviews and their inline
comments, and request reviewers. Request a human review only when explicitly
authorized; the skill's prohibition on initiating Codex reviews still applies.
Its help exposes no thread-reply or thread-resolution operation and no structured
inventory comparable to the GitHub adapter. Evaluate feedback through
`pr view <id> comments` and `pr review <id> list --comments`. After a valid fix,
amend, revalidate, synchronize, and
push through `$al-rebase`, then leave a reasoned top-level
`pr comment`. For invalid feedback, leave the concrete reason as a top-level
comment. Do not claim that either comment is an in-thread reply, and leave the
original thread or unresolved state intact. Report these review limitations
explicitly rather than silently skipping or falsely completing them.
