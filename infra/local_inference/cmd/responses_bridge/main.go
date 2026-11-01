package main

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"

	"git.alwaldend.com/alwaldend/src/infra/local_inference/internal/responses"
	"github.com/spf13/cobra"
)

// main reports command failures without logging request contents.
func main() {
	if err := command().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// command configures the two loopback endpoints and starts the HTTP process.
func command() *cobra.Command {
	listen := "127.0.0.1:8080"
	upstream := "http://127.0.0.1:8081"
	cmd := &cobra.Command{Use: "responses_bridge", Args: cobra.NoArgs, SilenceUsage: true, SilenceErrors: true}
	cmd.Flags().StringVar(&listen, "listen", listen, "loopback listen address")
	cmd.Flags().StringVar(&upstream, "upstream", upstream, "loopback llama.cpp URL")
	cmd.RunE = func(cmd *cobra.Command, _ []string) error {
		target, err := url.Parse(upstream)
		if err != nil {
			return fmt.Errorf("parse inference URL: %w", err)
		}
		if target.Scheme != "http" || target.User != nil || target.Path != "" || target.RawQuery != "" || target.Fragment != "" {
			return fmt.Errorf("inference URL must be a loopback HTTP origin")
		}
		for _, address := range []string{listen, target.Host} {
			host, _, err := net.SplitHostPort(address)
			if err != nil {
				return fmt.Errorf("parse loopback endpoint %q: %w", address, err)
			}
			if host != "127.0.0.1" {
				return fmt.Errorf("endpoint must use 127.0.0.1")
			}
		}
		listener, err := net.Listen("tcp", listen)
		if err != nil {
			return fmt.Errorf("listen for Codex: %w", err)
		}
		defer listener.Close()
		fmt.Fprintln(cmd.OutOrStdout(), listener.Addr().String())
		server := http.Server{Handler: responses.Handler(target), ReadHeaderTimeout: 10 * time.Second, IdleTimeout: 120 * time.Second}
		if err := server.Serve(listener); err != nil {
			return fmt.Errorf("serve Responses bridge: %w", err)
		}
		return nil
	}
	return cmd
}
