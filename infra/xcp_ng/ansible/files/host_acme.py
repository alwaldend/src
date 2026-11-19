"""Manage XCP-ng ACME certificates without adding packages to dom0 (Python 3.6)."""

import base64
import datetime
import fcntl
import hashlib
import json
import os
from pathlib import Path
import re
import ssl
import subprocess
import sys
import tarfile
import tempfile
import urllib.request


ROOT = Path("/etc/xcp-ng-acme")
CERT = ROOT / "certificates/host1.crt"
KEY = ROOT / "certificates/host1.key"
CHAIN = ROOT / "certificates/host1.issuer.crt"
CURRENT = Path("/etc/xensource/xapi-ssl.pem")


def write_file(path, content, mode):
    if path.is_symlink():
        raise RuntimeError(
            "Refusing to replace a symlink in the ACME installation"
        )
    if path.exists() and path.read_bytes() == content:
        changed = (
            path.stat().st_mode & 0o777
        ) != mode or path.stat().st_uid != 0
        os.chmod(str(path), mode)
        os.chown(str(path), 0, 0)
        return changed
    fd, temporary = tempfile.mkstemp(prefix=".acme-", dir=str(path.parent))
    try:
        with os.fdopen(fd, "wb") as output:
            output.write(content)
            os.fchmod(output.fileno(), mode)
        os.replace(temporary, str(path))
    finally:
        if os.path.exists(temporary):
            os.unlink(temporary)
    return True


def install(config):
    if ROOT.is_symlink():
        raise RuntimeError(
            "The ACME installation directory must not be a symlink"
        )
    ROOT.mkdir(mode=0o700, exist_ok=True)
    os.chmod(str(ROOT), 0o700)
    os.chown(str(ROOT), 0, 0)
    webroot = Path(config["webroot"])
    if not webroot.is_dir() or webroot.stat().st_uid != 0:
        raise RuntimeError(
            "The inspected XAPI webroot must exist and be root-owned"
        )
    download = config.pop("lego_download")
    unit = config.pop("service_unit")
    timer = config.pop("timer_unit")
    ca = config.pop("ca_certificate")
    pin = download["integrity"]
    changed = False
    marker = ROOT / "lego.integrity"
    if (
        not (ROOT / "lego").is_file()
        or not marker.exists()
        or marker.read_text() != pin
    ):
        algorithm, expected = pin.split("-", 1)
        if algorithm != "sha256":
            raise RuntimeError("Expected a SHA256-pinned Lego archive")
        digest = hashlib.sha256()
        with tempfile.TemporaryFile(dir=str(ROOT)) as archive:
            with urllib.request.urlopen(
                download["url"], timeout=60
            ) as response:
                total = 0
                while True:
                    chunk = response.read(1024 * 1024)
                    if not chunk:
                        break
                    total += len(chunk)
                    if total > 32 * 1024 * 1024:
                        raise RuntimeError(
                            "Lego archive exceeded the expected size limit"
                        )
                    digest.update(chunk)
                    archive.write(chunk)
            if base64.b64encode(digest.digest()).decode("ascii") != expected:
                raise RuntimeError("Lego archive checksum mismatch")
            archive.seek(0)
            with tarfile.open(fileobj=archive, mode="r:gz") as source:
                for name, mode in (("lego", 0o700), ("LICENSE", 0o600)):
                    member = source.getmember(name)
                    if not member.isfile() or member.size > 128 * 1024 * 1024:
                        raise RuntimeError("Unexpected Lego archive member")
                    changed |= write_file(
                        ROOT / name, source.extractfile(member).read(), mode
                    )
        changed |= write_file(marker, pin.encode("ascii"), 0o600)
    subprocess.run(
        [str(ROOT / "lego"), "--version"], check=True, stdout=subprocess.PIPE
    )
    changed |= write_file(
        ROOT / "host_acme.py", Path(__file__).read_bytes(), 0o700
    )
    changed |= write_file(ROOT / "ca.crt", ca.encode("utf-8"), 0o600)
    changed |= write_file(
        ROOT / "config.json",
        json.dumps(config, sort_keys=True).encode("utf-8"),
        0o600,
    )
    changed |= write_file(
        Path("/etc/systemd/system/xcp-ng-acme-renew.service"),
        unit.encode("utf-8"),
        0o644,
    )
    changed |= write_file(
        Path("/etc/systemd/system/xcp-ng-acme-renew.timer"),
        timer.encode("utf-8"),
        0o644,
    )
    if changed:
        subprocess.run(["/usr/bin/systemctl", "daemon-reload"], check=True)
    print(
        json.dumps(
            {
                "changed": changed,
                "certificate_exists": CERT.is_file(),
                "account_exists": account_exists(config),
            }
        )
    )


