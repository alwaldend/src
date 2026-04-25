{
  "cloudflare": {
    "TYPE": "CLOUDFLAREAPI",
    "accountid": "{{ .Secrets.dns_token.cloudflare_account_id }}",
    "apitoken": "{{ .Secrets.dns_token.cloudflare_api_token }}"
  },
  "bind": {
    "TYPE": "BIND",
    "directory": "infra/dns/zones"
  }
}
