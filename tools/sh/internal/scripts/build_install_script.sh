#!/usr/bin/env sh

set -eu

archive="${1}"
output=${2}

cat - >"${output}" <<EOF
#!/usr/bin/env sh
set -eux
temp=\$(mktemp -d)
archive="\${temp}/archive.tar"
trap 'rm -r "\${temp}"' EXIT

echo '$(base64 --wrap 0 "${archive}")' | base64 --decode >"\${archive}"
tar -xf "\${archive}" -C "\${temp}"
for script in "\${temp}"/*.sh; do
    name=\$(basename "\${script}")
    name="\${name%.sh}"
    if [ "\${name}" = "bashrc" ]; then
        target="\${HOME}/.bashrc"
    else
        target="\${HOME}/.local/bin/\${name}"
    fi
    mkdir -p "\$(dirname "\${target}")"
    mv "\${script}" "\${target}"
    chmod u+x "\${target}"
done
EOF
