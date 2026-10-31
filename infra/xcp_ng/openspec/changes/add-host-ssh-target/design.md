## Decision

Use al_binary_run with the existing xcpng_host injector label and host sshpass
and OpenSSH clients. The host administrator account is root. Keep password
authentication in the environment and host keys in the controller's existing
known_hosts. Disable host SSH configuration loading to avoid the local
system SSH proxy configuration permission failure.

No host settings, VM allocation, or PCI assignments are changed.
