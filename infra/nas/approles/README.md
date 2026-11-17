# AppRole storage projects

Each `<approle>/` directory is a standalone Terraform root with its own state,
plan/apply targets and fixed namespace under `default/encrypted/approles`.
It owns the namespace, ACLs, shares and child datasets through the shared
[dataset module](../../../projects/tf_modules/truenas_dataset/README.md).
No parent input or remote-state dependency is needed.

The [shared NAS root](../truenas/README.md) owns accounts, the encrypted parent
and appliance-wide settings. It does not load AppRole projects. Group membership
creates an account; a namespace is created when its storage project is applied.
Add another project only when declaring that AppRole's storage.

Projects use the NAS operator identity, the shared connection injector and a
separate backend under that identity's `tf_backend/approles/<approle>` path.
This separates state and lifecycle; it does not restrict TrueNAS administrator
API permissions to a dataset path. The consumer AppRole receives only its own
NAS account credential, not the provisioning key.

Apply the shared NAS prerequisites first, then the selected AppRole storage
project. The [development project](src_infra_dev/README.md) owns its namespace
and runner cache. VM provisioning consumes the resulting Xen Orchestra storage
repository through its own development identity.
