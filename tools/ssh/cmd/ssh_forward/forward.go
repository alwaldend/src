package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
)

const (
	// startupTimeout bounds SSH authentication and destination verification.
	startupTimeout = 10 * time.Second
	// connectionLimit bounds concurrent forwarded connections for a listener.
	connectionLimit = 64
)

// forwardConfig describes one authenticated SSH forwarding request.
type forwardConfig struct {
	// host is the SSH server address, including its port.
	host string
	// user is the remote SSH principal.
	user string
	// target is the TCP address as resolved by the SSH server.
	target string
	// ca is the trusted SSH host CA public key in authorized-keys format.
	ca string
	// keyPath points to the borrowed SSH client private key.
	keyPath string
	// certPath points to the borrowed signed SSH client certificate.
	certPath string
}

// forward owns one authenticated SSH session and its loopback listener.
type forward struct {
	// client carries forwarded TCP channels.
	client *ssh.Client
	// listener accepts connections on an invocation-local port.
	listener net.Listener
	// cancel closes accepted local connections when the tunnel ends.
	cancel context.CancelFunc
}

// readCredential bounds the amount read from a caller-supplied SSH key or certificate.
func readCredential(path string) (_ []byte, result error) {
	const maxCredentialBytes = 1 << 20
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open SSH credential %q: %w", path, err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			result = errors.Join(result, fmt.Errorf("close SSH credential %q: %w", path, err))
		}
	}()
	data, err := io.ReadAll(io.LimitReader(file, maxCredentialBytes+1))
	if err != nil {
		return nil, fmt.Errorf("read SSH credential %q: %w", path, err)
	}
	if len(data) > maxCredentialBytes {
		return nil, fmt.Errorf("SSH credential %q exceeds %d bytes", path, maxCredentialBytes)
	}
	return data, nil
}

// clientConfig loads the supplied keypair and pins the SSH host certificate CA.
func (c forwardConfig) clientConfig() (*ssh.ClientConfig, error) {
	if c.user == "" || c.keyPath == "" || c.certPath == "" || c.ca == "" {
		return nil, errors.New("SSH user, client key and certificate, and host CA are required")
	}
	ca, _, _, rest, err := ssh.ParseAuthorizedKey([]byte(c.ca))
	if err != nil {
		return nil, fmt.Errorf("parse SSH host CA: %w", err)
	}
	if len(bytes.TrimSpace(rest)) != 0 {
		return nil, errors.New("expected exactly one SSH host CA")
	}
	key, err := readCredential(c.keyPath)
	if err != nil {
		return nil, fmt.Errorf("read SSH client key: %w", err)
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, fmt.Errorf("parse SSH client key: %w", err)
	}
	data, err := readCredential(c.certPath)
	if err != nil {
		return nil, fmt.Errorf("read SSH client certificate: %w", err)
	}
	publicKey, _, _, _, err := ssh.ParseAuthorizedKey(data)
	if err != nil {
		return nil, fmt.Errorf("parse SSH client certificate: %w", err)
	}
	cert, ok := publicKey.(*ssh.Certificate)
	if !ok || cert.CertType != ssh.UserCert {
		return nil, errors.New("SSH client requires a user certificate")
	}
	identity, err := ssh.NewCertSigner(cert, signer)
	if err != nil {
		return nil, fmt.Errorf("match SSH client certificate to its key: %w", err)
	}
	checker := &ssh.CertChecker{
		IsHostAuthority: func(key ssh.PublicKey, _ string) bool {
			return bytes.Equal(key.Marshal(), ca.Marshal())
		},
	}
	return &ssh.ClientConfig{
		User: c.user, Auth: []ssh.AuthMethod{ssh.PublicKeys(identity)},
		HostKeyCallback: func(host string, addr net.Addr, key ssh.PublicKey) error {
			certificate, ok := key.(*ssh.Certificate)
			if !ok || len(certificate.ValidPrincipals) == 0 {
				return errors.New("SSH server must present a host certificate with an explicit principal")
			}
			if err := checker.CheckHostKey(host, addr, key); err != nil {
				return fmt.Errorf("verify SSH host certificate: %w", err)
			}
			return nil
		},
	}, nil
}

