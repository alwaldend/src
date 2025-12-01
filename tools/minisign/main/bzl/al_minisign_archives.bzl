AL_MINISIGN_ARCHIVES = {
    "0.12": {
        "linux": {
            "url": "https://github.com/jedisct1/minisign/releases/download/0.12/minisign-0.12-linux.tar.gz",
            "integrity": "sha256-mlmbSLput7HoDxLza5TOynwAt6UXPJXD78iNmCKVfnM=",
            "toolchains": [
                {
                    "name": "toolchain_x86_64",
                    "binary": "minisign-linux/x86_64/minisign",
                    "exec_compatible_with": [
                        "@platforms//os:linux",
                        "@platforms//cpu:x86_64",
                    ],
                },
                {
                    "name": "toolchain_aarch64",
                    "binary": "minisign-linux/aarch64/minisign",
                    "exec_compatible_with": [
                        "@platforms//os:linux",
                        "@platforms//cpu:aarch64",
                    ],
                },
            ],
        },
        "win64": {
            "url": "https://github.com/jedisct1/minisign/releases/download/0.12/minisign-0.12-win64.zip",
            "integrity": "sha256-N7YANE4gwZMUsugoE9sr/cxAi3e4dvdyeInb1G1TlHk=",
            "toolchains": [
                {
                    "name": "toolchain_x86_64",
                    "binary": "minisign-win64/x86_64/minisign.exe",
                    "exec_compatible_with": [
                        "@platforms//os:windows",
                        "@platforms//cpu:x86_64",
                    ],
                },
                {
                    "name": "toolchain_aarch64",
                    "binary": "minisign-win64/aarch64/minisign.exe",
                    "exec_compatible_with": [
                        "@platforms//os:windows",
                        "@platforms//cpu:aarch64",
                    ],
                },
            ],
        },
        "macos": {
            "url": "https://github.com/jedisct1/minisign/releases/download/0.12/minisign-0.12-macos.zip",
            "integrity": "sha256-iQALGVNXZfnP/GWmXWSoIPQz7224AgZn91cOBr9qrGM=",
            "toolchains": [
                {
                    "name": "toolchain",
                    "binary": "minisign",
                    "exec_compatible_with": [
                        "@platforms//os:macos",
                        "@platforms//cpu:x86_64",
                    ],
                },
            ],
        },
    },
}
