{
  "cloudflare": {
    "TYPE": "CLOUDFLAREAPI",
    "accountid": "{{ .Res.cloudflare.Data.cloudflare_account_id }}",
    "apitoken": "{{ .Res.cloudflare.Data.cloudflare_api_token }}"
  },
  "bind": {
    "TYPE": "BIND",
    "directory": "infra/dns/zones"
  },
  "mikrotik": {
    "TYPE": "MIKROTIK",
    "host": "https://router1.dc1.alwaldend.com",
    "username": "{{ .Res.mikrotik.Data.username }}",
    "password": "{{ .Res.mikrotik.Data.password }}",
  }
}
