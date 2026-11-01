## Context

The existing service map drives HTTPS routers and transports with explicit
backend SNI. The default TLS option requires verified client certificates.

## Goals / Non-Goals

Add ComfyUI using this map. Preserve the shared templates and authentication.
Live ingress deployment is outside this source change.

## Decisions

Use `comfyui-host-bot` beside `t3code-host-bot`, targeting the dedicated `dc1`
name declared by Host Bot. Using the public service name as the target could
route back to ingress, so it is unsuitable.

## Risks / Trade-offs

- The backend DNS and local proxy must be deployed before the route is usable;
  keep their change linked in the proposal and validate the full rendered chain.
