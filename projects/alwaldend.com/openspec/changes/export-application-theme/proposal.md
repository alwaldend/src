## Why

The X profile cleaner's localhost interface needs the site's maintained visual
theme without embedding the entire Hugo site or copying its palette.

## What Changes

- Export a Sass application entry point using the shared palette, typography,
  and pinned Bootstrap components.
- Keep shared palette selectors in one partial used by both the site and apps.
- Generalize the existing Sass compilation rule for both theme consumers.

## Capabilities

### New Capabilities

- `application-theme`: Offline application styling derived from the site sources.

### Modified Capabilities

None.

## Impact

Site assets and their Sass build rule, consumed by projects/x_profile_cleaner.
No dependency upgrades, site publication, or infrastructure changes.
