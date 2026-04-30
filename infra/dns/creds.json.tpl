{
  "cloudflare": {
    "TYPE": "CLOUDFLAREAPI",
    "accountid": "{{ env "CLOUDFLARE_ACCOUNT_ID" }}",
    "apitoken": "{{ env "CLOUDFLARE_API_TOKEN" }}"
  },
  "bind": {
    "TYPE": "BIND",
    "directory": "infra/dns/zones"
  }
}
