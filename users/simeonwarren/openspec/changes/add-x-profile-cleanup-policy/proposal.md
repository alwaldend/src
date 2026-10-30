## Why

Keep the requested X cleanup policy in the user directory so it can be loaded
by the CLI or browser without rebuilding the rules each time.

## What Changes

- Add a Lua policy selecting posts and reposts older than 30 days.
- Preserve the six supplied threads with descriptive comments.
- Document how to open the policy in the local browser interface.
- Add a separate single-post policy for the explicitly authorized live lifecycle.

## Capabilities

### New Capabilities

- `x-profile-cleanup-policy`: User-owned age rules and thread preservation.

### Modified Capabilities

None.

## Impact

User configuration and documentation change without new dependencies or cleaner
runtime changes. The separately requested live browser demonstration authorizes
only deletion of the single target named in its policy; the age-based policy
must not be applied as part of that demonstration.
