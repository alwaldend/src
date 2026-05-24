{
  "cloudflare": {
    "TYPE": "CLOUDFLAREAPI",
    "accountid": "{{ .Secrets.cloudflare.cloudflare_account_id }}",
    "apitoken": "{{ .Secrets.cloudflare.cloudflare_api_token }}"
  },
  "bind": {
    "TYPE": "BIND",
    "directory": "infra/dns/zones"
  },
  "mikrotik": {
    "TYPE": "MIKROTIK",
    "host": "https://router1.dc1.alwaldend.com",
    "username": "{{ .Secrets.mikrotik.username }}",
    "password": "{{ .Secrets.mikrotik.password }}",
  }
}
