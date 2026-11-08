// Package main injects referenced Vault resources into commands through AL.
// SSH credentials use separate temporary private-key, public-key and certificate
// files, which are removed by the injector lifecycle. Secret values are supplied
// at runtime and are not persisted in repository source.
package main
