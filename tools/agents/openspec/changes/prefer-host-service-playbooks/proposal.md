## Why

A full deployment should apply the existing host baseline before configuring
component services.

## What Changes

Prefer the shared host role, Traefik reverse proxy, and service roles in
default web-service deployments. Document the host inputs and when a
deployment can omit the proxy or run service roles alone.
Document the certificate name and enrollment pattern for guest proxies behind
shared ingress and explain how to configure direct LAN access.

## Impact

Updates the canonical infrastructure skill and its evaluation fixture.
