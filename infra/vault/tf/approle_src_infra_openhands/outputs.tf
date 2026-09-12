output "entity_id" {
  description = "AppRole entity id"
  value       = module.approle.entity_id
}

output "group_id" {
  description = "AppRole group id"
  value       = module.approle.group_id
}

output "ssh_policy" {
  description = "Policy allowing host key signing for the component domain"
  value       = module.ssh.policy
}
