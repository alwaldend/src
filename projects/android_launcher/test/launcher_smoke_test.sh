#!/usr/bin/env bash
set -euo pipefail

if [ -z "${RUNFILES_DIR:-}" ]; then
    RUNFILES_DIR="$0.runfiles"
fi

adb="${ADB:-${ANDROID_HOME:-$HOME/Android/Sdk}/platform-tools/adb}"
if [ ! -x "$adb" ]; then
    echo >&2 "adb not found at $adb; set ADB or install Android platform-tools"
    exit 1
fi

apk="$(find -L "${RUNFILES_DIR:?}" -type f -name '*.apk' | grep -v '_unsigned' | head -1)"
if [ -z "$apk" ]; then
    apk="$(find -L "${RUNFILES_DIR:?}" -type f -name '*.apk' | head -1)"
fi
if [ -z "$apk" ]; then
    echo >&2 "launcher APK not found in runfiles"
    exit 1
fi

package_name="${PACKAGE_NAME:?PACKAGE_NAME must be set}"
activity_name="${ACTIVITY_NAME:?ACTIVITY_NAME must be set}"

"$adb" devices
"$adb" install -r "$apk"
"$adb" shell am force-stop "$package_name"
"$adb" shell am start -n "$package_name/$activity_name"

sleep 5

pid="$("$adb" shell pidof "$package_name" | tr -d '\r')"
if [ -z "$pid" ]; then
    echo >&2 "Launcher activity did not start successfully"
    exit 1
fi

"$adb" shell am force-stop "$package_name"
