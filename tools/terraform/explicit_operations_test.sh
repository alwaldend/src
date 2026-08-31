#!/usr/bin/env bash

set -o errexit -o nounset -o pipefail

definitions="$1"
if grep -Eq '^ *"": *\[' "${definitions}"; then
    echo "unnamed Terraform operation must not exist" >&2
    exit 1
fi
if ! grep -Eq '^ *"apply": *\["apply"\]' "${definitions}"; then
    echo "explicit Terraform apply replacement is missing" >&2
    exit 1
fi