// openForward authenticates SSH, probes the destination, and publishes a listener.
func openForward(ctx context.Context, cfg forwardConfig) (_ *forward, failure error) {
	if _, _, err := net.SplitHostPort(cfg.host); err != nil {
		return nil, fmt.Errorf("parse SSH server address: %w", err)
	}
	if _, _, err := net.SplitHostPort(cfg.target); err != nil {
		return nil, fmt.Errorf("parse SSH forward destination: %w", err)
	}
	settings, err := cfg.clientConfig()
	if err != nil {
		return nil, fmt.Errorf("prepare SSH authentication: %w", err)
	}
	startup, cancelStartup := context.WithTimeout(ctx, startupTimeout)
	defer cancelStartup()
	raw, err := (&net.Dialer{}).DialContext(startup, "tcp", cfg.host)
	if err != nil {
		return nil, fmt.Errorf("dial SSH server %q: %w", cfg.host, err)
	}
	stopHandshake := context.AfterFunc(startup, func() { _ = raw.Close() })
	defer stopHandshake()
	conn, channels, requests, err := ssh.NewClientConn(raw, cfg.host, settings)
	if err != nil {
		_ = raw.Close()
		return nil, fmt.Errorf("authenticate SSH server %q: %w", cfg.host, err)
	}
	client := ssh.NewClient(conn, channels, requests)
	defer func() {
		if failure != nil {
			_ = client.Close()
		}
	}()
	destination, err := client.DialContext(startup, "tcp", cfg.target)
	if err != nil {
		return nil, fmt.Errorf("reach SSH forward destination %q: %w", cfg.target, err)
	}
	if err := destination.Close(); err != nil {
		return nil, fmt.Errorf("close SSH destination probe: %w", err)
	}
	stopHandshake()
	if err := startup.Err(); err != nil {
		return nil, fmt.Errorf("SSH forward startup expired: %w", err)
	}
	listener, err := (&net.ListenConfig{}).Listen(startup, "tcp4", "127.0.0.1:0")
	if err != nil {
		return nil, fmt.Errorf("allocate SSH loopback listener: %w", err)
	}
	tunnelCtx, cancel := context.WithCancel(context.Background())
	tunnel := &forward{client: client, listener: listener, cancel: cancel}
	go tunnel.serve(tunnelCtx, cfg.target)
	go func() {
		_ = client.Wait()
		cancel()
		_ = listener.Close()
	}()
	return tunnel, nil
}

// stop closes both ends of the tunnel, including idle accepted connections.
func (f *forward) stop(context.Context) error {
	f.cancel()
	var result error
	if err := f.listener.Close(); err != nil && !errors.Is(err, net.ErrClosed) {
		result = fmt.Errorf("close SSH forward listener: %w", err)
	}
	if err := f.client.Close(); err != nil && !errors.Is(err, net.ErrClosed) && !errors.Is(err, io.EOF) {
		result = errors.Join(result, fmt.Errorf("close SSH client connection: %w", err))
	}
	return result
}

// serve accepts bounded concurrent connections until the listener closes.
func (f *forward) serve(ctx context.Context, target string) {
	inFlight := make(chan struct{}, connectionLimit)
	for {
		local, err := f.listener.Accept()
		if err != nil {
			return
		}
		select {
		case inFlight <- struct{}{}:
			go func() {
				defer func() { <-inFlight }()
				f.forwardConnection(ctx, local, target)
			}()
		default:
			_ = local.Close()
		}
	}
}

// forwardConnection connects a local socket to one SSH TCP channel.
func (f *forward) forwardConnection(ctx context.Context, local net.Conn, target string) {
	defer local.Close()
	stopLocal := context.AfterFunc(ctx, func() { _ = local.Close() })
	defer stopLocal()
	remote, err := f.client.DialContext(ctx, "tcp", target)
	if err != nil {
		return
	}
	defer remote.Close()
	sent := make(chan struct{})
	go func() {
		copyStream(remote, local)
		close(sent)
	}()
	copyStream(local, remote)
	<-sent
}

// copyStream propagates EOF while allowing traffic in the opposite direction.
func copyStream(destination, source net.Conn) {
	_, _ = io.Copy(destination, source)
	if writer, ok := destination.(interface{ CloseWrite() error }); ok {
		_ = writer.CloseWrite()
	} else {
		_ = destination.Close()
	}
}
