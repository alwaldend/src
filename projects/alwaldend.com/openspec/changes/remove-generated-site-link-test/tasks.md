## Remove the test

- [x] Delete `projects/alwaldend.com/test/site/site_test.go` and its BUILD
      target, including the now-empty test directory.
- [x] Drop the test from the project README validation commands.
- [x] Withdraw the baseline requirement that named the generated-site test.

## Validate

- [x] `//projects/alwaldend.com:site_test` still passes.
- [x] `//infra/src/openspec/validation:projects_alwaldend_com_validate_test`
      passes.
- [x] No surviving reference to `projects/alwaldend.com/test/site` remains.
