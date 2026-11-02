## Acceptance

The canonical project-layout skill now requires landing pages for new
first-party projects under projects/, including repository-only helpers.
It routes site mechanics to repo-hugo and preserves the distinction between
generated output in an existing owner and creating a new project. Its existing
generated-SDK case now explicitly uses an existing project; a new case covers
landing onboarding for a new repository-only helper.

The skill library passed its validation aspect. The offline Promptfoo config,
exact-state skill discovery, tools/agents and bridge OpenSpec checks, DNS suite,
and shared-site build test passed: nine tests. Semantic lint passed. These
checks validate packaging and evaluation configuration, not agent behavior;
no credentialed or live model evaluation ran.

The bridge's existing change records the landing page, registry membership,
site build, and representative HTML verification. No site was deployed.
