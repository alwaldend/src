---
title: Repository AL configuration
description: Shared configuration for repository command wrappers
---

`//tools/al:config` packages the repository's shared AL configuration. The
root `//:al` label remains a compatibility alias. The Lua source files stay
at the repository root so existing CLI configuration discovery and
`require("al_lib")` imports retain their paths.
