// Package internal supervises isolated Molecule scenarios and their owned guests.
package internal

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"syscall"
	"time"

	"github.com/bazelbuild/rules_go/go/runfiles"
)

type configuration struct {
	Label          string            `json:"label"`
	Scenario       string            `json:"scenario"`
	Files          map[string]string `json:"files"`
	Mappings       map[string]string `json:"mappings"`
	Tools          map[string]string `json:"tools"`
	Lifecycle      map[string]string `json:"lifecycle"`
	Image          string            `json:"image"`
	BIOS           string            `json:"bios"`
	CPUs           int               `json:"cpus"`
	MemoryMB       int               `json:"memory_mb"`
	DiskGB         int               `json:"disk_gb"`
	DataDisksGB    []int             `json:"data_disks_gb"`
	Ports          []int             `json:"ports"`
	Accelerator    string            `json:"accelerator"`
	BootTimeout    int               `json:"boot_timeout"`
	TCGBootTimeout int               `json:"tcg_boot_timeout"`
	AnsibleArgs    []string          `json:"ansible_args"`
	PythonImports  []string          `json:"python_imports"`
	Services       []string          `json:"services"`
}

type phaseResult struct {
	Step    int    `json:"step"`
	Name    string `json:"name"`
	Success bool   `json:"success"`
}

type report struct {
	RunID       string            `json:"run_id"`
	Label       string            `json:"label"`
	Rerun       string            `json:"rerun"`
	RuntimeDir  string            `json:"runtime_dir"`
	Success     bool              `json:"success"`
	Interrupted bool              `json:"interrupted"`
	Accelerator string            `json:"accelerator"`
	Ports       map[string]int    `json:"ports"`
	QEMU        processIdentity   `json:"qemu"`
	Inputs      map[string]string `json:"inputs"`
	Tools       map[string]string `json:"tools"`
	Phases      []phaseResult     `json:"phases"`
	Error       string            `json:"error,omitempty"`
	Cleanup     struct {
		Success bool   `json:"success"`
		Error   string `json:"error,omitempty"`
	} `json:"cleanup"`
}

type manifest struct {
	Version      int             `json:"version"`
	RunID        string          `json:"run_id"`
	Root         string          `json:"root"`
	Phase        string          `json:"phase"`
	Config       configuration   `json:"config"`
	Environment  []string        `json:"environment"`
	Ports        map[string]int  `json:"ports"`
	QEMU         processIdentity `json:"qemu"`
	Child        processIdentity `json:"child"`
	Supervisor   processIdentity `json:"supervisor"`
	PIDNamespace string          `json:"pid_namespace"`
}

type processIdentity struct {
	PID   int    `json:"pid"`
	Start string `json:"start"`
}

