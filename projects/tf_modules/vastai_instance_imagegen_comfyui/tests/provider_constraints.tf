# Offline tests use the repository-built prerelease provider. Consumers own
# their provider constraints; this override is packaged only in the tests.
terraform {
  required_providers {
    vastai = {
      source  = "alwaldend/vastai"
      version = "0.0.0-dev"
    }
  }
}
