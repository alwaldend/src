package ingester_test

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bazelbuild/rules_go/go/runfiles"
)

var (
	adbPath      = flag.String("adb", "", "Path to the adb executable")
	packageName  = flag.String("package_name", "", "Android package name to install and launch")
	activityName = flag.String("activity_name", "", "Launcher activity to start")
)

func adbCommand(t *testing.T, args ...string) *exec.Cmd {
	t.Helper()
	return exec.Command(*adbPath, args...)
}

func adbOutput(t *testing.T, args ...string) string {
	t.Helper()
	output, err := adbCommand(t, args...).CombinedOutput()
	if err != nil {
		t.Fatalf("adb %s: %v\n%s", strings.Join(args, " "), err, output)
	}
	return string(output)
}

func TestIngesterSmoke(t *testing.T) {
	apkPath, err := runfiles.Rlocation("_main/projects/activitywatch_ingester_android/main/java/ingester_binary.apk")
	if err != nil {
		t.Fatalf("resolve ingester APK: %v", err)
	}
	apkFile, err := os.Open(apkPath)
	if err != nil {
		t.Fatalf("open ingester APK: %v", err)
	}
	defer apkFile.Close()

	adbOutput(t, "devices")
	adbOutput(t, "install", "-r", apkPath)

	t.Cleanup(func() {
		_ = adbCommand(t, "shell", "am", "force-stop", *packageName).Run()
	})

	adbOutput(t, "shell", "am", "force-stop", *packageName)
	adbOutput(t, "shell", "am", "start", "-n", fmt.Sprintf("%s/%s", *packageName, *activityName))

	var pid string
	deadline := time.Now().Add(30 * time.Second)
	for pid == "" && time.Now().Before(deadline) {
		output, err := adbCommand(t, "shell", "pidof", *packageName).Output()
		if err == nil {
			pid = strings.TrimSpace(string(output))
		} else if pid == "" {
			time.Sleep(500 * time.Millisecond)
		}
	}
	if pid == "" {
		t.Fatalf("ingester activity did not start successfully")
	}

	t.Logf("ingester process running with pid %s", pid)
}

func init() {
	flag.Parse()
	if *packageName == "" {
		*packageName = os.Getenv("PACKAGE_NAME")
	}
	if *activityName == "" {
		*activityName = os.Getenv("ACTIVITY_NAME")
	}
	if *adbPath == "" {
		home, _ := os.UserHomeDir()
		candidate := filepath.Join(home, "Android", "Sdk", "platform-tools", "adb")
		if info, err := os.Stat(candidate); err == nil && info.Mode()&0o111 != 0 {
			*adbPath = candidate
		}
	}
	if *adbPath == "" {
		panic("adb not found; set -adb or install Android platform-tools at ~/Android/Sdk/platform-tools/adb")
	}
	if *packageName == "" {
		panic("-package_name is required")
	}
	if *activityName == "" {
		panic("-activity_name is required")
	}
}
