#!/usr/bin/env sh

nvim_lazy_sync() {
    nvim --headless "+Lazy! sync" +qa
    nvim_lazy_update_lockfile
}

nvim_lazy_update_lockfile() {
    _lock_path=".config/nvim/lazy-lock.json"
    _lockfile_source="${HOME}/.local/share/dotfiles/stow/${_lock_path}"
    _lockfile_target="$(git rev-parse --show-toplevel)/dotfiles/nvim/${_lock_path}"
    cp -f "${_lockfile_source}" "${_lockfile_target}"
}

d_ssh_keygen() {
    _name="${1:?missing key name}"
    ssh-keygen -t ed25519 -f "${HOME}/.ssh/${_name}" -C "${_name}"
}

#######################################
# Run make in parallel
# Arguments:
#   Additional arguments for make
#######################################
d_make() {
    NUMCPUS="$(grep -c '^processor' /proc/cpuinfo)"
    export NUMCPUS
    time nice make "-j${NUMCPUS}" "--load-average=${NUMCPUS}" "${@}"
}

#######################################
# Check if you are running in a container
# Returns:
#   0 if you are in a container, 1 otherwise
#######################################
d_is_container() {
    grep container_t </proc/1/attr/current
}

#########################################
# Check for an available venv and activate it
# Args:
#   None
#######################################
d_venv() {
    poetry_venv=$(poetry env info -p 2>/dev/null || true)
    for venv in "${poetry_venv}" ./.venv ./venv ~/.venv ~/venv; do
        if [ -s "${venv}/bin/activate" ]; then
            . "${venv}/bin/activate"
            break
        fi
    done
}

#######################################
# Get a docker image with digest
# Args:
#   1: image tag
# Outputs:
#   image with digest
#######################################
d_docker_image_with_digest() {
    image="${1:?missing image}"
    digest=$(docker inspect --format='{{index .RepoDigests 0}}' "${image}")
    echo "${image}@${digest##*@}"
}

#######################################
# Start a tmux session
# If it doesn't exist, start one
# Globals:
#   TERM_PROGRAM - used to check if running inside tmux
#######################################
d_tmux_start() {
    session_name="DEFAULT"
    if [ "${TERM_PROGRAM}" = "tmux" ]; then
        return
    fi
    if ! tmux has-session -t "${session_name}"; then
        tmux new-session -d -s "${session_name}"
        tmux send-keys tmux-new-window Enter
    fi
    tmux attach-session -t "${session_name}"
}

#######################################
# PROMPT_COMMAND, run before every prompt
#######################################
d_prompt_command() {
    # append history lines from this session to the history file
    history -a
}

#######################################
# Show RAM info
# Outputs:
#   RAM info
#######################################
d_info_ram() {
    sudo dmidecode --type memory
}

#######################################
# Source scripts
# Args:
#   Script paths
#######################################
d_source_scripts() {
    for script in "${@}"; do
        if [ -s "${script}" ]; then
            . "${script}"
        fi
    done
}

#######################################
# Add a path to the front of PATH
# Globals:
#   PATH
# Args:
#   paths that will be added to PATH
#######################################
d_add_to_path_front() {
    for path in "${@}"; do
        case ":${PATH:=${path}}:" in
        *":${path}:"*) ;;
        *) PATH="${path}:${PATH}" ;;
        esac
    done
}

## create a certificate secret
d_kubernetes_secret_create_sertificate() {
    name="${1:?missing name}"
    namespace="${2:-default}"
    key="${3:-${name}.key}"
    certificate="${4:-${name}.crt}"

    kubectl create secret tls "${name}" \
        --key "${key}" \
        --cert "${certificate}" \
        --namespace "${namespace}"

}

## modify current context
d_kubernetes_context() {
    ns=$(kubectl get ns | fzf)
    kubectl config set-context --current --namespace="${ns}"
}

## create a service account and everything needed
d_kubernetes_service_account_full() {
    namespace="${1}"
    account="${2}"
    role="${3}"
    rolebinding="${4}"
    roletype="${5:-Role}"

    kubernetes_service_account "${namespace}" "${account}"
    kubernetes_rolebinding "${namespace}" "${account}" "${role}" "${rolebinding}"
}

d_kubernetes_rolebinding() {
    namespace="${1}"
    account="${2}"
    role="${3}"
    rolebinding="${4:-${role}-rolebinding}"
    roletype="${5:-Role}"
    apigroup="${6:-rbac.authorization.k8s.io}"
    if [ "${roletype}" = "ClusterRole" ]; then
        apigroup="cluster.${apigroup}"
    fi

    dotfiles_kubernetes_apply "
        apiVersion: rbac.authorization.k8s.io/v1
        kind: RoleBinding
        metadata:
          name: ${rolebinding}
          namespace: ${namespace}
        roleRef:
          apiGroup: ${apigroup}
          kind: Role
          name: ${role}
        subjects:
        - namespace: ${namespace}
          kind: ServiceAccount
          name: ${account}
    "
}

d_kubernetes_service_account() {
    namespace="${1}"
    account="${2}"
    dotfiles_kubernetes_apply "
        apiVersion: v1
        kind: ServiceAccount
        metadata:
            name: ${account}
            namespace: ${namespace}
    "
}

d_kubernetes_apply() {
    echo "${1:?missing yaml}" | kubectl apply -f -
}

# get api url
d_kubernetes_api() {
    kubectl cluster-info | grep 'Kubernetes master' | awk '/http/ { print $NF }'
}

## check open ports
d_net_ports_list() {
    sudo ss -tulpn | grep LISTEN
}

d_net_connections_list() {
    sudo netstat -nputw
}

d_net_check_dns() {
    dig "${@}" +nostats +nocomments +nocmd
}

# add windows 10 uefi entry to grub
d_setup_grub_add_windows_10_uefi() {
    # exec tail -n +4 $0
    # this line needs to be in the file, without it
    # commands will not be recognized
    source_grub
    echo "input where the EFI partition is mounted"
    read -r partition
    fs_uuid=$(sudo grub-probe --target=fs_uuid "${partition}/EFI/Microsoft/Boot/bootmgfw.efi")
    entry=$(
        cat - <<EOF
        # Windows 10 EFI entry
        if [ "\${grub_platform}" == "efi" ]; then
            menuentry "Microsoft Windows 10 UEFI" {
                insmod part_gpt
                insmod fat
                insmod search_fs_uuid
                insmod chain
                search --fs-uuid --set=root "${fs_uuid}"
                chainloader /EFI/Microsoft/Boot/bootmgfw.efi
            }
        fi
EOF
    )
    echo "${entry}"
    echo
    echo "is that okay?"
    read -r answer
    if [ "${answer}" != "n" ] && [ "${answer}" != "N" ]; then
        echo "${entry}" | sudo tee -a "/etc/grub.d/40_custom"
        source_grub
    else
        setup_grub_add_windows_10_uefi
    fi
}

if [ "${#}" -gt 0 ]; then
    set -o errexit -o nounset -o xtrace
    "${@}"
fi