func Main(args []string) error {
	if len(args) == 2 && (args[0] == "create" || args[0] == "stop" || args[0] == "destroy") {
		if args[0] == "create" {
			return create(args[1])
		}
		if args[0] == "stop" {
			return stop(args[1], false)
		}
		return stop(args[1], true)
	}
	flags := flag.NewFlagSet("molecule-runner", flag.ContinueOnError)
	configPath := flags.String("config", os.Args[0]+".json", "declared scenario configuration")
	parent := flags.String("state-parent", os.Getenv("TEST_TMPDIR"), "private runtime parent")
	outputs := flags.String("outputs", os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"), "sanitized results directory")
	mode := flags.String("acceptance-mode", "", "runner acceptance injection")
	accelerator := flags.String("accelerator", "", "explicit accelerator override")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if *parent == "" || *outputs == "" {
		return errors.New("run as a Bazel test or supply --state-parent and --outputs")
	}
	if *mode != "" && *mode != "failure" && *mode != "cancellation" {
		return errors.New("unknown acceptance mode")
	}
	var config configuration
	if err := readJSON(*configPath, &config); err != nil {
		return fmt.Errorf("read declared configuration: %w", err)
	}
	config.Label = strings.TrimPrefix(config.Label, "@@")
	if *accelerator != "" {
		config.Accelerator = *accelerator
	}
	if runtime.GOOS != "linux" || runtime.GOARCH != "amd64" {
		return errors.New("requires Linux x86-64")
	}
	if config.CPUs < 1 || config.MemoryMB < 512 || config.DiskGB < 5 {
		return errors.New("invalid guest resource declaration")
	}
	if !filepath.IsLocal(config.Scenario) || filepath.Base(config.Scenario) != config.Scenario {
		return errors.New("scenario name must be a single local directory name")
	}
	for _, service := range config.Services {
		if !regexp.MustCompile(`^[A-Za-z0-9_.@:-]+$`).MatchString(service) {
			return errors.New("invalid diagnostic service name")
		}
	}
	for _, size := range config.DataDisksGB {
		if size < 1 {
			return errors.New("data disk sizes must be positive")
		}
	}
	if err := os.MkdirAll(*parent, 0o700); err != nil {
		return err
	}
	if err := os.MkdirAll(*outputs, 0o700); err != nil {
		return err
	}
	root, err := os.MkdirTemp(*parent, "molecule-")
	if err != nil {
		return err
	}
	root, err = filepath.Abs(root)
	if err != nil {
		return err
	}
	var entropy [16]byte
	if _, err := rand.Read(entropy[:]); err != nil {
		return err
	}
	id := hex.EncodeToString(entropy[:])
	m := manifest{Version: 1, RunID: id, Root: root, Phase: "preflight", Config: config, Ports: map[string]int{}, Supervisor: identity(os.Getpid())}
	m.PIDNamespace, err = os.Readlink("/proc/self/ns/pid")
	if err != nil {
		return err
	}
	r := report{RunID: id, Label: config.Label, Rerun: "bazel_agent bazel test " + config.Label, RuntimeDir: root, Inputs: map[string]string{}, Tools: map[string]string{}}
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	if err := saveManifest(m); err != nil {
		return err
	}
	fmt.Fprintln(os.Stderr, "Molecule private runtime:", root)
	err = execute(ctx, &m, &r, *mode)
	r.Interrupted = ctx.Err() != nil
	// Recovery reads the latest manifest: lifecycle subprocesses update it.
	latest := m
	_ = readJSON(filepath.Join(root, "manifest.json"), &latest)
	r.Ports = latest.Ports
	r.QEMU = latest.QEMU
	err = errors.Join(err, collect(root, *outputs, latest))
	cleanupErr := stop(root, true)
	r.Interrupted = ctx.Err() != nil
	if r.Interrupted && err == nil {
		err = context.Canceled
		r.Error = err.Error()
	}
	r.Cleanup.Success = cleanupErr == nil
	if cleanupErr != nil {
		r.Cleanup.Error = cleanupErr.Error()
	}
	r.Success = err == nil && cleanupErr == nil && !r.Interrupted
	if err != nil {
		r.Error = err.Error()
	}
	if writeErr := writeJSON(filepath.Join(*outputs, "result.json"), r); writeErr != nil {
		return errors.Join(err, cleanupErr, writeErr)
	}
	return errors.Join(err, cleanupErr)
}

func execute(ctx context.Context, m *manifest, r *report, mode string) error {
	config := &m.Config
	ports := map[int]bool{}
	for _, port := range config.Ports {
		if port < 1 || port > 65535 || ports[port] {
			return errors.New("prerequisite: guest ports must be distinct TCP ports between 1 and 65535")
		}
		ports[port] = true
	}
	if !ports[22] {
		return errors.New("prerequisite: guest ports must include SSH port 22")
	}
	if config.BootTimeout < 1 || config.TCGBootTimeout < 1 {
		return errors.New("prerequisite: boot deadlines must be positive")
	}
	for _, tool := range []string{"ssh", "ssh-keygen"} {
		if _, err := os.Stat("/usr/bin/" + tool); err != nil {
			return fmt.Errorf("prerequisite: install OpenSSH client tool %s through the owning host configuration", tool)
		}
	}
	switch config.Accelerator {
	case "auto":
		f, err := os.OpenFile("/dev/kvm", os.O_RDWR, 0)
		if err == nil {
			f.Close()
			config.Accelerator = "kvm"
		} else {
			config.Accelerator = "tcg"
		}
	case "kvm":
		f, err := os.OpenFile("/dev/kvm", os.O_RDWR, 0)
		if err != nil {
			return fmt.Errorf("prerequisite: requested KVM unavailable: %w", err)
		}
		f.Close()
	case "tcg":
	default:
		return errors.New("prerequisite: accelerator must be auto, kvm, or tcg")
	}
	r.Accelerator = config.Accelerator
	fmt.Fprintln(os.Stderr, "Molecule:", config.Label, "accelerator:", config.Accelerator, "run:", m.RunID)
	rf, err := runfiles.New()
	if err != nil {
		return err
	}
	resolve := func(key string) (string, error) {
		p, err := rf.Rlocation(key)
		if err != nil {
			return "", err
		}
		return filepath.Abs(p)
	}
	config.Image, err = resolve(config.Image)
	if err != nil {
		return fmt.Errorf("prerequisite: declared guest image unavailable: %w", err)
	}
	config.BIOS, err = resolve(config.BIOS)
	if err != nil {
		return err
	}
	for key, value := range config.Tools {
		config.Tools[key], err = resolve(value)
		if err != nil {
			return err
		}
	}
	for _, dir := range []string{"bin", "home", "tmp", "project", "project/collections", "project/inventory"} {
		if err := os.MkdirAll(filepath.Join(m.Root, dir), 0o700); err != nil {
			return err
		}
	}
	for name, path := range config.Tools {
		if err := os.Symlink(path, filepath.Join(m.Root, "bin", name)); err != nil {
			return err
		}
	}
	self, err := os.Executable()
	if err != nil {
		return err
	}
	m.Environment = []string{"PATH=" + filepath.Join(m.Root, "bin") + ":/usr/bin:/bin", "HOME=" + filepath.Join(m.Root, "home"), "TMPDIR=" + filepath.Join(m.Root, "tmp"), "LANG=C.UTF-8", "LC_ALL=C.UTF-8", "PYTHONNOUSERSITE=1", "ANSIBLE_NOCOLOR=1", "NO_COLOR=1", "ANSIBLE_HOST_KEY_CHECKING=False", "ANSIBLE_RETRY_FILES_ENABLED=False", "ANSIBLE_COLLECTIONS_PATH=" + filepath.Join(m.Root, "project/collections"), "ANSIBLE_CONFIG=" + filepath.Join(m.Root, "project/ansible.cfg"), "ANSIBLE_LOCAL_TEMP=" + filepath.Join(m.Root, "tmp/ansible"), "ANSIBLE_SSH_ARGS=-F /dev/null -o ControlMaster=no -o UserKnownHostsFile=/dev/null -o IdentitiesOnly=yes", "MOLECULE_RUNNER=" + self, "MOLECULE_RUNNER_STATE=" + m.Root, "MOLECULE_ACCEPTANCE_MODE=" + mode, "XDG_CACHE_HOME=" + filepath.Join(m.Root, "home/.cache")}
	for _, entry := range rf.Env() {
		m.Environment = append(m.Environment, entry)
	}
	var imports []string
	for _, key := range config.PythonImports {
		path, err := resolve(key)
		if err != nil {
			return err
		}
		imports = append(imports, path)
	}
	m.Environment = append(m.Environment, "PYTHONPATH="+strings.Join(imports, string(os.PathListSeparator)))
	m.Environment = append(m.Environment, "ANSIBLE_PIPELINING=True")
	m.Environment = append(m.Environment, "ANSIBLE_CALLBACK_PLUGINS="+filepath.Join(m.Root, "project/callback_plugins"), "ANSIBLE_CALLBACKS_ENABLED=events", "MOLECULE_EVENTS_FILE="+filepath.Join(m.Root, "events.jsonl"))
	if err := os.WriteFile(filepath.Join(m.Root, "project/ansible.cfg"), []byte("[defaults]\ninterpreter_python=auto_silent\nhost_key_checking=False\nretry_files_enabled=False\n[ssh_connection]\npipelining=True\n"), 0o600); err != nil {
		return err
	}
	for name, path := range map[string]string{"image": config.Image, "bios": config.BIOS} {
		info, err := os.Stat(path)
		if err != nil {
			return fmt.Errorf("prerequisite %s unavailable: %w", name, err)
		}
		if !info.Mode().IsRegular() {
			return fmt.Errorf("prerequisite %s must be a regular declared file", name)
		}
		digest, err := hashFile(path)
		if err != nil {
			return fmt.Errorf("prerequisite %s: %w", name, err)
		}
		r.Inputs[name] = digest
	}
	for name, args := range map[string][]string{"qemu": {"--version"}, "qemu-img": {"--version"}, "molecule": {"--version"}, "ansible": {"--version"}, "seed": {"--version"}} {
		output, err := command(ctx, m, config.Tools[name], args...).CombinedOutput()
		if err != nil {
			return fmt.Errorf("prerequisite: %s version check failed: %w", name, err)
		}
		r.Tools[name] = strings.Split(strings.TrimSpace(string(output)), "\n")[0]
		digest, err := hashFile(config.Tools[name])
		if err != nil {
			return err
		}
		r.Inputs["tool/"+name] = digest
	}
	imageInfo, err := command(ctx, m, config.Tools["qemu-img"], "info", "--output=json", config.Image).Output()
	if err != nil {
		return fmt.Errorf("prerequisite: cannot inspect guest image: %w", err)
	}
	var info struct {
		Format  string `json:"format"`
		Backing string `json:"backing-filename"`
	}
	if err := json.Unmarshal(imageInfo, &info); err != nil {
		return err
	}
	if info.Format != "qcow2" || info.Backing != "" {
		return errors.New("prerequisite: guest image must be a standalone qcow2 image")
	}
	for destination, key := range config.Mappings {
		source, err := resolve(key)
		if err != nil {
			return err
		}
		dest, err := safeJoin(filepath.Join(m.Root, "project"), destination)
		if err != nil {
			return err
		}
		if err := copyInput(source, dest, r.Inputs, "mapping/"+destination); err != nil {
			return err
		}
	}
	for destination, key := range config.Files {
		source, err := resolve(key)
		if err != nil {
			return err
		}
		dest, err := safeJoin(filepath.Join(m.Root, "project/molecule"), destination)
		if err != nil {
			return err
		}
		if err := copyInput(source, dest, r.Inputs, "scenario/"+destination); err != nil {
			return err
		}
	}
	scenarioDir := filepath.Join(m.Root, "project/molecule", config.Scenario)
	for name, key := range config.Lifecycle {
		source, err := resolve(key)
		if err != nil {
			return err
		}
		dest := filepath.Join(scenarioDir, name)
		if name == "events.py" {
			dest = filepath.Join(m.Root, "project/callback_plugins", name)
		}
		if err := copyInput(source, dest, r.Inputs, "lifecycle/"+name); err != nil {
			return err
		}
	}
	if err := os.WriteFile(filepath.Join(m.Root, "project/inventory/hosts.yml"), []byte("all:\n  hosts:\n    localhost:\n      ansible_connection: local\n      ansible_remote_tmp: "+filepath.Join(m.Root, "tmp/local")+"\n"), 0o600); err != nil {
		return err
	}
	if err := writeMolecule(m, scenarioDir); err != nil {
		return err
	}
	if err := saveManifest(*m); err != nil {
		return err
	}
	phases := []string{"create", "prepare", "converge", "idempotence", "verify"}
	if _, err := os.Stat(filepath.Join(scenarioDir, "update.yml")); err == nil {
		phases = append(phases, "side_effect", "converge", "idempotence", "verify")
	}
	phases = append(phases, "destroy")
	for index, phase := range phases {
		if err := ctx.Err(); err != nil {
			return err
		}
		if err := readJSON(filepath.Join(m.Root, "manifest.json"), m); err != nil {
			return err
		}
		m.Phase = phase
		if err := saveManifest(*m); err != nil {
			return err
		}
		fmt.Fprintln(os.Stderr, "Molecule phase:", phase)
		if phase == "destroy" {
			if err := collect(m.Root, filepath.Join(m.Root, "diagnostics"), *m); err != nil {
				return fmt.Errorf("retain guest diagnostics: %w", err)
			}
		}
		err := runPhase(ctx, m, phase, index+1)
		r.Phases = append(r.Phases, phaseResult{Step: index + 1, Name: phase, Success: err == nil})
		if err != nil {
			return fmt.Errorf("%s failed (see phase and assertion artifacts): %w", phase, err)
		}
	}
	if digest, err := hashFile(config.Image); err != nil || digest != r.Inputs["image"] {
		return errors.New("base image changed during execution")
	}
	return nil
}

func writeMolecule(m *manifest, dir string) error {
	sequence := map[string]any{}
	for _, phase := range []string{"create", "prepare", "converge", "idempotence", "verify", "destroy", "side_effect"} {
		sequence[phase+"_sequence"] = []string{phase}
	}
	sequence["test_sequence"] = []string{"create", "prepare", "converge", "idempotence", "verify", "destroy"}
	playbooks := map[string]string{}
	for _, phase := range []string{"create", "prepare", "converge", "verify", "destroy"} {
		playbooks[phase] = phase + ".yml"
	}
	playbooks["side_effect"] = "update.yml"
	args := append([]string{"--inventory=" + filepath.Join(m.Root, "project/inventory")}, m.Config.AnsibleArgs...)
	config := map[string]any{"dependency": map[string]any{"name": "galaxy", "enabled": false}, "ansible": map[string]any{"executor": map[string]any{"backend": "ansible-playbook", "args": map[string]any{"ansible_playbook": args}}, "playbooks": playbooks}, "scenario": sequence, "verifier": map[string]string{"name": "ansible"}}
	return writeJSON(filepath.Join(dir, "molecule.yml"), config)
}

func runPhase(ctx context.Context, m *manifest, phase string, step int) error {
	log, err := os.OpenFile(filepath.Join(m.Root, fmt.Sprintf("%02d-%s.log", step, phase)), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		return err
	}
	defer log.Close()
	cmd := command(ctx, m, m.Config.Tools["molecule"], strings.ReplaceAll(phase, "_", "-"), "--scenario-name", m.Config.Scenario)
	cmd.Env = append(cmd.Env, "MOLECULE_RUNNER_PHASE="+phase)
	cmd.Env = append(cmd.Env, fmt.Sprintf("MOLECULE_RUNNER_STEP=%d", step))
	cmd.Dir = filepath.Join(m.Root, "project")
	cmd.Stdout, cmd.Stderr = log, log
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true, Pdeathsig: syscall.SIGKILL}
	cmd.Cancel = func() error { return syscall.Kill(-cmd.Process.Pid, syscall.SIGTERM) }
	cmd.WaitDelay = 10 * time.Second
	if err := cmd.Start(); err != nil {
		return err
	}
	m.Child = identity(cmd.Process.Pid)
	if err := updateManifest(m.Root, func(current *manifest) { current.Child = m.Child }); err != nil {
		_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		_ = cmd.Wait()
		return err
	}
	err = cmd.Wait()
	if ctx.Err() != nil {
		_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	}
	return err
}

