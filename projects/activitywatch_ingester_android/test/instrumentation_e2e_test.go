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
	e2eAdbPath     = flag.String("e2e_adb", "", "Path to adb")
	e2eAppApk      = flag.String("e2e_app_apk", "", "Path to app APK")
	e2eTestApk     = flag.String("e2e_test_apk", "", "Path to instrumentation APK")
	e2eTestPackage = flag.String("e2e_test_package", "", "Instrumentation package")
	e2eAppPackage  = flag.String("e2e_app_package", "", "Application package")
	e2eTestRunner  = flag.String("e2e_test_runner", "", "Instrumentation runner class")
)

func e2eAdbCommand(t *testing.T, args ...string) *exec.Cmd {
	t.Helper()
	path := os.Getenv("E2E_ADB")
	if path == "" {
		path = *e2eAdbPath
	}
	if path == "" {
		home, _ := os.UserHomeDir()
		path = filepath.Join(home, "Android", "Sdk", "platform-tools", "adb")
	}
	return exec.Command(path, args...)
}

func e2eAdbOutput(t *testing.T, args ...string) string {
	t.Helper()
	cmd := e2eAdbCommand(t, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("adb %s: %v\n%s", strings.Join(args, " "), err, out)
	}
	return string(out)
}

func TestInstrumentationE2E(t *testing.T) {
	appApk, err := runfiles.Rlocation("_main/projects/activitywatch_ingester_android/main/java/ingester_binary.apk")
	if err != nil {
		t.Fatalf("resolve app APK: %v", err)
	}
	testApk, err := runfiles.Rlocation("_main/projects/activitywatch_ingester_android/test/instrumentation_apk.apk")
	if err != nil {
		t.Fatalf("resolve test APK: %v", err)
	}

	e2eAdbOutput(t, "devices")
	e2eAdbOutput(t, "install", "-r", appApk)
	e2eAdbOutput(t, "install", "-r", testApk)

	instrumentCmd := []string{
		"shell", "am", "instrument", "-w", "-e", "class",
		"com.alwaldend.src.projects.activitywatch_ingester_android.ActivityWatchInstrumentationTest",
		fmt.Sprintf("%s/%s", *e2eTestPackage, *e2eTestRunner),
	}

	var output string
	deadline := time.Now().Add(2 * time.Minute)
	for time.Now().Before(deadline) {
		output = e2eAdbOutput(t, instrumentCmd...)
		if strings.Contains(output, "OK (") || strings.Contains(output, "FAILURES") {
			break
		}
		time.Sleep(2 * time.Second)
	}

	if strings.Contains(output, "FAILURES") {
		t.Fatalf("instrumentation test failed:\n%s", output)
	}
	if !strings.Contains(output, "OK (") {
		t.Fatalf("instrumentation did not complete within timeout:\n%s", output)
	}
	t.Logf("instrumentation output:\n%s", output)
}

func init() {
	flag.Parse()
	if *e2eTestRunner == "" {
		*e2eTestRunner = "androidx.test.runner.AndroidJUnitRunner"
	}
	if *e2eTestPackage == "" {
		*e2eTestPackage = os.Getenv("E2E_TEST_PACKAGE")
	}
}
