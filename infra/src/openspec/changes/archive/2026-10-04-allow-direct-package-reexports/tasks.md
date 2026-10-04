## 1. Correct re-exports

- [x] 1.1 Extend fixtures for direct package exports and forbidden imports.
- [x] 1.2 Update the skill and checker and simplify all affected exports.
- [x] 1.3 Regenerate BUILD declarations and validate consumers.

Validation: import fixtures, recommender HTTP E2E, Autoscroll consumer tests,
and offline skill configuration passed. Gazelle regenerated the recommender
packages; existing Autoscroll declarations were preserved for its legacy layout.
