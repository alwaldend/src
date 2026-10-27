## 1. Native platform resolution

- [x] 1.1 Configure the complementary uv environments and regenerate the lock
      through `uv_lock.update`; compare existing versions and artifact hashes.
- [x] 1.2 Inspect unpatched native import metadata for Linux, macOS, and Windows;
      retain PyObjC only on macOS and record repeatable artifacts.
- [x] 1.3 Document the configuration, pass Autoscroll and full-site tests, and
      complete the shared dependency repair's publication gates.
