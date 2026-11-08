## Decisions

Use AL's owning built-in-resource change and public res_<type> API. Shared config
registers the Vault provider once and propagates its runfiles. Built-in-only
fixtures carry no provider executable. Preserve Vault paths, labels and command
arguments; normalize only the deliberate generic resource wire representation
when comparing configuration output. Keep runtime protocol facts with AL.
