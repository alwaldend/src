mock_provider "vastai" {}
override_module {
  target = module.imagegen
  outputs = {
    instance_id         = null
    ssh_host            = null
    ssh_port            = null
    ssh_forward_command = null
    local_url           = null
    https_url           = null
  }
}
variables {
  image               = "vastai/comfy:test"
  workspace_directory = "/srv/imagegen"
}
run "alternate_workspace" {
  command = plan
}
run "reject_workspace_escape" {
  command = plan
  variables { workspace_directory = "/srv/../etc" }
  expect_failures = [var.workspace_directory]
}
run "reject_workspace_shell_syntax" {
  command = plan
  variables { workspace_directory = "/srv/imagegen;id" }
  expect_failures = [var.workspace_directory]
}
