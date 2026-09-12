resource "vault_identity_group" "dns_approles" {
  name     = "dns_approles"
  type     = "internal"
  policies = []
  member_entity_ids = [
    module.src_infra_mikrotik.entity_id,
    module.src_infra_nas.entity_id,
    module.src_projects_activitywatch_ingester_android.entity_id,
    module.src_projects_agents.entity_id,
    module.src_projects_al.entity_id,
    module.src_projects_android_launcher.entity_id,
    module.src_projects_ansible_collection.entity_id,
    module.src_projects_autoscroll.entity_id,
    module.src_projects_bazel_agent.entity_id,
    module.src_projects_ci_platform.entity_id,
    module.src_projects_dotfiles.entity_id,
    module.src_projects_hugo_landing.entity_id,
    module.src_projects_infinitime.entity_id,
    module.src_projects_kustomization.entity_id,
    module.src_projects_leetcode_downloader.entity_id,
    module.src_projects_mcp_cordis.entity_id,
    module.src_projects_nexus_security_plugin.entity_id,
    module.src_projects_renders.entity_id,
    module.src_projects_rules_binary_toolchain.entity_id,
    module.src_projects_rules_dnscontrol.entity_id,
    module.src_projects_rules_docs.entity_id,
    module.src_projects_rules_docs_gazelle.entity_id,
    module.src_projects_rules_hugo.entity_id,
    module.src_projects_rules_iso.entity_id,
    module.src_projects_rules_promptfoo.entity_id,
    module.src_projects_rules_promptfoo_gazelle.entity_id,
    module.src_projects_rules_skill_gazelle.entity_id,
    module.src_projects_rules_skills.entity_id,
    module.src_projects_rules_template.entity_id,
    module.src_projects_sri.entity_id,
    module.src_projects_tf_modules.entity_id,
    module.src_projects_useless_qt_gui.entity_id,
    module.src_users_simeonwarren_host_bot.entity_id,
  ]
  metadata = {
    comment = "DNS-only AppRoles without general infrastructure privileges"
  }
}
