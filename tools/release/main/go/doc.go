// Package main builds release bundles and publishes artifacts through OCI or
// SSH. It owns content transfer, archive validation, and static-site selection;
// callers own host provisioning, credentials, and authorization to publish.
package main