def account_exists(config):
    for filename in (ROOT / "accounts").rglob("*.json"):
        account = json.loads(filename.read_text())
        if (
            account.get("id") == "src_infra_xcp_ng_host1"
            and account.get("server") == config["server"]
            and account.get("registration")
        ):
            return True
    return False


def fingerprint(path):
    if not path.is_file():
        return None
    return subprocess.check_output(
        [
            "/usr/bin/openssl",
            "x509",
            "-in",
            str(path),
            "-noout",
            "-fingerprint",
            "-sha256",
        ]
    ).strip()


def record_failure(error, eab=None):
    detail = str(error)
    if eab is not None:
        for name in ("id", "key"):
            if eab.get(name):
                detail = detail.replace(eab[name], "[REDACTED]")
    detail = re.sub(
        r"-----BEGIN [^-]+-----.*?-----END [^-]+-----",
        "[REDACTED PEM]",
        detail,
        flags=re.DOTALL,
    )
    classification = "unknown"
    for needle, value in (
        ("x509", "tls_validation"),
        ("certificate verify failed", "tls_validation"),
        ("external account", "external_account_binding"),
        ("eab", "external_account_binding"),
        ("unauthorized", "acme_authorization"),
        ("rejectedidentifier", "domain_policy"),
        ("timeout", "timeout"),
        ("timed out", "timeout"),
        ("server_certificate", "xapi_certificate_validation"),
        ("flag", "client_arguments"),
    ):
        if needle in detail.lower():
            classification = value
            break
    status = re.search(r"\(status (\d+)\)", detail)
    stage = (
        "xapi_install"
        if detail.startswith("XAPI certificate installation")
        else "acme_client"
    )
    receipt = {
        "observed_at": datetime.datetime.utcnow().isoformat() + "Z",
        "classification": classification,
        "stage": stage,
        "returncode": int(status.group(1)) if status else None,
        "error_type": type(error).__name__,
        "detail": detail[:12000],
    }
    write_file(
        ROOT / "last-error.json", json.dumps(receipt).encode("utf-8"), 0o600
    )


def renew(eab=None):
    try:
        result = renew_certificate(eab)
    except Exception as error:
        record_failure(error, eab)
        raise RuntimeError(
            "Host ACME failed; inspect /etc/xcp-ng-acme/last-error.json"
        ) from None
    receipt = ROOT / "last-error.json"
    if receipt.exists():
        receipt.unlink()
    return result


