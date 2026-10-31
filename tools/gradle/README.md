# Gradle wrapper

`//tools/gradle:gradle-wrapper` runs the pinned Gradle distribution from
`third_party/net_gradle_gradle` with the Java runtime selected by the existing
`rules_java` toolchain. It does not download anything at run time.
The launcher derives runfile paths from those declared targets and includes
the existing `rules_shell` runtime; it does not hardcode external repository
names.

```sh
al bazel run //tools/gradle:gradle-wrapper -- --version
```

Gradle caches remain under ignored `out/gradle/`. Set `GRADLE_CACHE_ROOT` to
redirect them for a task.
