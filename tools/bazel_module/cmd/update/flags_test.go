package main

import "flag"

// testGit identifies the Git executable provided by each Bazel test target.
var testGit = flag.String("git", "", "Git executable runfile")
