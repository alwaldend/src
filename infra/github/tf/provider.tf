terraform {
  required_providers {
    github = {
      source  = "integrations/github"
      version = "6.10.2"
    }
  }
  backend "http" {
  }
}

provider "github" {
  owner = local.github_organization_name
}
