resource "consul_acl_policy" "nodes" {
  name  = "nodes"
  rules = <<-RULE
    agent_prefix "" {
      policy = "write"
    }
    node_prefix "" {
      policy = "write"
    }
    service_prefix "" {
      policy = "read"
    }
    session_prefix "" {
      policy = "read"
    }
  RULE
}
