## Decisions

This shared agent documentation contract spans the canonical Go and Python skills, whose instructions remain in their existing owners. Python definitions use immediate body docstrings, while variables, aliases, and attributes use string literal docstrings immediately after their declarations or assignments, including nested contexts and methods outside `__init__`. Incomplete extraction support is a tool limitation, not a documentation exemption or a reason to move assignments or their documentation. Attribute docstrings do not become the assigned value runtime `__doc__`. Document project architecture in README files, and generated APIs through handwritten package docs. Offline eval validation checks packaging and configuration, not model behavior.

## Acceptance

Verify package documentation and target inclusion, formatting, semantic lint, and relevant builds. Publish independently against master.
