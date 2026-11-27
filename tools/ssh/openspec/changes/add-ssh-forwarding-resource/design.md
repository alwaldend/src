## Decisions

- Borrow the client key and certificate from a declared AL resource dependency.
- Verify the SSH host certificate against the configured CA and host principal.
- Use an OS-assigned IPv4 loopback port for each resource. Verify the remote
  destination before returning the local address.
- Tie the listener and SSH connection to the resource provider lifecycle.

## Verification

Verify data forwarding with a local certified SSH server. Exercise a wrong
host CA, an unavailable destination, separate ports for simultaneous
invocations, and socket cleanup after exit.
