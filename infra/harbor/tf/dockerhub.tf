resource "harbor_registry" "dockerhub" {
  provider_name = "docker-hub"
  name          = "dockerhub"
  description   = "Dockerhub"
  endpoint_url  = "https://hub.docker.com"
}

resource "harbor_project" "dockerhub" {
  name                           = "dockerhub"
  registry_id                    = harbor_registry.dockerhub.registry_id
  proxy_cache_local_on_not_found = true
}

resource "harbor_project_member_group" "dockerhub_src_infra_harbor_users" {
  project_id = harbor_project.dockerhub.id
  group_name = "src_infra_harbor_users"
  group_id   = 3
  role       = "guest"
  type       = "oidc"
}
