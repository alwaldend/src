package evaluation

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	vastclient "git.alwaldend.com/alwaldend/src/internal/vastai_sdk/pkg/client"
)

// prepareStartup compiles the canonical template without providers, plans or state.
func prepareStartup(ctx context.Context, options Options, profiles registry, model model, qwenGIF bool, directory string, login credentials) (result string, resultErr error) {
	if options.StartupFile != "" {
		data, err := os.ReadFile(options.StartupFile)
		if err != nil {
			return "", fmt.Errorf("read explicit startup recipe: %w", err)
		}
		command, err := vastclient.StartupCommand(string(data), true)
		if err != nil {
			return "", fmt.Errorf("encode explicit startup recipe: %w", err)
		}
		return command, nil
	}
	password := login.password
	username := login.username
	if !regexp.MustCompile(`^[A-Za-z0-9_-]{1,32}$`).MatchString(username) || !regexp.MustCompile(`^[A-Za-z0-9]{32,64}$`).MatchString(password) {
		return "", errors.New("evaluation requires the component's Vault-injected HTTPS username/password")
	}
	var runtime map[string]any
	if err := readJSON(options.RuntimeDefaults, &runtime); err != nil {
		return "", fmt.Errorf("read shared runtime declaration: %w", err)
	}
	directories, ok := runtime["model_directories"].(map[string]any)
	if !ok {
		return "", errors.New("runtime declaration omits model directories")
	}
	if model.CheckpointDirectory != "" {
		directories["checkpoint"] = model.CheckpointDirectory
	}
	downloads := []map[string]string{}
	for _, filename := range model.Downloads {
		entry, found := profiles.Downloads[filename]
		destination, valid := directories[entry.Kind].(string)
		u, err := url.Parse(entry.URL)
		if !found || !valid || err != nil || u.Scheme != "https" || u.Hostname() == "" || u.User != nil || u.RawQuery != "" || filepath.Base(filename) != filename || filename == "." || filename == ".." || strings.ContainsAny(filename, "'\r\n\\") || !strings.HasPrefix(destination, "models/") || strings.Contains(destination, "..") {
			return "", fmt.Errorf("model asset %q has an invalid publisher or destination", filename)
		}
		downloads = append(downloads, map[string]string{"url": entry.URL, "dest": fmt.Sprint(runtime["model_root"]) + "/" + destination + "/" + filename})
	}
	repositories := []map[string]string{}
	var enabled []string
	names := append([]string{}, model.Nodes...)
	sort.Strings(names)
	for _, name := range names {
		node, found := profiles.Nodes[name]
		if !found || !regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_-]*$`).MatchString(name) || !regexp.MustCompile(`^https://github[.]com/[A-Za-z0-9_-]+/[A-Za-z0-9_.-]+$`).MatchString(node.URL) || !regexp.MustCompile(`^[0-9a-f]{40}$`).MatchString(node.Ref) {
			return "", fmt.Errorf("node pack %q is not a valid approved registry entry", name)
		}
		repositories = append(repositories, map[string]string{"url": node.URL, "ref": node.Ref, "dest": "/workspace/ComfyUI/custom_nodes/" + name})
		enabled = append(enabled, name)
	}
	var files []map[string]string
	// firstPartySource selects a module without coupling model-specific code to other bundles.
	type firstPartySource struct {
		// name is the first-party module directory installed by the provisioner.
		name string
		// path identifies its existing reviewed source.
		path string
	}
	sources := []firstPartySource{{"chained_img2img_gif", options.GIFSource}, {"prompt_sequence", options.SequenceSource}}
	if qwenGIF {
		sources = append(sources, firstPartySource{"qwen_chained_img2img_gif", options.QwenGIFSource})
	}
	for _, source := range sources {
		data, err := os.ReadFile(source.path)
		if err != nil {
			return "", fmt.Errorf("read first-party node %q: %w", source.name, err)
		}
		files = append(files, map[string]string{"path": "/workspace/ComfyUI/custom_nodes/" + source.name + "/__init__.py", "content": string(data), "permissions": "0600"})
		enabled = append(enabled, source.name)
	}
	packages := []map[string]any{}
	if model.Requirements != "" {
		data, err := os.ReadFile(options.Requirements)
		if err != nil {
			return "", fmt.Errorf("read approved node requirements: %w", err)
		}
		files = append(files, map[string]string{"path": "/etc/imagegen-python-requirements.txt", "content": string(data), "permissions": "0600"})
		packages = append(packages, map[string]any{"venv": "/venv/main", "requirements": []string{"/etc/imagegen-python-requirements.txt"}, "args": "--require-hashes --no-deps"})
	}
	sort.Strings(enabled)
	runtime["arguments"] = fmt.Sprint(runtime["arguments"]) + " --whitelist-custom-nodes " + strings.Join(enabled, " ")
	manifest, err := json.Marshal(map[string]any{"version": 1, "downloads": downloads, "git_repos": repositories, "write_files": files, "pip_packages": packages})
	if err != nil {
		return "", fmt.Errorf("encode native provisioner manifest: %w", err)
	}
	parameters, err := json.Marshal(map[string]any{"model_manifest": string(manifest), "hf_token_b64": base64.StdEncoding.EncodeToString([]byte(login.huggingface)), "runtime": runtime, "runtime_files": map[string]string{}, "exposed": true, "web_username": username, "web_password": password})
	if err != nil {
		return "", fmt.Errorf("encode template parameters: %w", err)
	}
	if err := os.MkdirAll(directory, 0o700); err != nil {
		return "", fmt.Errorf("create empty renderer directory: %w", err)
	}
	defer func() {
		if err := os.RemoveAll(directory); err != nil {
			resultErr = errors.Join(resultErr, fmt.Errorf("remove renderer directory: %w", err))
		}
	}()
	path, err := filepath.Abs(options.StartupTemplate)
	if err != nil {
		return "", fmt.Errorf("resolve startup template: %w", err)
	}
	quotedPath, _ := json.Marshal(path)
	expression := "base64encode(templatefile(" + string(quotedPath) + ", jsondecode(base64decode(\"" + base64.StdEncoding.EncodeToString(parameters) + "\"))))\n"
	child := processCommand(ctx, options.Terraform, "-chdir="+directory, "console", "-no-color")
	child.Env = workloadEnvironment()
	child.Stdin = strings.NewReader(expression)
	var output bytes.Buffer
	child.Stdout = &output
	child.Stderr = &bytes.Buffer{}
	if err := child.Run(); err != nil {
		return "", fmt.Errorf("render startup template (secret-bearing diagnostics suppressed): %w", err)
	}
	var encoded string
	if err := json.Unmarshal(bytes.TrimSpace(output.Bytes()), &encoded); err != nil {
		return "", fmt.Errorf("decode template renderer string: %w", err)
	}
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", fmt.Errorf("decode rendered startup: %w", err)
	}
	command, err := vastclient.StartupCommand(string(decoded), true)
	if err != nil {
		return "", fmt.Errorf("encode compiled startup: %w", err)
	}
	return command, nil
}