func command(ctx context.Context, m *manifest, name string, args ...string) *exec.Cmd {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Env = m.Environment
	cmd.Dir = m.Root
	return cmd
}

func safeJoin(root, name string) (string, error) {
	if !filepath.IsLocal(name) {
		return "", fmt.Errorf("unsafe declared path %q", name)
	}
	return filepath.Join(root, name), nil
}

func copyInput(source, dest string, hashes map[string]string, key string) error {
	info, err := os.Stat(source)
	if err != nil {
		return err
	}
	if info.IsDir() {
		entries, err := os.ReadDir(source)
		if err != nil {
			return err
		}
		for _, entry := range entries {
			if err := copyInput(filepath.Join(source, entry.Name()), filepath.Join(dest, entry.Name()), hashes, key+"/"+entry.Name()); err != nil {
				return err
			}
		}
		return nil
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("unsupported input type: %s", key)
	}
	if err := os.MkdirAll(filepath.Dir(dest), 0o700); err != nil {
		return err
	}
	in, err := os.Open(source)
	if err != nil {
		return err
	}
	defer in.Close()
	mode := os.FileMode(0o600)
	if info.Mode()&0o111 != 0 {
		mode = 0o700
	}
	out, err := os.OpenFile(dest, os.O_CREATE|os.O_EXCL|os.O_WRONLY, mode)
	if err != nil {
		return err
	}
	h := sha256.New()
	_, copyErr := io.Copy(io.MultiWriter(out, h), in)
	closeErr := out.Close()
	hashes[key] = hex.EncodeToString(h.Sum(nil))
	return errors.Join(copyErr, closeErr)
}

func hashFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func readJSON(path string, v any) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, v)
}

func writeJSON(path string, v any) error {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(path+".new", append(b, '\n'), 0o600); err != nil {
		return err
	}
	return os.Rename(path+".new", path)
}
func saveManifest(m manifest) error { return writeJSON(filepath.Join(m.Root, "manifest.json"), m) }

func updateManifest(root string, update func(*manifest)) error {
	lock, err := os.OpenFile(filepath.Join(root, "manifest.lock"), os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		return err
	}
	defer lock.Close()
	if err := syscall.Flock(int(lock.Fd()), syscall.LOCK_EX); err != nil {
		return err
	}
	defer syscall.Flock(int(lock.Fd()), syscall.LOCK_UN)
	m, err := loadOwned(root)
	if err != nil {
		return err
	}
	update(&m)
	return saveManifest(m)
}

// Keep only task/play headings and recap counters. Ansible result bodies and
// raw serial output can contain fixture credentials and stay in private state.
func collect(root, outputs string, m manifest) error {
	if err := os.MkdirAll(outputs, 0o700); err != nil {
		return err
	}
	var retainedErrors []error
	write := func(name string, data []byte) {
		if err := os.WriteFile(filepath.Join(outputs, name), data, 0o600); err != nil {
			retainedErrors = append(retainedErrors, err)
		}
	}
	coverage := filepath.Join(root, "project/molecule", m.Config.Scenario, "coverage.json")
	if data, err := os.ReadFile(coverage); err == nil {
		write("coverage.json", data)
	}
	if data, err := os.ReadFile(filepath.Join(root, "events.jsonl")); err == nil {
		write("events.jsonl", data)
	}
	if outputs != filepath.Join(root, "diagnostics") {
		entries, _ := os.ReadDir(filepath.Join(root, "diagnostics"))
		for _, entry := range entries {
			if data, err := os.ReadFile(filepath.Join(root, "diagnostics", entry.Name())); err == nil {
				write(entry.Name(), data)
			} else {
				retainedErrors = append(retainedErrors, err)
			}
		}
	}
	entries, _ := os.ReadDir(root)
	for _, entry := range entries {
		if !strings.HasSuffix(entry.Name(), ".log") || entry.Name() == "serial.log" {
			continue
		}
		data, err := os.ReadFile(filepath.Join(root, entry.Name()))
		if err != nil {
			retainedErrors = append(retainedErrors, err)
			continue
		}
		var selected []string
		for _, line := range strings.Split(string(data), "\n") {
			if strings.HasPrefix(line, "PLAY [") || strings.HasPrefix(line, "TASK [") || strings.HasPrefix(line, "RUNNING HANDLER [") || strings.HasPrefix(line, "PLAY RECAP") || (strings.HasPrefix(line, "molecule") && strings.Contains(line, "changed=")) {
				selected = append(selected, line)
			}
		}
		write(entry.Name(), []byte(strings.Join(selected, "\n")+"\n"))
	}
	if alive(m.QEMU) {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		cmd := sshCommand(ctx, &m, "rpm -qa --qf '%{NAME}-%{VERSION}-%{RELEASE}.%{ARCH}\\n'")
		if data, err := cmd.Output(); err == nil {
			lines := strings.Split(strings.TrimSpace(string(data)), "\n")
			sort.Strings(lines)
			data, err := json.MarshalIndent(lines, "", "  ")
			if err != nil {
				retainedErrors = append(retainedErrors, err)
			} else {
				write("packages.json", append(data, '\n'))
			}
		}
		if len(m.Config.Services) > 0 {
			cmd := sshCommand(ctx, &m, "systemctl show --property=Id,ActiveState,SubState,Result,ExecMainCode,ExecMainStatus,NRestarts,ExecMainStartTimestampMonotonic "+strings.Join(m.Config.Services, " "))
			if data, err := cmd.Output(); err == nil {
				write("services.txt", data)
			}
		}
	}
	return errors.Join(retainedErrors...)
}
