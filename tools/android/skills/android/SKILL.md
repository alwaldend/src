---
name: android
description: >-
  Build, package, and publish Android apps from this Bazel monorepo. Covers
  rules_android, the Gradle parallel build required for F-Droid, and the
  dependency lock conversion between Bazel and Gradle.
---

# Android apps in the monorepo

## Build system

The primary build system is Bazel (`rules_android` + `rules_kotlin` via
Bzlmod). Every Android app must also maintain a parallel Gradle build
suitable for F-Droid submission. The Gradle build is a derived artifact
generated from the Bazel dependency lock; it is not independently authored.

### Bazel build

Each Android app lives at `projects/<app>/` with a root `BUILD.bazel`
containing an `android_binary` target. Inspect the app's
`include.MODULE.bazel` for its Maven dependencies, which are managed by
`rules_jvm_external` with a `maven_lock.json` lock file.

```sh
bazel_agent bazel build //projects/<app>/main/java:<binary_target>
```

### Gradle parallel build (required for F-Droid)

F-Droid's buildserver has no Bazel support. Every Android app in this repo
must maintain a parallel Gradle build that consumes the same pinned
dependencies so an F-Droid reviewer sees a conventional Gradle project.

**Generated Gradle inputs.** A Go converter
(`tools/maven_install_gradle_converter/cmd/gradle_lock_gen`) reads the `maven_lock.json` lock file
and produces `gradle/libs.versions.toml` and
`gradle/verification-metadata.xml`. The `al_gradle_lock` Starlark macro
(`tools/maven_install_gradle_converter/defs.bzl`) wires these through
`write_source_file`, so a Bazel test (`*_write_test`) fails when the checked-in
Gradle files go stale relative to the Bazel lock.

To regenerate after a dependency change:

```sh
bazel_agent bazel run //projects/<app>/gradle:update_toml_write
bazel_agent bazel run //projects/<app>/gradle:update_verification_write
```

To verify freshness in CI:

```sh
bazel_agent bazel test //projects/<app>/gradle:update_toml_write_test
bazel_agent bazel test //projects/android_launcher/gradle:update_verification_write_test
```

**Manual Gradle files.** `build.gradle.kts` and `settings.gradle.kts` are
hand-authored (not generated). They reference the version catalog from
`gradle/libs.versions.toml` and apply the Kotlin/Compose plugins. Keep the
source-set layout aligned with the Bazel source layout (`main/java/`,
`main/res/`, `AndroidManifest.xml`).

## F-Droid submission

- F-Droid does not support Bazel. Use the Gradle parallel build for
  fdroiddata metadata (`gradle: yes` in the `Builds:` entry).
- The app must have real `versionCode`/`versionName` values (not placeholders)
  and a Fastlane metadata directory at
  `fastlane/metadata/android/<package>/en-US/` with `title.txt`,
  `short_description.txt`, `full_description.txt`, and images.
- The Fastlane metadata is required before F-Droid submission. Local build
  validation can proceed without it, but publication cannot.
- F-Droid signs every APK itself; you provide the unsigned APK built from a
  pinned commit. The Gradle build must produce a working, reproducible APK.

### Local F-Droid build check

Use `tools/fdroid/fdroid-wrapper` for a local `fdroidserver` environment. Keep
the fdroiddata checkout and generated APK under ignored `out/<task>/`; never
place scratch fdroidserver state in the app source tree.

1. Create `out/<task>/fdroiddata/{metadata,build,tmp,unsigned,repo}` and
   `config.yml` pointing `sdk_path` at the local Android SDK. Resolve
   `jdk_path` from the rules_java remote JDK consumed by
   `//tools/gradle:gradle-wrapper` so the Bazel-pinned JDK remains the only
   build input.
2. Create a metadata YAML whose `Repo` points to a Git repository containing
   the app and whose `Builds.commit` is a real commit in that repository.
3. Initialize and commit the fdroiddata metadata in Git. fdroidserver derives
   reproducibility metadata from Git and fails if the metadata is untracked.
4. Run the build from `out/<task>/fdroiddata`:

```sh
../../../tools/fdroid/fdroid-wrapper readmeta
../../../tools/fdroid/fdroid-wrapper build <package>:<versionCode> --on-server --no-tarball
```

