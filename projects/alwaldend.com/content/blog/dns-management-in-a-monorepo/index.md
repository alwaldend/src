---
title: DNS management in a monorepo
description: Evolution of DNS management in my monorepo
date: 2026-09-13
tags:
  - dns
  - dnscontrol
  - terraform
  - repo
---

Buying a domain is one of the best purchases you can make because they are very cheap and owning one offers you more independence - you can have your email not tied to a specific provider, you can have a site, etc.

After I bought a domain, I needed a way to manage DNS of that domain. Since using the registrar's UI to do that is not the best idea, I searched for an IaC tool that can do that for me - I settled on [DNSControl](https://dnscontrol.org/).

## Single-file config

Initially, I just needed to set up email and GitHub Pages, so I used the most basic setup imaginable.

```
infra/dns/
├── dnsconfig.js # DNS config
├── creds.json.tpl # Credential config
└── zones/
    └── alwaldend.com.zone # BIND file
```

I just wrote records in DNSControl's subset of JavaScript and templated secrets from the environment using `envsubst`. It worked alright, and DNSControl rendered a BIND file so I could see the computed DNS records as text.

## Single-file config with secret injection

At some point, I migrated my secrets to HashiCorp Vault because storing secrets in environment variables is a big no-no. That didn't change the layout. The credential config just became a Go template that used secrets from Vault.

## Per-project config

Once I started setting up infrastructure, I hit an annoying limitation - extracting records from `.js` is extremely inconvenient, which means that you have to duplicate IPs for VMs instead of just extracting them from the config. And managing all DNS records in a single `.js` file is not very convenient either. To solve it, I split the main config into several static per-project configs and then just parsed them in the main `.js` file.

The layout became like this:

```
infra/dns/ # The main DNS project
infra/project/ # Another infra project
├── dnsconfig.json # Static DNS config
```

And the config looked like this:

```json
{
  "domains": {
    "default": {
      "records": {
        "host1_a": {
          "A": {
            "name": "host1.consul1.dc1",
            "address": "192.168.10.20"
          }
        }
      }
    }
  }
}
```

This was a bad config design as will become clear later.

## Split-horizon DNS

Later, I decided to set up [Split-horizon DNS](https://en.wikipedia.org/wiki/Split-horizon_DNS) so I can have different records for internal and external DNS servers. Thankfully, DNSControl supports this, so there were no problems from this side, but configs had to be redone because they grouped records by domains which did not translate well to DNSControl's DNS providers.

I simplified the config and added `dsp` fields to select which DNS provider the record belongs to.

```json
{
  "records": {
    "host1": {
      "A": {
        "name": "host1.threexui",
        "address": "192.168.10.80"
      },
      "dsp": ["dc1"]
    }
  }
}
```

## Terraform

That setup worked fine for a while, but then I stumbled upon an obvious flaw during a clanker session - there is no way to work in parallel (without mucking around with `IGNORE` and `NO_PURGE`). DNSControl controls the whole zone, so if you have different sessions applying different configs, then those changes override each other.

That seems like an easy problem to fix - just split the domain into different zones, and it was easy for my local DNS (a MikroTik router) because DNSControl's provider has zone hints, but there is no way to split the domain on the external DNS (Cloudflare) without buying an Enterprise plan or without delegating to yet another DNS provider.

I thought about my options and decided to just switch to Terraform because I already have a standard way of handling Terraform in the repo. This approach fixed the problem, because now projects own only their own DNS records, so you can work on them in parallel.

To implement DNS control in Terraform, I wrote a Terraform module which parses the config and updates DNS, then used it in each project. I also added a linter that checks that projects do not overlap DNS records. The final layout looks like this:

```
projects/
├── project name/
│   ├── dnsconfig.json # DNS config, unchanged
│   └── tf/ # A Terraform directory
└── tf_modules/dns_records/ # Reusable Terraform module

infra/dns/
├── dnsconfig.json # Generic DNS like email
├── tf/ # A Terraform directory
└── cmd/lint/ # A linter that checks configs for collisions
```

## Future

There are no restrictions on DNS management - each project uses the same API keys and can modify the whole domain. This is not ideal, of course, but it's not that big of a deal because I trust all developers involved (Me and clankers) to not fuck it up. I might revisit this problem in the future, but for now this setup is good enough.