def renew_certificate(eab=None):
    os.umask(0o077)
    with (ROOT / "renew.lock").open("a") as lock:
        fcntl.flock(lock, fcntl.LOCK_EX)
        config = json.loads((ROOT / "config.json").read_text())
        environment = os.environ.copy()
        environment["LEGO_CA_CERTIFICATES"] = str(ROOT / "ca.crt")
        for name in ("LEGO_EAB", "LEGO_EAB_KID", "LEGO_EAB_HMAC"):
            environment.pop(name, None)
        if eab is not None:
            environment.update(
                LEGO_EAB="true",
                LEGO_EAB_KID=eab["id"],
                LEGO_EAB_HMAC=eab["key"],
            )
        command = [
            str(ROOT / "lego"),
            "run",
            "--path",
            str(ROOT),
            "--cert.name",
            "host1",
            "--account-id",
            "src_infra_xcp_ng_host1",
            "--accept-tos",
            "--key-type",
            "RSA4096",
            "--no-bundle",
            "--force-cert-domains",
            "--http",
            "--http.webroot",
            config["webroot"],
            "--server",
            config["server"],
            "--no-random-sleep",
        ]
        for domain in config["domains"]:
            command.extend(["--domains", domain])
        issued_before = fingerprint(CERT)
        result = subprocess.run(
            command,
            env=environment,
            stdout=subprocess.PIPE,
            stderr=subprocess.STDOUT,
            timeout=240,
        )
        installed = False
        # Reconcile even after an ACME failure: a prior issuance may have succeeded
        # while installing its certificate failed. No new order is needed to retry.
        if CERT.is_file() and fingerprint(CERT) != fingerprint(CURRENT):
            installation = subprocess.run(
                [
                    "/opt/xensource/bin/xe",
                    "host-server-certificate-install",
                    "certificate=" + str(CERT),
                    "private-key=" + str(KEY),
                    "certificate-chain=" + str(CHAIN),
                ],
                stdout=subprocess.PIPE,
                stderr=subprocess.STDOUT,
                timeout=120,
            )
            if installation.returncode != 0:
                raise RuntimeError(
                    "XAPI certificate installation failed (status {}): {}".format(
                        installation.returncode,
                        installation.stdout.decode("utf-8", errors="replace"),
                    )
                )
            installed = True
        if result.returncode != 0:
            raise RuntimeError(
                "Lego certificate request failed (status {}): ".format(
                    result.returncode
                )
                + result.stdout.decode("utf-8", errors="replace")
            )
        if fingerprint(CERT) != fingerprint(CURRENT):
            raise RuntimeError(
                "XAPI certificate does not match the issued certificate"
            )
        result = {
            "changed": installed or issued_before != fingerprint(CERT),
            "installed": installed,
        }
        print(json.dumps(result))
        return result


def diagnose():
    receipt = ROOT / "last-error.json"
    result = {
        "certificate_exists": CERT.is_file(),
        "error_receipt_exists": receipt.is_file(),
    }
    if receipt.is_file():
        error = json.loads(receipt.read_text())
        result.update(
            {
                key: error[key]
                for key in (
                    "observed_at",
                    "classification",
                    "error_type",
                    "stage",
                    "returncode",
                )
            }
        )
    print(json.dumps(result))


def verify():
    config = json.loads((ROOT / "config.json").read_text())
    context = ssl.create_default_context(cafile=str(ROOT / "ca.crt"))
    for domain in config["domains"]:
        with urllib.request.urlopen(
            "https://" + domain + "/", context=context, timeout=20
        ) as response:
            if response.status != 200:
                raise RuntimeError("Unexpected XAPI HTTPS status")
    enabled = (
        subprocess.run(
            [
                "/usr/bin/systemctl",
                "is-enabled",
                "--quiet",
                "xcp-ng-acme-renew.timer",
            ],
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
        ).returncode
        == 0
    )
    active = (
        subprocess.run(
            [
                "/usr/bin/systemctl",
                "is-active",
                "--quiet",
                "xcp-ng-acme-renew.timer",
            ],
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
        ).returncode
        == 0
    )
    subprocess.run(
        ["/usr/bin/systemctl", "enable", "--now", "xcp-ng-acme-renew.timer"],
        check=True,
    )
    print(
        json.dumps(
            {
                "verified_domains": config["domains"],
                "timer_enabled": True,
                "changed": not enabled or not active,
            }
        )
    )


if __name__ == "__main__":
    os.umask(0o077)
    if os.geteuid() != 0:
        raise RuntimeError("XCP-ng certificate management requires root")
    if sys.argv[1] == "install":
        install(json.loads(sys.argv[2]))
    elif sys.argv[1] == "renew":
        renew()
    elif sys.argv[1] == "renew-result":
        try:
            result = renew()
        except Exception:
            print(json.dumps({"success": False, "changed": False}))
        else:
            print(json.dumps({"success": True, "changed": result["changed"]}))
    elif sys.argv[1] == "verify":
        verify()
    elif sys.argv[1] == "diagnose":
        diagnose()
    else:
        raise RuntimeError("Unknown XCP-ng certificate operation")
