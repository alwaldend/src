## Purpose

Define Host Bot's ComfyUI DNS and HTTPS forwarding contract alongside T3 Code.

## ADDED Requirements

### Requirement: ComfyUI host routing and DNS

Host Bot SHALL declare `comfyui.host-bot.simeonwarren.users.alwaldend.com`
with a site-local A record pointing to Host Bot and a public CNAME pointing to
ingress. Its dedicated `dc1` backend name SHALL point to Host Bot in both views.
Local Traefik SHALL match the service hostname and forward to
`http://127.0.0.1:8188`, using the existing Vault TLS resolver.

#### Scenario: Render ComfyUI configuration

- **WHEN** the checked-in DNS and Traefik definitions are evaluated
- **THEN** local and public DNS follow the T3 Code view pattern
- **AND** the ComfyUI HTTPS router forwards to loopback port 8188
- **AND** the existing T3 Code route is preserved
