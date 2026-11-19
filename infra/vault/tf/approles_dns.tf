locals {
  dns_approle_context = {
    secrets           = vault_mount.secrets.path
    backend           = vault_auth_backend.approle.path
    backend_accessor  = vault_auth_backend.approle.accessor
    member_entity_ids = [vault_identity_entity.simeonwarren.id]
  }
}

module "src_infra_mikrotik" {
  source  = "./approles/src_infra_mikrotik"
  context = local.dns_approle_context
}

module "src_infra_nas" {
  source  = "./approles/src_infra_nas"
  context = local.dns_approle_context
}

module "src_projects_activitywatch_ingester_android" {
  source  = "./approles/src_projects_activitywatch_ingester_android"
  context = local.dns_approle_context
}

module "src_projects_agents" {
  source  = "./approles/src_projects_agents"
  context = local.dns_approle_context
}

module "src_projects_al" {
  source  = "./approles/src_projects_al"
  context = local.dns_approle_context
}

module "src_projects_android_launcher" {
  source  = "./approles/src_projects_android_launcher"
  context = local.dns_approle_context
}

module "src_projects_ansible_collection" {
  source  = "./approles/src_projects_ansible_collection"
  context = local.dns_approle_context
}

module "src_projects_autoscroll" {
  source  = "./approles/src_projects_autoscroll"
  context = local.dns_approle_context
}

module "src_projects_bazel_agent" {
  source  = "./approles/src_projects_bazel_agent"
  context = local.dns_approle_context
}

module "src_projects_ci_platform" {
  source  = "./approles/src_projects_ci_platform"
  context = local.dns_approle_context
}

module "src_projects_dotfiles" {
  source  = "./approles/src_projects_dotfiles"
  context = local.dns_approle_context
}

module "src_projects_hugo_landing" {
  source  = "./approles/src_projects_hugo_landing"
  context = local.dns_approle_context
}

module "src_projects_infinitime" {
  source  = "./approles/src_projects_infinitime"
  context = local.dns_approle_context
}

module "src_projects_kustomization" {
  source  = "./approles/src_projects_kustomization"
  context = local.dns_approle_context
}

module "src_projects_leetcode_downloader" {
  source  = "./approles/src_projects_leetcode_downloader"
  context = local.dns_approle_context
}

module "src_projects_mcp_cordis" {
  source  = "./approles/src_projects_mcp_cordis"
  context = local.dns_approle_context
}

module "src_projects_nexus_security_plugin" {
  source  = "./approles/src_projects_nexus_security_plugin"
  context = local.dns_approle_context
}

module "src_projects_renders" {
  source  = "./approles/src_projects_renders"
  context = local.dns_approle_context
}

module "src_projects_rules_binary_toolchain" {
  source  = "./approles/src_projects_rules_binary_toolchain"
  context = local.dns_approle_context
}

module "src_projects_rules_dnscontrol" {
  source  = "./approles/src_projects_rules_dnscontrol"
  context = local.dns_approle_context
}

module "src_projects_rules_docs" {
  source  = "./approles/src_projects_rules_docs"
  context = local.dns_approle_context
}

module "src_projects_rules_docs_gazelle" {
  source  = "./approles/src_projects_rules_docs_gazelle"
  context = local.dns_approle_context
}

module "src_projects_rules_hugo" {
  source  = "./approles/src_projects_rules_hugo"
  context = local.dns_approle_context
}

module "src_projects_rules_iso" {
  source  = "./approles/src_projects_rules_iso"
  context = local.dns_approle_context
}

module "src_projects_rules_promptfoo" {
  source  = "./approles/src_projects_rules_promptfoo"
  context = local.dns_approle_context
}

module "src_projects_rules_promptfoo_gazelle" {
  source  = "./approles/src_projects_rules_promptfoo_gazelle"
  context = local.dns_approle_context
}

module "src_projects_rules_skill_gazelle" {
  source  = "./approles/src_projects_rules_skill_gazelle"
  context = local.dns_approle_context
}

module "src_projects_rules_skills" {
  source  = "./approles/src_projects_rules_skills"
  context = local.dns_approle_context
}

module "src_projects_rules_template" {
  source  = "./approles/src_projects_rules_template"
  context = local.dns_approle_context
}

module "src_projects_sri" {
  source  = "./approles/src_projects_sri"
  context = local.dns_approle_context
}

module "src_projects_tf_modules" {
  source  = "./approles/src_projects_tf_modules"
  context = local.dns_approle_context
}

module "src_projects_useless_qt_gui" {
  source  = "./approles/src_projects_useless_qt_gui"
  context = local.dns_approle_context
}

module "src_users_simeonwarren_host_bot" {
  source  = "./approles/src_users_simeonwarren_host_bot"
  context = local.dns_approle_context
}
