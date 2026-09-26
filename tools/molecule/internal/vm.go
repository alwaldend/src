package internal

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func identity(pid int) processIdentity {
	data, err := os.ReadFile(fmt.Sprintf("/proc/%d/stat", pid))
	if err != nil {
		return processIdentity{}
	}
	end := strings.LastIndex(string(data), ")")
	if end < 0 {
		return processIdentity{}
	}
	fields := strings.Fields(string(data)[end+1:])
	if len(fields) < 20 || fields[0] == "Z" {
		return processIdentity{}
	}
	return processIdentity{PID: pid, Start: fields[19]}
}

func alive(p processIdentity) bool { return p.PID > 1 && p.Start != "" && identity(p.PID) == p }

func loadOwned(root string) (manifest, error) {
	var m manifest
	abs, err := filepath.Abs(root)
	if err != nil {
		return m, err
	}
	info, err := os.Lstat(abs)
	if err != nil {
		return m, err
	}
	if !info.IsDir() || info.Mode()&0o077 != 0 {
		return m, errors.New("refusing non-private or symlinked run directory")
	}
	stat, ok := info.Sys().(*syscall.Stat_t)
	if !ok || stat.Uid != uint32(os.Getuid()) {
		return m, errors.New("run directory owner mismatch")
	}
	if err := readJSON(filepath.Join(abs, "manifest.json"), &m); err != nil {
		return m, err
	}
	if m.Version != 1 || m.Root != abs || len(m.RunID) != 32 || !strings.HasPrefix(filepath.Base(abs), "molecule-") {
		return m, errors.New("invalid ownership manifest")
	}
	namespace, err := os.Readlink("/proc/self/ns/pid")
	if err != nil {
		return m, err
	}
	if namespace != m.PIDNamespace {
		return m, errors.New("run belongs to a different PID namespace; recover inside the original namespace or let Bazel destroy its sandbox")
	}
	return m, nil
}

