// xo_config reconciles the XO OIDC plugin through its JSON-RPC API.
package main

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/url"
	"os"
	"sort"
	"time"

	"golang.org/x/net/websocket"
)

type rpc struct {
	conn *websocket.Conn
	id   int
}

// Errors intentionally exclude server messages, payloads and connection URLs:
// these may echo credentials or plugin configuration.
func (r *rpc) call(method string, params any, result any) error {
	r.id++
	if err := r.conn.SetDeadline(time.Now().Add(45 * time.Second)); err != nil {
		return errors.New("setting RPC deadline failed")
	}
	if err := websocket.JSON.Send(r.conn, map[string]any{"jsonrpc": "2.0", "id": r.id, "method": method, "params": params}); err != nil {
		return fmt.Errorf("%s: sending request failed", method)
	}
	for n := 0; n < 1000; n++ {
		var response struct {
			ID     int             `json:"id"`
			Result json.RawMessage `json:"result"`
			Error  json.RawMessage `json:"error"`
		}
		if err := websocket.JSON.Receive(r.conn, &response); err != nil {
			return fmt.Errorf("%s: receiving response failed", method)
		}
		if response.ID != r.id {
			continue
		}
		if len(response.Error) > 0 && string(response.Error) != "null" {
			return fmt.Errorf("%s: server rejected request (details suppressed)", method)
		}
		if result != nil {
			if err := json.Unmarshal(response.Result, result); err != nil {
				return fmt.Errorf("%s: unexpected response format", method)
			}
		}
		return nil
	}
	return fmt.Errorf("%s: response notification limit reached", method)
}

func connect(endpoint, token string) (*rpc, error) {
	u, err := url.Parse(endpoint)
	if err != nil || u.Scheme != "wss" || u.Host == "" || u.User != nil || u.RawQuery != "" || u.Fragment != "" {
		return nil, errors.New("XOA_URL must be a wss URL without credentials, query or fragment")
	}
	if token == "" {
		return nil, errors.New("XOA_TOKEN is required")
	}
	if u.Path == "" || u.Path == "/" {
		u.Path = "/api/"
	}
	origin := *u
	origin.Scheme = "https"
	origin.Path = "/"
	cfg, err := websocket.NewConfig(u.String(), origin.String())
	if err != nil {
		return nil, errors.New("invalid websocket configuration")
	}
	cfg.Dialer = &net.Dialer{Timeout: 15 * time.Second}
	if os.Getenv("XOA_INSECURE") == "true" {
		cfg.TlsConfig = &tls.Config{InsecureSkipVerify: true}
	} // Explicit bootstrap opt-in, matching the provider.
	conn, err := websocket.DialConfig(cfg)
	if err != nil {
		return nil, errors.New("connecting to XO failed")
	}
	conn.MaxPayloadBytes = 32 << 20
	r := &rpc{conn: conn}
	if err := r.call("session.signInWithToken", map[string]string{"token": token}, nil); err != nil {
		conn.Close()
		return nil, err
	}
	return r, nil
}

type plugin struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Version  string `json:"version"`
	Loaded   bool   `json:"loaded"`
	Autoload bool   `json:"autoload"`
}

func plugins(r *rpc) ([]plugin, error) {
	var result []plugin
	err := r.call("plugin.get", map[string]any{}, &result)
	return result, err
}

func oidcPlugin(r *rpc) (plugin, error) {
	all, err := plugins(r)
	if err != nil {
		return plugin{}, err
	}
	for _, p := range all {
		if p.ID == "auth-oidc" {
			return p, nil
		}
	}
	return plugin{}, errors.New("auth-oidc plugin is not installed")
}

