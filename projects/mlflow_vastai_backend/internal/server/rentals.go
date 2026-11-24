package server

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	vastai "git.alwaldend.com/alwaldend/src/internal/vastai_sdk/pkg/client"
	"git.alwaldend.com/alwaldend/src/projects/mlflow_vastai_backend/api/jobpb"
	"google.golang.org/protobuf/proto"
)

// errNoRental indicates that an allowed selector has no eligible rental.
var errNoRental = errors.New("no eligible rental")

// errRentalLost identifies a rental confirmed absent by Vast.ai.
var errRentalLost = errors.New("assigned rental no longer exists")

// sshHost limits instance addresses to hostnames and IP literals.
var sshHost = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9.:-]*$`)

// rentals discovers operator-managed instances and invokes their worker clients.
type rentals struct {
	// api supplies read-only discovery.
	api *vastai.Client
	// allowed lists IDs or exact labels in order of preference.
	allowed []string
	// user is the SSH login name.
	user string
}

// selectRental returns a running rental selected from the configured pool.
func (rentals rentals) selectRental(ctx context.Context, worker string) (*vastai.Instance, error) {
	var listed []vastai.InstanceIdentity
	for _, selector := range rentals.allowed {
		if worker != "" && selector != worker {
			continue
		}
		id, err := strconv.ParseInt(selector, 10, 64)
		numeric := err == nil && id > 0
		if !numeric {
			if listed == nil {
				listed, err = rentals.api.Instances(ctx)
				if err != nil {
					return nil, fmt.Errorf("discover rental labels: %w", err)
				}
			}
			id = 0
			for _, instance := range listed {
				if instance.Label == selector {
					if id != 0 {
						return nil, fmt.Errorf("worker label %q is ambiguous: %w", selector, errNoRental)
					}
					id = instance.ID
				}
			}
		}
		if id == 0 {
			continue
		}
		instance, err := rentals.resolve(ctx, id)
		if errors.Is(err, errRentalLost) || errors.Is(err, errNoRental) {
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("resolve eligible rental: %w", err)
		}
		if !numeric && instance.Label != selector {
			continue
		}
		return instance, nil
	}
	return nil, errNoRental
}

// resolve refreshes the assigned rental and verifies its SSH endpoint.
func (rentals rentals) resolve(ctx context.Context, id int64) (*vastai.Instance, error) {
	instance, err := rentals.api.Get(ctx, id)
	if errors.Is(err, vastai.ErrNotFound) {
		return nil, fmt.Errorf("resolve rental %d: %w", id, errRentalLost)
	}
	if err != nil {
		return nil, fmt.Errorf("resolve rental %d: %w", id, err)
	}
	if instance.Status == nil || *instance.Status != "running" ||
		instance.SSHHost == nil || !sshHost.MatchString(*instance.SSHHost) ||
		instance.SSHPort == nil || *instance.SSHPort < 1 || *instance.SSHPort > 65535 {
		return nil, fmt.Errorf("rental %d has no running SSH endpoint: %w", id, errNoRental)
	}
	return instance, nil
}

// call exchanges one protobuf message with a worker over a bounded SSH session.
func (rentals rentals) call(ctx context.Context, instance *vastai.Instance, operation string, message proto.Message) (*jobpb.Job, error) {
	payload, err := proto.Marshal(message)
	if err != nil {
		return nil, fmt.Errorf("encode worker request: %w", err)
	}
	ctx, cancel := context.WithTimeout(ctx, 12*time.Second)
	defer cancel()
	command := exec.CommandContext(ctx, "ssh", "-T", "-o", "BatchMode=yes",
		"-o", "StrictHostKeyChecking=yes", "-o", "ClearAllForwardings=yes",
		"-o", "ForwardAgent=no", "-o", "ConnectTimeout=5", "-l", rentals.user,
		"-p", strconv.FormatInt(*instance.SSHPort, 10), "--", *instance.SSHHost,
		"mlflow-vastai worker "+operation)
	command.Stdin = bytes.NewReader(payload)
	command.WaitDelay = time.Second
	for _, variable := range command.Environ() {
		if !strings.HasPrefix(variable, "VAST_API_KEY=") && !strings.HasPrefix(variable, "MLFLOW_TRACKING_TOKEN=") {
			command.Env = append(command.Env, variable)
		}
	}
	output, err := command.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("open SSH response: %w", err)
	}
	if err := command.Start(); err != nil {
		return nil, fmt.Errorf("start worker SSH: %w", err)
	}
	const maximumReply = 64 << 10
	reply, readErr := io.ReadAll(io.LimitReader(output, maximumReply+1))
	if readErr != nil || len(reply) > maximumReply {
		cancel()
	}
	waitErr := command.Wait()
	if readErr != nil || waitErr != nil {
		return nil, fmt.Errorf("receive worker reply: %w", errors.Join(readErr, waitErr))
	}
	if len(reply) > maximumReply {
		return nil, errors.New("SSH reply exceeded the job size limit")
	}
	job := &jobpb.Job{}
	if err := proto.Unmarshal(reply, job); err != nil {
		return nil, fmt.Errorf("decode worker reply: %w", err)
	}
	return job, nil
}
