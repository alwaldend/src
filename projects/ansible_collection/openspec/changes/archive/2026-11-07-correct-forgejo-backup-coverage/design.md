## Context

See proposal.md for the accepted review defects. The existing controlled
activation and complete-recovery requirements remain authoritative.

## Decisions

Proceed with a conservative literal INI subset. The pinned Forgejo provider
uses go-ini with continuation disabled; its parser removes surrounding quotes.
Accept plain and simply single-, double- or backtick-quoted one-line values.
Reject ambiguous escapes, mixed quotes, inline comments and multiline values
instead of resolving a different persistence path. Normalize before Python's
standard INI parser; do not add another parser dependency or attempt a complete
independent implementation of go-ini.

Resolve the three documented legacy storage settings through their modern
replacement when nonempty. Otherwise retain the legacy value, then the normal
storage default. Require explicit local storage paths to be absolute, matching
the documented storage contract; advanced named storage profiles require a
separately coordinated procedure. Compare effective subsystem identities as
well as covered path sets to reject implicit storage relocation.

Require all three startup files to be covered by overridden backup roots
before stopping. Keep the user's explicit archive scope rather than silently
expanding it.

## Failure cases and acceptance before corrections

- Existing quoted repositories outside default roots must reject an upgrade
  with omitted coverage before any stop or archive.
- Independently omit legacy LFS, avatar and repository-avatar storage roots;
  each omission must reject before downtime.
- Independently omit executable, configuration and unit roots; each omission
  must reject before downtime.
- Changing legacy storage to another covered root must reject relocation;
  modern settings must override unused legacy values.
- A complete archive must restore the external repository and storage files,
  prior executable and configuration after those trees are actually removed.
- Explicit local subsystems retain their application-data defaults. Global
  `[storage].PATH` requires a separately coordinated migration and backup and
  must reject before downtime, including an explicitly empty setting.
- Unsupported quoting or multiline syntax must reject before downtime.

The existing real-binary integration scenario will exercise these cases before
the implementation is corrected. Preserve the systemctl fixture boundary and
sanitized result artifacts; no unit tests or live services are introduced.

## Risks and scope

Unusual INI syntax or storage profiles now fail safely before activation.
Document the supported subset. Do not deploy, automatically restore, change
production credentials or broaden external-storage support.

Sources: pinned
[INI parser](https://raw.githubusercontent.com/go-ini/ini/v1.67.0/parser.go),
[Forgejo provider](https://codeberg.org/forgejo/forgejo/src/tag/v17.0.1/modules/setting/config_provider.go),
and [storage compatibility](https://forgejo.org/docs/v17.0/admin/setup/storage/#legacy-settings).

The attempted global-root fixture exposed the pinned release's initial
migration sanity guard, which rejects `[storage].PATH` for fresh/older databases
while allowing it for newer migration histories. The role conservatively
requires per-subsystem absolute paths across both cases rather than inferring
historical storage migrations. See the pinned
[migration guard](https://codeberg.org/forgejo/forgejo/src/tag/v17.0.1/services/forgejo/sanity_v1TOv5_0_1Included.go).
