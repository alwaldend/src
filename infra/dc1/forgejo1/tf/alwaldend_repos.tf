resource "forgejo_repository" "alwaldend_src" {
  owner          = forgejo_organization.alwaldend.name
  name           = "src"
  description    = "Source code"
  website        = "https://alwaldend.com/"
  default_branch = "master"
  clone_addr     = "https://github.com/alwaldend/src.git"
}
