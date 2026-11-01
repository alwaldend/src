#!/bin/sh
# This script only checks prerequisites; the playbook owns the transaction.
set -eu
PATH=/usr/sbin:/usr/bin:/sbin:/bin
LC_ALL=C
export PATH LC_ALL
reject() {
    printf 'UPDATE_BLOCKED=%s\n' "$1"
    exit 1
}
query() {
    timeout --signal=TERM --kill-after=5s 30s xe "$@"
}
test "$(id -u)" = 0 || reject root_required
test -x /usr/bin/yum || reject yum_missing
test -x /opt/xensource/bin/xe-toolstack-restart || reject toolstack_restart_missing
release=$(rpm -q --qf '%{VERSION}' xcp-ng-release)
case "$release" in
8.3 | 8.3.0) ;;
*) reject unsupported_xcpng_release ;;
esac
host=$(query host-list --minimal)
printf '%s\n' "$host" | grep -Eq '^[0-9a-f-]{36}$' || reject single_host_pool_required
pool=$(query pool-list --minimal)
printf '%s\n' "$pool" | grep -Eq '^[0-9a-f-]{36}$' || reject single_pool_required
test "$(query pool-param-get uuid="$pool" param-name=master)" = "$host" || reject pool_coordinator_required
test "$(query pool-param-get uuid="$pool" param-name=ha-enabled)" = false || reject ha_enabled
host_enabled=$(query host-param-get uuid="$host" param-name=enabled)
case "$host_enabled" in
true | false) ;;
*) reject unknown_host_enabled_state ;;
esac
pending_tasks=$(query task-list status=pending --minimal) || reject pending_task_query_failed
cancelling_tasks=$(query task-list status=cancelling --minimal) || reject cancelling_task_query_failed
xostor=$(query sr-list type=linstor --minimal) || reject storage_query_failed
attached_media=$(query vbd-list type=CD empty=false currently-attached=true --minimal) || reject attached_media_query_failed
test -z "$pending_tasks" || reject pending_xapi_tasks
test -z "$cancelling_tasks" || reject cancelling_xapi_tasks
test -z "$xostor" || reject xostor_requires_coordinated_maintenance
test -z "$attached_media" || reject attached_cd_media
if rpm -q yum-cron >/dev/null 2>&1; then
    reject installed_yum_cron_requires_review
fi
if test "${XCPNG_USE_ORIGIN_MIRROR:-no}" = yes; then
    /usr/bin/python -Es - <<'APPLIANCE_ORIGIN_CHECK_EOF'
"""Read native YUM configuration without setting up repository metadata."""
import yum

_base = yum.YumBase()  # type: yum.YumBase
"""Vendor Python 2 YUM loader, including effective reposdir configuration."""
for _repository_id in ('xcp-ng-base', 'xcp-ng-updates'):
    _repository = _base.repos.getRepo(_repository_id)  # type: yum.yumRepo.YumRepository
    """Configured repository before a transaction overrides its baseurl."""
    if _repository.mirrorlist or _repository.metalink:
        raise SystemExit('UPDATE_BLOCKED=origin_override_requires_baseurl_only')
APPLIANCE_ORIGIN_CHECK_EOF
fi
printf 'HOST_UUID=%s\nHOST_ENABLED=%s\n' "$host" "$host_enabled"
printf 'UPDATE_PREFLIGHT_COMPLETE\n'
