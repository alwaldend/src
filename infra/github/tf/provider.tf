terraform {
  required_providers {
    github = {
      source  = "integrations/github"
      version = "6.6.0"
    }
  }
  backend "http" {
  }
}

provider "github" {
}
