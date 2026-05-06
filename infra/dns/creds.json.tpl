{
  "cloudflare": {
    "TYPE": "CLOUDFLAREAPI",
    "accountid": "{{ .Secret.cloudflare_account_id }}",
    "apitoken": "{{ .Secret.cloudflare_api_token }}"
  },
  "bind": {
    "TYPE": "BIND",
    "directory": "infra/dns/zones"
  }
}
