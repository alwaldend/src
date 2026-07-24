{
  "global": {
    "TYPE": "CLOUDFLAREAPI",
    "accountid": "{{ .Res.cloudflare.Data.cloudflare_account_id }}",
    "apitoken": "{{ .Res.cloudflare.Data.cloudflare_api_token }}"
  },
  "global_bind": {
    "TYPE": "BIND",
    "directory": "infra/dns/zones"
  },
  "dc1": {
    "TYPE": "MIKROTIK",
    "host": "https://router1.dc1.alwaldend.com",
    "username": "{{ .Res.mikrotik.Data.username }}",
    "password": "{{ .Res.mikrotik.Data.password }}",
  },
  "dc1_bind": {
    "TYPE": "BIND",
    "directory": "infra/dns/zones"
  }
}
