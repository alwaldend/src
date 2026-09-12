## 1. Integrate the project DNS stage

- [x] 1.1 Package the owner declaration and shared global module; verify the Terraform target builds with complete runfiles.
- [x] 1.2 Wire the project AppRole, state backend, and Cloudflare injection to the stage label; inspect AL configuration for project ownership and absence of RouterOS dependencies.
- [x] 1.3 Preserve disabled record ownership and the optional zone input; inspect both source defaults and document that operational provider prerequisites still apply.

## 2. Validate integration

- [x] 2.1 Generate provider locks through the owning Terraform workflow and pass configured Terraform formatting, Bazel packaging, and OpenSpec validation.
