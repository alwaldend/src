#!/usr/bin/env sh

set -eux

result="#!/usr/bin/env sh\n"
result+="set -eux\n"
output=${1}
shift
for file in "${@}"; do
    base64=$(base64 --wrap 0 "${file}") 
    name=$(basename "${file}")
    name="${name%.sh}"
    if [ "${name}" = "bashrc"; then
        target="\${HOME}/.bashrc"
    else
        target="\${HOME}/.local/bin/${name}"
    fi
    result+="\n# ${name}\n"
    result+="mkdir -p \"\$(dirname \"${target}\")\"\n"
    result+="base64='${base64}'\n"
    result+="base64 --decode <<<\"${base64}\" >\"${target}\"\n"
    result+="chmod u+x \"${target}\""
done
echo -e "${result}" >"${output}"
