terraform {
  required_version = ">= 1.9"
}

variable "zone" {
  type        = string
  default     = "alwaldend.com"
  description = "DNS zone containing record names; a final dot is optional."

  validation {
    condition     = can(regex("^[a-zA-Z0-9-]+(\\.[a-zA-Z0-9-]+)+\\.?$", var.zone))
    error_message = "zone must be a DNS domain name."
  }
}

variable "document" {
  type        = any
  description = "Decoded canonical dnsconfig.json; each logical key contains scalar record-type members and dsp destinations."

  validation {
    condition     = try(toset(keys(var.document)) == toset(["records"]) && can(keys(var.document.records)), false)
    error_message = "document must contain exactly one records object."
  }

  validation {
    condition = try(alltrue([
      for key, record in var.document.records :
      can(regex("^[a-zA-Z0-9_.-]+$", key)) &&
      length(setsubtract(toset(keys(record)), toset(["A", "AAAA", "CNAME", "NS", "MX", "TXT", "dsp"]))) == 0 &&
      length(setintersection(toset(keys(record)), toset(["A", "AAAA", "CNAME", "NS", "MX", "TXT"]))) > 0 &&
      length(record.dsp) > 0 &&
      can(tolist(record.dsp)) &&
      alltrue([for view in record.dsp : contains(["global", "dc1", "all"], view)])
    ]), false)
    error_message = "Each logical record requires a stable key, supported type members, and a nonempty dsp list containing global, dc1, or all."
  }

  validation {
    condition = try(alltrue(flatten([
      for key, record in var.document.records : [
        for type, member in record : type == "dsp" ? true : (
          length(setsubtract(toset(keys(member)), toset(concat(
            ["name", "ttl"],
            contains(["A", "AAAA", "NS"], type) ? ["address"] : [],
            contains(["CNAME", "MX"], type) ? ["target"] : [],
            type == "MX" ? ["priority"] : [],
            type == "TXT" ? ["content"] : [],
          )))) == 0 &&
          try(tostring(member.name) == member.name && length(member.name) > 0, false) &&
          try(
            type == "TXT" ? member.content != null && tostring(member.content) == member.content :
            contains(["A", "AAAA", "NS"], type) ? tostring(member.address) == member.address && length(member.address) > 0 :
            tostring(member.target) == member.target && length(member.target) > 0,
            false,
          ) &&
          try(member.ttl == null ? false : member.ttl == floor(member.ttl) && member.ttl >= 60 && member.ttl <= 86400, !contains(keys(member), "ttl")) &&
          (type == "MX" ? try(member.priority == floor(member.priority) && member.priority >= 0 && member.priority <= 65535, false) : true)
        )
      ]
    ])), false)
    error_message = "Type members must be scalar objects with name and address/target/content, only supported fields, an integer TTL between 60 and 86400, and an MX priority between 0 and 65535."
  }

  validation {
    condition = try(alltrue(flatten([
      for key, record in var.document.records : [
        for type, member in record : type == "dsp" ? true : (
          member.name == "@" || (
            can(regex("^(\\*\\.)?[_a-zA-Z0-9-]+(\\.[_a-zA-Z0-9-]+)*\\.?$", member.name)) &&
            (!endswith(member.name, ".") || lower(trimsuffix(member.name, ".")) == lower(trimsuffix(var.zone, ".")) || endswith(lower(trimsuffix(member.name, ".")), ".${lower(trimsuffix(var.zone, "."))}"))
          )
        )
      ]
    ])), false)
    error_message = "Record names must be relative DNS names, @, or absolute names within zone."
  }

  validation {
    condition = try(alltrue(flatten([
      for key, record in var.document.records : [
        for type, member in record :
        type == "A" ? can(cidrnetmask("${member.address}/32")) :
        type == "AAAA" ? can(cidrhost("${member.address}/128", 0)) && strcontains(member.address, ":") :
        contains(["CNAME", "MX", "NS"], type) ? can(regex("^(@|[_a-zA-Z0-9-]+(\\.[_a-zA-Z0-9-]+)*\\.?)$", try(member.target, member.address))) : true
      ]
    ])), false)
    error_message = "A and AAAA values must match their address families; CNAME, MX, and NS targets must be domain names or @."
  }
}

locals {
  zone = lower(trimsuffix(var.zone, "."))
  # A and NS retain the effective legacy default. AAAA and CNAME carried
  # explicit TTL(600); the legacy A call placed that modifier outside A().
  default_ttls = { A = 300, AAAA = 600, CNAME = 600, NS = 300, MX = 300, TXT = 300 }
  members = try(flatten([
    for logical_key, record in var.document.records : [
      for type, member in record : [
        for view in sort(distinct(flatten([
          for destination in record.dsp : destination == "all" ? ["global", "dc1"] : [destination]
          ]))) : {
          key      = "${logical_key}/${type}/${view}"
          name     = lower(member.name)
          type     = type
          value    = type == "TXT" ? member.content : contains(["A", "AAAA", "NS"], type) ? member.address : member.target
          ttl      = try(member.ttl, local.default_ttls[type])
          priority = type == "MX" ? member.priority : null
          view     = view
          proxied  = false
        }
      ] if type != "dsp"
    ]
  ]), [])
  qualified_members = [for record in local.members : merge(record, {
    name = record.name == "@" ? local.zone : (
      trimsuffix(record.name, ".") == local.zone || endswith(trimsuffix(record.name, "."), ".${local.zone}") ?
      trimsuffix(record.name, ".") : "${record.name}.${local.zone}"
    )
    value = contains(["CNAME", "MX", "NS"], record.type) ? (
      record.value == "@" ? local.zone : (
        endswith(record.value, ".") || lower(record.value) == local.zone || endswith(lower(record.value), ".${local.zone}") ?
        lower(trimsuffix(record.value, ".")) : "${lower(record.value)}.${local.zone}"
      )
    ) : record.value
  })]
}

output "zone" {
  description = "Canonical DNS zone used to normalize the declarations."
  value       = local.zone
}

output "normalized_records" {
  description = "Canonical per-provider records indexed by logical key/type/view; DNS names and host targets have no final dot."
  value = { for record in local.qualified_members : record.key => merge(record, {
    relative_name = record.name == local.zone ? "@" : trimsuffix(record.name, ".${local.zone}")
  }) }
}
