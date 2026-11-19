## Design

Declare each destination under `forgejo.push_mirrors` in its repository
catalog file. The catalog projects the destination entries to one flat
map for Terraform with stable repository-and-forge keys.