func create(root string) error {
	m, err := loadOwned(root)
	if err != nil {
		return err
	}
	if alive(m.QEMU) {
		return nil
	}
	config := m.Config
	deadline := config.BootTimeout
	if config.Accelerator == "tcg" {
		deadline = config.TCGBootTimeout
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(deadline)*time.Second)
	defer cancel()
	if err := command(ctx, &m, "/usr/bin/ssh-keygen", "-q", "-t", "ed25519", "-N", "", "-C", "molecule-test", "-f", filepath.Join(root, "ssh_key")).Run(); err != nil {
		return fmt.Errorf("generate disposable key: %w", err)
	}
	pub, err := os.ReadFile(filepath.Join(root, "ssh_key.pub"))
	if err != nil {
		return err
	}
	userData := fmt.Sprintf("#cloud-config\npreserve_hostname: true\nusers:\n  - name: molecule\n    groups: [wheel]\n    sudo: ['ALL=(ALL) NOPASSWD:ALL']\n    shell: /bin/bash\n    lock_passwd: true\n    ssh_authorized_keys:\n      - %s\nssh_pwauth: false\ndisable_root: true\n", strings.TrimSpace(string(pub)))
	if err := os.WriteFile(filepath.Join(root, "user-data"), []byte(userData), 0o600); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(root, "meta-data"), []byte("instance-id: "+m.RunID+"\n"), 0o600); err != nil {
		return err
	}
	if err := command(ctx, &m, config.Tools["seed"], filepath.Join(root, "seed.iso"), filepath.Join(root, "user-data"), filepath.Join(root, "meta-data")).Run(); err != nil {
		return fmt.Errorf("create seed: %w", err)
	}
	if err := command(ctx, &m, config.Tools["qemu-img"], "create", "-f", "qcow2", "-F", "qcow2", "-b", config.Image, filepath.Join(root, "root.qcow2"), strconv.Itoa(config.DiskGB)+"G").Run(); err != nil {
		return fmt.Errorf("create root overlay: %w", err)
	}
	for i, size := range config.DataDisksGB {
		if err := command(ctx, &m, config.Tools["qemu-img"], "create", "-f", "qcow2", filepath.Join(root, fmt.Sprintf("data-%d.qcow2", i)), strconv.Itoa(size)+"G").Run(); err != nil {
			return fmt.Errorf("create virtual data disk: %w", err)
		}
	}
	if len(config.Ports) == 0 {
		return errors.New("guest ports must include SSH port 22")
	}
	for attempt := 0; attempt < 5; attempt++ {
		var reservations []net.Listener
		network := "user,id=network"
		for _, port := range config.Ports {
			if port < 1 || port > 65535 {
				return errors.New("guest port out of range")
			}
			l, err := net.Listen("tcp4", "127.0.0.1:0")
			if err != nil {
				return err
			}
			reservations = append(reservations, l)
			m.Ports[strconv.Itoa(port)] = l.Addr().(*net.TCPAddr).Port
			network += fmt.Sprintf(",hostfwd=tcp:127.0.0.1:%d-:%d", m.Ports[strconv.Itoa(port)], port)
		}
		if m.Ports["22"] == 0 {
			for _, l := range reservations {
				l.Close()
			}
			return errors.New("guest ports must include SSH port 22")
		}
		cpu := "max"
		if config.Accelerator == "kvm" {
			cpu = "host"
		}
		args := []string{"-name", "molecule-" + m.RunID, "-machine", "q35", "-accel", config.Accelerator, "-cpu", cpu, "-smp", strconv.Itoa(config.CPUs), "-m", strconv.Itoa(config.MemoryMB), "-nodefaults", "-display", "none", "-serial", "file:serial.log", "-qmp", "unix:qmp.sock,server=on,wait=off", "-pidfile", "qemu.pid", "-daemonize", "-drive", "file=root.qcow2,format=qcow2,if=virtio", "-drive", "file=seed.iso,format=raw,if=virtio,readonly=on", "-netdev", network, "-device", "virtio-net-pci,netdev=network,romfile="}
		for i := range config.DataDisksGB {
			args = append(args, "-drive", fmt.Sprintf("file=data-%d.qcow2,format=qcow2,if=virtio", i))
		}
		// Persist ownership before launch. Recovery can validate a pidfile written
		// during the narrow interval before the start time is recorded.
		if err := updateManifest(root, func(current *manifest) { current.Ports = m.Ports }); err != nil {
			return err
		}
		for _, l := range reservations {
			l.Close()
		}
		output, launchErr := command(ctx, &m, config.Tools["qemu"], args...).CombinedOutput()
		if launchErr != nil {
			if strings.Contains(string(output), "Could not set up host forwarding rule") {
				continue
			}
			return fmt.Errorf("QEMU launch failed: %s: %w", strings.TrimSpace(string(output)), launchErr)
		}
		pidBytes, err := os.ReadFile(filepath.Join(root, "qemu.pid"))
		if err != nil {
			return err
		}
		pid, err := strconv.Atoi(strings.TrimSpace(string(pidBytes)))
		if err != nil {
			return err
		}
		m.QEMU = identity(pid)
		if !alive(m.QEMU) {
			return errors.New("QEMU exited before boot")
		}
		if err := updateManifest(root, func(current *manifest) { current.QEMU = m.QEMU }); err != nil {
			return err
		}
		for {
			if !alive(m.QEMU) {
				return errors.New("QEMU exited while waiting for SSH")
			}
			if err := sshCommand(ctx, &m, "true").Run(); err == nil {
				break
			}
			select {
			case <-ctx.Done():
				return errors.New("guest SSH boot deadline exceeded")
			case <-time.After(2 * time.Second):
			}
		}
		inventory := map[string]any{"all": map[string]any{"hosts": map[string]any{"localhost": map[string]any{"ansible_connection": "local", "ansible_remote_tmp": filepath.Join(root, "tmp/local")}, "molecule": map[string]any{"ansible_host": "127.0.0.1", "ansible_port": m.Ports["22"], "ansible_user": "molecule", "ansible_ssh_private_key_file": filepath.Join(root, "ssh_key"), "ansible_python_interpreter": "/usr/bin/python3", "molecule_ports": m.Ports, "molecule_state": root}}}}
		return writeJSON(filepath.Join(root, "project/inventory/hosts.yml"), inventory)
	}
	return errors.New("could not allocate loopback forwarding ports after five attempts")
}

