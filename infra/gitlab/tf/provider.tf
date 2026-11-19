terraform {
  required_providers {
    gitlab = {
      source  = "gitlabhq/gitlab"
      version = "19.2.1"
    }
  }
  backend "http" {
  }
}

provider "gitlab" {
}
