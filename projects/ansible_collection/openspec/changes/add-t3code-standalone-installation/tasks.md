## 1. Role

- [x] 1.1 Add the opt-in installer and direct service invocation; verify
      archive layout, checksum, packaged role, and rendered unit.
- [x] 1.2 Protect environment writes and use change-triggered activation;
      verify playbook syntax and unchanged-run behavior in source.
- [x] 1.3 Remove the redundant daemon-reload task while ensuring unit changes
      reload systemd before the initial start and before a notified restart.
- [x] 1.4 Let archive comparison repair incomplete standalone installations;
      remove binary-existence guards and verify retry and unchanged behavior.
