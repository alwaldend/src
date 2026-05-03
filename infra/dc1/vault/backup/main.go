package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"git.alwaldend.com/alwaldend/src/tools/al/pkg/al"
)

var rcloneFlag = flag.String("rclone", "", "Rclone binary")
var vaultFlag = flag.String("vault", "", "Vault binary")

func main() {
	now := time.Now()
	flag.Parse()
	tmpDir := al.Must(os.MkdirTemp("", "infra_dc1_vault_backup_*"))
	defer os.RemoveAll(tmpDir)
	vaultName := fmt.Sprintf("vault-backup-%s.snap", now.Format(time.RFC3339))
	vaultPath := filepath.Join(tmpDir, vaultName)
	bucket := os.Getenv("AL_VAULT_BACKUP_BUCKET")
	if bucket == "" {
		log.Fatalln("missing bucket name")
	}
	al.Check(al.RunCommand(al.CommandArgs{Name: *rcloneFlag, Args: []string{"config", "dump"}}))
	al.Check(al.RunCommand(al.CommandArgs{Name: *vaultFlag, Args: []string{"operator", "raft", "snapshot", "save", vaultPath}}))
	al.Check(al.RunCommand(al.CommandArgs{Name: *rcloneFlag, Args: []string{"copyto", "-v", "--check-first", "--s3-no-check-bucket", vaultPath, fmt.Sprintf("remote:%s/vault/%s", bucket, vaultName)}}))
}
