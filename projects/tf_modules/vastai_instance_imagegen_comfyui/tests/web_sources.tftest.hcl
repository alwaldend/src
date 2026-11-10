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
  image                   = "vastai/comfy:test"
  custom_node_sources     = { prompt_sequence = "NODE_CLASS_MAPPINGS = {}" }
  custom_node_web_scripts = { prompt_sequence = "// first-party browser script" }
}

run "stage_first_party_browser_script" {
  command = plan
  assert {
    condition = local.first_party_files == [
      { path = "/workspace/ComfyUI/custom_nodes/prompt_sequence/__init__.py", content = "NODE_CLASS_MAPPINGS = {}", permissions = "0600" },
      { path = "/workspace/ComfyUI/custom_nodes/prompt_sequence/web/index.js", content = "// first-party browser script", permissions = "0600" },
    ]
    error_message = "First-party Python and browser files must share the approved node directory."
  }
  assert {
    condition     = local.enabled_custom_nodes == tolist(["prompt_sequence"])
    error_message = "Browser files must not broaden the custom-node whitelist."
  }
}

run "reject_unapproved_script_directory" {
  command = plan
  variables {
    custom_node_web_scripts = { unknown = "// unapproved" }
  }
  expect_failures = [var.custom_node_web_scripts]
}
