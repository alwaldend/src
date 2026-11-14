// Package ansible verifies profile selection and packaged workflow installation
// through the actual deployment playbook with disposable local inventories.
// Workflow-only runs avoid model downloads. The full profile case uses local
// substitutes and verifies playbook JSON translation. No case accesses
// Vault or contacts a rental.
package ansible
