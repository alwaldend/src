---
title: Al proto
description: Protobuf contracts
---

The root configuration contains name-keyed `vault_conn`, `vault_auth`,
`plugin_configs`, `plugins` and `inject` maps. Entry values carry their settings;
the keys own identity. Injection groups keep their resource maps and selectors.
Their explicit `extends` reference determines override precedence without relying
on iteration order. Arguments, resource dependencies and result files retain
their ordered sequence semantics.

Plugin startup carries the selected executable key in `PluginStartRequest.name`
and selected calls in the plugin's named `calls` map. Startup responses advertise
resource types alongside an optional generic result. Provider fetches carry
typed resource identity, opaque arguments, configured input data and ordered
dependency results.
