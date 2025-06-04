#!/usr/bin/env sh
set -o errexit -o nounset -o xtrace

output_path="${1:?missing path}"
archive="${2:?missing tar location}"
install="${3:?missing install script location}"

archive_base64=$(base64 -w 0 <"${archive}")
install_base64=$(base64 -w 0 <"${install}")

cat - >"${output_path}" <<EOF
#!/usr/bin/env sh
set -o errexit -o nounset -o xtrace
install=\$(mktemp)
trap "rm -rf '\${install}'" EXIT
archive=\$(mktemp)
trap "rm -rf '\${archive}'" EXIT
chmod u+x "\${install}"
install_base64='${install_base64}'
archive_base64='${archive_base64}'

echo "\${install_base64}" | base64 -d >"\${install}"
echo "\${archive_base64}" | base64 -d >"\${archive}"
"\${install}" "\${archive}" "\${@}"
EOF
