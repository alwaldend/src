// Package main implements Vault certificate login with PKCS#11 hardware and
// atomic private token storage. AL owns Vault endpoint and CA injection; native
// p11tool and curl own hardware discovery and TLS signing. This command never
// prints authentication responses or tokens and changes no host configuration.
package main
