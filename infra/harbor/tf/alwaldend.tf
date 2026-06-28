resource "harbor_project" "alwaldend" {
  name                   = "alwaldend"
  vulnerability_scanning = true
  auto_sbom_generation   = true
}

resource "harbor_project_member_group" "alwaldend_src_infra_harbor_users" {
  project_id = harbor_project.alwaldend.id
  group_name = "src_infra_harbor_users"
  group_id   = 3
  role       = "guest"
  type       = "oidc"
}