func sshCommand(ctx context.Context, m *manifest, remote string) *exec.Cmd {
	args := []string{"-F", "/dev/null", "-o", "BatchMode=yes", "-o", "StrictHostKeyChecking=no", "-o", "UserKnownHostsFile=/dev/null", "-o", "IdentitiesOnly=yes", "-o", "ConnectTimeout=5", "-o", "LogLevel=ERROR", "-i", filepath.Join(m.Root, "ssh_key"), "-p", strconv.Itoa(m.Ports["22"]), "molecule@127.0.0.1", remote}
	return command(ctx, m, "/usr/bin/ssh", args...)
}

func stop(root string, remove bool) error {
	m, err := loadOwned(root)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	if remove && alive(m.Supervisor) && m.Supervisor.PID != os.Getpid() {
		return errors.New("refusing recovery while the supervisor is still alive; signal it instead")
	}
	// Validate both kernel process start time and a unique QEMU command token.
	if m.QEMU.PID == 0 {
		if b, err := os.ReadFile(filepath.Join(root, "qemu.pid")); err == nil {
			pid, _ := strconv.Atoi(strings.TrimSpace(string(b)))
			m.QEMU = identity(pid)
		}
	}
	if alive(m.QEMU) {
		args, err := os.ReadFile(fmt.Sprintf("/proc/%d/cmdline", m.QEMU.PID))
		if err != nil {
			return err
		}
		if !strings.Contains(string(args), "\x00molecule-"+m.RunID+"\x00") {
			return errors.New("QEMU ownership mismatch; refusing to terminate process")
		}
		// QEMU daemonization changes cwd to /. Verify its open overlay instead.
		fds, err := os.ReadDir(fmt.Sprintf("/proc/%d/fd", m.QEMU.PID))
		if err != nil {
			return err
		}
		ownedDisk := false
		for _, fd := range fds {
			path, _ := os.Readlink(fmt.Sprintf("/proc/%d/fd/%s", m.QEMU.PID, fd.Name()))
			if path == filepath.Join(m.Root, "root.qcow2") {
				ownedDisk = true
			}
		}
		if !ownedDisk {
			return errors.New("QEMU overlay ownership mismatch")
		}
		if err := terminate(m.QEMU, false); err != nil {
			return err
		}
	}
	if remove {
		// A killed supervisor can leave descendants after the group leader has
		// died. Match the captured process group and run-specific environment,
		// not a process name or an unverified recycled group ID.
		entries, err := os.ReadDir("/proc")
		if err != nil {
			return err
		}
		for _, entry := range entries {
			pid, err := strconv.Atoi(entry.Name())
			if err != nil || pid <= 1 {
				continue
			}
			data, err := os.ReadFile(filepath.Join("/proc", entry.Name(), "stat"))
			if err != nil {
				continue
			}
			end := strings.LastIndex(string(data), ")")
			if end < 0 {
				continue
			}
			fields := strings.Fields(string(data)[end+1:])
			if len(fields) < 20 || fields[2] != strconv.Itoa(m.Child.PID) {
				continue
			}
			start, _ := strconv.ParseUint(fields[19], 10, 64)
			leaderStart, _ := strconv.ParseUint(m.Child.Start, 10, 64)
			if leaderStart == 0 || start < leaderStart {
				continue
			}
			env, err := os.ReadFile(filepath.Join("/proc", entry.Name(), "environ"))
			if err != nil {
				continue
			}
			if !bytes.Contains(append([]byte{0}, env...), []byte("\x00MOLECULE_RUNNER_STATE="+m.Root+"\x00")) {
				continue
			}
			if err := terminate(processIdentity{PID: pid, Start: fields[19]}, false); err != nil {
				return err
			}
		}
	}
	if remove {
		return os.RemoveAll(m.Root)
	}
	return nil
}

func terminate(p processIdentity, group bool) error {
	if !alive(p) {
		return nil
	}
	pid := p.PID
	if group {
		pid = -pid
	}
	if err := syscall.Kill(pid, syscall.SIGTERM); err != nil && err != syscall.ESRCH {
		return err
	}
	for i := 0; i < 50; i++ {
		if !alive(p) {
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
	if alive(p) {
		if err := syscall.Kill(pid, syscall.SIGKILL); err != nil && err != syscall.ESRCH {
			return err
		}
	}
	for i := 0; i < 50; i++ {
		if !alive(p) {
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
	return errors.New("owned process survived cleanup deadline")
}