Do not debug local fdroidserver failures by repeatedly inspecting vendored
fdroidserver internals. First validate the four local prerequisites: metadata
parses, `Repo` is a real Git repository, `commit` exists, and fdroiddata is a
Git repository containing the metadata commit.

fdroidserver scans `<subdir>/build/outputs/apk/release/` for the unsigned APK.
When the Gradle application module is a subdirectory, add a `Copy` task that
publishes `*.apk` from `app/build/outputs/apk/release/` to the equivalent
directory under the metadata `subdir`, and finalize `packageRelease` with it.
Release `lintVital` can resolve AGP toolchain dependencies that are absent
from the shared lock; set `lint { checkReleaseBuilds = false }` for the
F-Droid release path unless those toolchain artifacts are intentionally added
to `verification-metadata.xml`.

### Manual instrumentation smoke test

Use `//projects/<app>/test:<app>_smoke_test` as a manual instrumentation test
when rules_android Bzlmod cannot provide `android_device`. It installs the
Bazel-built APK on an already running device or emulator, starts the launcher
activity, waits, and verifies that the process is still alive. Run it
explicitly:

```sh
bazel_agent bazel run //projects/<app>/test:<app>_smoke_test
```

Do not run it as a CI test; it requires local Android platform tools and a
connected device. Add one per Android app before F-Droid submission because
`build_test` does not prove that the installed APK starts successfully.

### Debug vs release APK signing

`assembleRelease` produces an unsigned APK that Android refuses to install
outside the F-Droid build server (which signs it server-side). Use
`assembleDebug` for local device and emulator testing; Gradle signs the debug
APK with the debug keystore so it installs directly.

```sh
bazel_agent bazel run //tools/gradle:gradle-wrapper -- assembleDebug
```

Run this from the app project directory with `ANDROID_HOME` set to the local
Android SDK. The debug APK lands at
`app/build/outputs/apk/debug/app-debug.apk`.

### Troubleshooting "App not installed"

The most common cause is a signature mismatch with an existing install.
Uninstall the old version first:

```sh
adb uninstall <package>
```

Other causes to check:

- **Unsigned release APK** — see the signing section above; use the debug
  variant for local installs.
- **ABI mismatch** — the APK may lack a native library for the device's CPU.
  Check with `aapt dump badging app.apk | rg native-code`.
- **Insufficient storage** — check `adb shell df /data`.
- **Corrupt APK** — verify it parses: `aapt dump badging app.apk`.
- **Unknown sources** — enable "Install unknown apps" for your file manager or
  adb source in Settings → Apps → Special access.

Watch `PackageManager` errors while attempting the install:

```sh
adb logcat | rg 'PackageManager|INSTALL_FAILED'
```

## Standard repo practices

- Put all Android source under `projects/<app>/main/java/` with resources in
  `projects/<app>/main/res/` and the manifest at `projects/<app>/`.
- Keep the `android_binary` target in `main/java/BUILD.bazel` with a
  `build_test` for CI verification.
- Declare Maven dependencies in the app's `include.MODULE.bazel` via
  `rules_jvm_external`, lock them in `maven_lock.json`, and regenerate the
  Gradle outputs after any change.
- Add any Gradle-only build inputs under `projects/<app>/app/` and declare
  them in BUILD files so they participate in repository validation. Do not
  leave compatibility symlinks between the conventional `app/src/main/proto`
  path and the Bazel layout; the parallel Gradle build must point its
  source set at the canonical Bazel path.
- Every Android app must provide a Fastlane metadata directory at
  `fastlane/metadata/android/<package>/en-US/` with `title.txt`,
  `short_description.txt`, `full_description.txt`, icon, and feature graphic
  assets so F-Droid metadata can be generated without divergence.
- The `al_gradle_lock` macro and the `gradle_lock_gen` Go binary are
  project-owned; do not create a shared Android tooling project for this.
- Android packages require a runtime smoke test before submission. Run the
  Bazel unit tests for repository behavior and install the built APK on an
  Android device or emulator, launch its launcher activity, and confirm the
  app starts without crashing. `build_test` proves the target builds; it does
  not prove the APK works. Gradle-generated `verification-metadata.xml` is
  excluded from Prettier by filename glob because its exact formatting is
  owned by Gradle.
