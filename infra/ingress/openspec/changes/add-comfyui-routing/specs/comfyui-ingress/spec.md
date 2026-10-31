## Purpose

Define authenticated HTTPS ingress forwarding for the ComfyUI service on Host Bot.

## ADDED Requirements

### Requirement: Route ComfyUI through authenticated ingress

Ingress SHALL route `comfyui.host-bot.simeonwarren.users.alwaldend.com` to
`https://dc1.comfyui.host-bot.simeonwarren.users.alwaldend.com`, preserving
the request Host and using the service hostname for backend TLS SNI.
ComfyUI SHALL inherit ingress's default client-certificate authentication.

#### Scenario: Render ComfyUI ingress

- **WHEN** the ingress service map is rendered
- **THEN** ComfyUI has a host router and HTTPS backend with matching TLS SNI
- **AND** default TLS requires a verified client certificate
- **AND** the T3 Code ingress route remains unchanged
