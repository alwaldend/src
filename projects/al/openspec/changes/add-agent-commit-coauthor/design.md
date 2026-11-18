## Context

AL prepares an exact commit with `git commit-tree`. Forgejo delivery uses the
documented native Git procedure.

## Decisions

Read `al.deliveryCoAuthor` from the account's global Git configuration during
preparation. Limit the query to this key in the global files and preserve
transport isolation. Forgejo delivery reads the same key when composing the
message.

Keep configured attribution in the final Git trailer block so Git and forges
can parse it. Preserve the final LLM disclaimer and the commit author identity.
