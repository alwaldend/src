locals {
  signing_key_tool  = abspath("${path.module}/../../../tools/vault/cmd/ssh_key/ssh_key")
  signing_key_mount = "secrets"
  signing_key_path  = "alwaldend.com/vault1/approles/src_infra_dc1_forgejo1/commit_signing"
}
