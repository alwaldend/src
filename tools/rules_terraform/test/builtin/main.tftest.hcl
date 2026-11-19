run "builtin_lifecycle" {
  command = apply

  assert {
    condition     = output.value == "packaged Terraform"
    error_message = "The declared Terraform configuration was not executed."
  }
}
