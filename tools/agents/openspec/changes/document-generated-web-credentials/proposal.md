## Why

Infrastructure using a shared web login needs consistent secret ownership
and limited credential distribution.

## What Changes

Document Terraform-generated logins, scoped Vault access for administrators,
and bcrypt-only delivery to the reverse proxy.

## Impact

Updates the infrastructure skill and its offline evaluation fixtures.
