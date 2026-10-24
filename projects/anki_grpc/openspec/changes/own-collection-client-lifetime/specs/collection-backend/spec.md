## REMOVED Requirements

### Requirement: Borrow prepared backend lifetime

**Reason**: Generic clients own transports, not reusable collection phases.
**Migration**: Use generated native collection RPCs or Anki as Code's prepared
collection handle; do not wrap one Client lifecycle in another.

### Requirement: Direct databases and native package restore

**Reason**: Database paths, archive restore and temporary storage are application
collection concerns.
**Migration**: Use the application's collection scope for managed paths or own
native import/open/close calls explicitly.

### Requirement: Respect native consuming operations and failures

**Reason**: Managed collection readiness and file retention are not transport
responsibilities.
**Migration**: The application scope tracks calls needed for its own publication
and cleanup; raw RPC callers retain native lifetime and failure responsibility.
