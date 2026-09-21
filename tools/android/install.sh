#!/usr/bin/env sh

set -eu

if [ -z "${ANDROID_HOME:-}" ]; then
    echo "ANDROID_HOME must point to an Android SDK directory" >&2
    exit 1
fi

if [ -z "${RUNFILES_DIR:-}" ] && [ -d "${0}.runfiles" ]; then
    RUNFILES_DIR="${0}.runfiles"
    export RUNFILES_DIR
fi

sdkmanager="${1}"
packages="${2}"
shift 2

exec "${sdkmanager}" \
    --sdk_root="${ANDROID_HOME}" \
    --package_file="${packages}" \
    "${@}"