func applyOIDC(r *rpc, configuration string, out io.Writer) error {
	var conf map[string]any
	if json.Unmarshal([]byte(configuration), &conf) != nil || conf == nil {
		return errors.New("XO_OIDC_CONFIGURATION must be a JSON object")
	}
	for _, key := range []string{"clientID", "clientSecret", "discoveryURL"} {
		if value, ok := conf[key].(string); !ok || value == "" {
			return fmt.Errorf("OIDC configuration requires %s", key)
		}
	}
	p, err := oidcPlugin(r)
	if err != nil {
		return err
	}
	if err = r.call("plugin.configure", map[string]any{"id": p.ID, "configuration": conf}, nil); err != nil {
		return err
	}
	if !p.Autoload {
		if err = r.call("plugin.enableAutoload", map[string]string{"id": p.ID}, nil); err != nil {
			return err
		}
	}
	// Configuring a loaded plugin reloads it internally. Loading it again fails.
	if !p.Loaded {
		if err = r.call("plugin.load", map[string]string{"id": p.ID}, nil); err != nil {
			return err
		}
	}
	p, err = oidcPlugin(r)
	if err != nil {
		return err
	}
	if !p.Loaded || !p.Autoload {
		return errors.New("OIDC plugin postcondition failed: loaded and autoload must be true")
	}
	if raw := os.Getenv("XO_OIDC_AUTHORIZATION"); raw != "" {
		if err := applyAuthorization(r, raw, configuration); err != nil {
			return err
		}
	}
	return json.NewEncoder(out).Encode(p)
}

func inspect(r *rpc, out io.Writer) error {
	var objects map[string]map[string]any
	if err := r.call("xo.getAllObjects", map[string]any{}, &objects); err != nil {
		return err
	}
	allowed := map[string]bool{"pool": true, "host": true, "SR": true, "network": true, "PIF": true, "VM-template": true, "VM": true, "VBD": true, "VDI": true}
	var selected []map[string]any
	for _, obj := range objects {
		kind, _ := obj["type"].(string)
		if !allowed[kind] {
			continue
		}
		clean := map[string]any{}
		for _, key := range []string{"id", "type", "name_label", "$pool", "$container", "size", "physical_usage", "SR_type", "writable", "shared", "default_SR", "memory", "version", "software_version", "patches", "PIFs", "VBDs", "$VBDs", "VDI", "VM", "$SR", "is_cd_drive", "bootable", "position", "power_state", "$network", "network", "device", "vlan", "VLAN", "ip", "netmask", "gateway", "management", "currently_attached", "bridge"} {
			if v, ok := obj[key]; ok {
				clean[key] = v
			}
		}
		selected = append(selected, clean)
	}
	sort.Slice(selected, func(i, j int) bool { return fmt.Sprint(selected[i]["id"]) < fmt.Sprint(selected[j]["id"]) })
	ps, err := plugins(r)
	if err != nil {
		return err
	}
	var groups []struct {
		ID       string `json:"id"`
		Name     string `json:"name"`
		Provider string `json:"provider,omitempty"`
	}
	if err := r.call("group.getAll", map[string]any{}, &groups); err != nil {
		return err
	}
	var servers []struct {
		ID      string `json:"id"`
		Host    string `json:"host"`
		Enabled bool   `json:"enabled"`
		Status  string `json:"status"`
	}
	if err := r.call("server.getAll", map[string]any{}, &servers); err != nil {
		return err
	}
	return json.NewEncoder(out).Encode(map[string]any{"objects": selected, "plugins": ps, "groups": groups, "servers": servers})
}

func run(args []string, out io.Writer) error {
	if len(args) == 0 || (len(args) != 1 && !(len(args) == 2 && args[0] == "import-template")) || (args[0] != "oidc-identities" && args[0] != "inspect" && args[0] != "apply-oidc" && args[0] != "register-host" && args[0] != "import-template") {
		return errors.New("usage: xo_config oidc-identities|inspect|apply-oidc|register-host|import-template [image-path]")
	}
	if args[0] == "oidc-identities" {
		return externalOIDCIdentities(os.Stdin, out)
	}
	r, err := connect(os.Getenv("XOA_URL"), os.Getenv("XOA_TOKEN"))
	if err != nil {
		return err
	}
	defer r.conn.Close()
	if args[0] == "inspect" {
		return inspect(r, out)
	}
	if args[0] == "register-host" {
		return registerHostEnv(r, out)
	}
	if args[0] == "import-template" {
		imagePath := os.Getenv("XO_TEMPLATE_IMAGE")
		if len(args) == 2 {
			imagePath = args[1]
		}
		return importTemplateEnv(r, imagePath, out)
	}
	return applyOIDC(r, os.Getenv("XO_OIDC_CONFIGURATION"), out)
}

func main() {
	if err := run(os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
