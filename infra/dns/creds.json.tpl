{
  "cloudflare": {
    "TYPE": "CLOUDFLAREAPI",
    "accountid": "{{ .Secret.data.data.cloudflare_account_id }}",
    "apitoken": "{{ .Secret.data.data.cloudflare_api_token }}"
  },
  "bind": {
    "TYPE": "BIND",
    "directory": "infra/dns/zones"
  }
}
