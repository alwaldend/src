## Decisions

Develop and test first-party rulesets as normal packages in the root module.
The root wildcard includes their targets. Dependency and toolchain
registrations stay with their owning includes.

DNSControl, systemd, and ISO are root packages and keep their existing archive
targets. For the other rulesets, the archive generator derives module names
from the project release manifests, maps package labels to their standalone
namespaces, and packages the license.

Root toolchain registrations live in the include files of their owning
components. Hugo checks in the root workspace use the pinned site toolchain.
The standalone Hugo module carries its development fixture.
