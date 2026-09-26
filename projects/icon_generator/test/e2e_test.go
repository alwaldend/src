package icon_generator_test

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"
)

var generator = flag.String("generator", "", "path to the built generator")

func executable(t *testing.T) string {
	t.Helper()
	if *generator == "" {
		t.Fatal("-generator is required")
	}
	if filepath.IsAbs(*generator) {
		return *generator
	}
	return filepath.Join(os.Getenv("TEST_SRCDIR"), os.Getenv("TEST_WORKSPACE"), *generator)
}

func run(t *testing.T, dir string, args ...string) (string, string, error) {
	t.Helper()
	ctx, cancel := context.WithTimeout(t.Context(), 15*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, executable(t), args...)
	cmd.Dir = dir
	cmd.Env = []string{"PATH="}
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	if err := cmd.Run(); err != nil {
		return stdout.String(), stderr.String(), fmt.Errorf("run generator: %w", err)
	}
	return stdout.String(), stderr.String(), nil
}

func success(t *testing.T, dir string, args ...string) string {
	t.Helper()
	_, stderr, err := run(t, dir, args...)
	if err != nil {
		t.Fatalf("%v: %s", err, stderr)
	}
	return stderr
}

func decode(t *testing.T, path string) image.Image {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	im, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("decode %s: %v", path, err)
	}
	return im
}

func pixel(im image.Image, x, y int) color.NRGBA {
	return color.NRGBAModel.Convert(im.At(x, y)).(color.NRGBA)
}

func digest(im image.Image) string {
	h := sha256.New()
	fmt.Fprintf(h, "%d,%d\n", im.Bounds().Dx(), im.Bounds().Dy())
	for y := 0; y < im.Bounds().Dy(); y++ {
		for x := 0; x < im.Bounds().Dx(); x++ {
			c := pixel(im, x, y)
			h.Write([]byte{c.R, c.G, c.B, c.A})
		}
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

func reportedSeed(t *testing.T, stderr string) string {
	t.Helper()
	m := regexp.MustCompile(`(?m)^seed: ([0-9]+)$`).FindStringSubmatch(stderr)
	if len(m) != 2 {
		t.Fatalf("missing effective seed: %q", stderr)
	}
	return m[1]
}

func TestDefaultsAndHelp(t *testing.T) {
	dir := t.TempDir()
	stdout, stderr, err := run(t, dir, "--help")
	if err != nil {
		t.Fatalf("help: %v: %s", err, stderr)
	}
	for _, name := range []string{"width", "height", "background", "colors", "pixel-size", "pixel-shape", "density", "clusterization", "clusterization-algorithm", "replication-inheritance", "replication-forward-bias", "replication-crowding", "replication-branching", "replication-field-scale", "replication-growth-variation", "replication-seeding-variation", "seed", "output"} {
		if !strings.Contains(stdout+stderr, "-"+name) {
			t.Errorf("help omits %s", name)
		}
	}
	entries, err := os.ReadDir(dir)
	if err != nil || len(entries) != 0 {
		t.Fatalf("help modified directory: %v, %v", entries, err)
	}
	seed := reportedSeed(t, success(t, dir))
	im := decode(t, filepath.Join(dir, "icon.png"))
	if im.Bounds() != image.Rect(0, 0, 256, 256) {
		t.Fatalf("default bounds: %v", im.Bounds())
	}
	success(t, dir, "--width", "256", "--height", "256", "--background", "#ffffff", "--colors", "#000000", "--pixel-size", "16", "--pixel-shape", "square", "--density", "0.4", "--clusterization", "0", "--seed", seed, "--output", "explicit.png")
	if digest(im) != digest(decode(t, filepath.Join(dir, "explicit.png"))) {
		t.Fatal("explicit defaults do not reproduce the default invocation")
	}
}

func TestGeometry(t *testing.T) {
	for _, tc := range []struct {
		name, shape string
		size        int
		mask        []string
	}{
		{"squares", "square", 2, []string{"######", "######", "######", "######"}},
		{"odd-circle", "circle", 5, []string{".###.", "#####", "#####", "#####", ".###."}},
		{"even-circle", "circle", 4, []string{".##.", "####", "####", ".##."}},
		{"clipped-circle", "circle", 4, []string{".##..", "#####", "#####"}},
		{"clipped-square", "square", 4, []string{"#####", "#####", "#####"}},
		{"one-circle", "circle", 1, []string{"###", "###"}},
		{"one-square", "square", 1, []string{"###", "###"}},
		{"oversized-circle", "circle", 4, []string{"."}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			success(t, dir, "--width", strconv.Itoa(len(tc.mask[0])), "--height", strconv.Itoa(len(tc.mask)), "--pixel-size", strconv.Itoa(tc.size), "--pixel-shape", tc.shape, "--density", "1", "--clusterization", "1", "--colors", "#123AbC", "--background", "transparent", "--seed", "0")
			im := decode(t, filepath.Join(dir, "icon.png"))
			if im.Bounds() != image.Rect(0, 0, len(tc.mask[0]), len(tc.mask)) {
				t.Fatalf("bounds: %v", im.Bounds())
			}
			for y, row := range tc.mask {
				for x, c := range row {
					want := color.NRGBA{}
					if c == '#' {
						want = color.NRGBA{R: 0x12, G: 0x3a, B: 0xbc, A: 255}
					}
					if got := pixel(im, x, y); got != want {
						t.Fatalf("pixel %d,%d: got %v want %v", x, y, got, want)
					}
				}
			}
		})
	}
}

func TestBackgroundAndPalette(t *testing.T) {
	for _, background := range []string{"#aBcDeF", "transparent"} {
		for _, shape := range []string{"square", "circle"} {
			for _, density := range []string{"0", "0.4", "1"} {
				t.Run(background+"/"+shape+"/"+density, func(t *testing.T) {
					dir := t.TempDir()
					success(t, dir, "--width", "120", "--height", "80", "--background", background, "--colors", "#123abc,#abcdef", "--pixel-shape", shape, "--pixel-size", "7", "--density", density, "--clusterization", "0.75", "--seed", "18446744073709551615")
					im := decode(t, filepath.Join(dir, "icon.png"))
					if im.Bounds() != image.Rect(0, 0, 120, 80) {
						t.Fatalf("bounds: %v", im.Bounds())
					}
					bg := color.NRGBA{R: 0xab, G: 0xcd, B: 0xef, A: 255}
					if background == "transparent" {
						bg = color.NRGBA{}
					}
					allowed := []color.NRGBA{bg}
					if density != "0" {
						allowed = append(allowed, color.NRGBA{R: 0x12, G: 0x3a, B: 0xbc, A: 255}, color.NRGBA{R: 0xab, G: 0xcd, B: 0xef, A: 255})
					}
					for y := 0; y < 80; y++ {
						for x := 0; x < 120; x++ {
							if c := pixel(im, x, y); !slices.Contains(allowed, c) {
								t.Fatalf("unexpected color at %d,%d: %v", x, y, c)
							}
						}
					}
				})
			}
		}
	}
}

func TestSeedReplay(t *testing.T) {
	for _, seed := range []string{"", "0", "42"} {
		t.Run("seed-"+seed, func(t *testing.T) {
			dir := t.TempDir()
			args := []string{"--width", "63", "--height", "45", "--pixel-size", "7", "--pixel-shape", "circle", "--colors", "#010203,#abcdef,#987654", "--clusterization", "0.75"}
			first := slices.Clone(args)
			if seed != "" {
				first = append(first, "--seed", seed)
			}
			effective := reportedSeed(t, success(t, dir, first...))
			if seed != "" && effective != seed {
				t.Fatalf("seed changed: %s", effective)
			}
			success(t, dir, append(args, "--seed", effective, "--output", "replay.png")...)
			if digest(decode(t, filepath.Join(dir, "icon.png"))) != digest(decode(t, filepath.Join(dir, "replay.png"))) {
				t.Fatal("seed replay differs")
			}
		})
	}
}

func TestInvalidOptions(t *testing.T) {
	cases := [][]string{
		{"--width", "0"},
		{"--height", "-1"},
		{"--width", "4097"},
		{"--height", "4097"},
		{"--width", "1.5"},
		{"--pixel-size", "0"},
		{"--pixel-size", "-1"},
		{"--pixel-size", "4097"},
		{"--density", "-0.1"},
		{"--density", "1.1"},
		{"--density", "NaN"},
		{"--density", "Inf"},
		{"--clusterization", "-0.1"},
		{"--clusterization", "1.1"},
		{"--clusterization", "NaN"},
		{"--clusterization", "Inf"},
		{"--clusterization", "-Inf"},
		{"--clusterization", "bad"},
		{"--clusterization-algorithm", "unknown"},
		{"--clusterization-algorithm", ""},
		{"--colors", ""},
		{"--colors", "#000000,"},
		{"--colors", "#gg0000"},
		{"--colors", "#123"},
		{"--colors", "transparent"},
		{"--background", "white"},
		{"--background", "#000000ff"},
		{"--pixel-shape", "triangle"},
		{"--seed", "-1"},
		{"--seed", "18446744073709551616"},
		{"--seed", "bad"},
		{"--unknown"},
		{"--width"},
		{"argument"},
		{"--output", ""},
	}
	for _, option := range []string{"--replication-inheritance", "--replication-forward-bias", "--replication-crowding", "--replication-branching", "--replication-growth-variation", "--replication-seeding-variation"} {
		for _, value := range []string{"-0.1", "1.1", "NaN", "Inf", "-Inf", "bad"} {
			cases = append(cases, []string{option, value})
		}
	}
	for _, value := range []string{"0", "-1", "4097", "1.5", "NaN", "bad"} {
		cases = append(cases, []string{"--replication-field-scale", value})
	}
	for _, args := range cases {
		t.Run(strings.Join(args, "="), func(t *testing.T) {
			dir := t.TempDir()
			_, stderr, err := run(t, dir, args...)
			if err == nil || stderr == "" {
				t.Fatalf("expected failure and diagnostic: %v %q", err, stderr)
			}
			entries, err := os.ReadDir(dir)
			if err != nil || len(entries) != 0 {
				t.Fatalf("invalid options created output: %v %v", entries, err)
			}
		})
	}
}

func TestClusterization(t *testing.T) {
	for _, size := range []image.Point{{64, 48}, {1, 1}, {1, 17}, {19, 1}, {3, 2}, {4, 4}, {5, 5}} {
		for _, density := range []string{"0", "0.2", "0.5", "0.8", "1"} {
			for _, seed := range []string{"0", "1", "2", "3", "42"} {
				t.Run(fmt.Sprintf("%dx%d/density-%s/seed-%s", size.X, size.Y, density, seed), func(t *testing.T) {
					dir := t.TempDir()
					initialCount := 0
					for i, factor := range []string{"0", "0.25", "0.75", "1"} {
						out := "factor-" + factor + ".png"
						success(t, dir, "--width", strconv.Itoa(size.X), "--height", strconv.Itoa(size.Y), "--pixel-size", "1", "--background", "transparent", "--colors", "#ffffff", "--density", density, "--clusterization", factor, "--clusterization-algorithm", "strands", "--seed", seed, "--output", out)
						im := decode(t, filepath.Join(dir, out))
						count, contacts, isolated, crowded, blocks := 0, 0, 0, 0, 0
						for y := 0; y < size.Y; y++ {
							for x := 0; x < size.X; x++ {
								if pixel(im, x, y).A == 0 {
									continue
								}
								count++
								neighbors := 0
								for dy := -1; dy <= 1; dy++ {
									for dx := -1; dx <= 1; dx++ {
										if (dx != 0 || dy != 0) && image.Pt(x+dx, y+dy).In(im.Bounds()) && pixel(im, x+dx, y+dy).A != 0 {
											neighbors++
										}
									}
								}
								if neighbors == 0 {
									isolated++
								}
								if neighbors > 4 {
									crowded++
								}
								if x+1 < size.X && y+1 < size.Y && pixel(im, x+1, y).A != 0 && pixel(im, x, y+1).A != 0 && pixel(im, x+1, y+1).A != 0 {
									blocks++
								}
								if x+1 < size.X && pixel(im, x+1, y).A != 0 {
									contacts++
								}
								if y+1 < size.Y && pixel(im, x, y+1).A != 0 {
									contacts++
								}
								if x+1 < size.X && y+1 < size.Y && pixel(im, x+1, y+1).A != 0 {
									contacts++
								}
								if x > 0 && y+1 < size.Y && pixel(im, x-1, y+1).A != 0 {
									contacts++
								}
							}
						}
						if i == 0 {
							initialCount = count
						}
						if count != initialCount {
							t.Fatalf("factor %s: count %d (initial %d)", factor, count, initialCount)
						}
						if size.X == 64 && density == "0.2" && (factor == "0.75" || factor == "1") {
							if isolated*10 > count || crowded*10 > count || blocks*25 > count {
								t.Fatalf("factor %s lacks open strands: cells=%d, pairs=%d, isolated=%d, crowded=%d, filled blocks=%d", factor, count, contacts, isolated, crowded, blocks)
							}
						}
					}
				})
			}
		}
	}
}

func TestOutputProtection(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "existing.png")
	before := []byte("preserve this file")
	if err := os.WriteFile(path, before, 0o600); err != nil {
		t.Fatal(err)
	}
	for _, out := range []string{path, dir, filepath.Join(dir, "missing", "icon.png")} {
		_, stderr, err := run(t, dir, "--output", out)
		if err == nil || !strings.Contains(stderr, out) {
			t.Fatalf("output failure: %v %q", err, stderr)
		}
	}
	got, err := os.ReadFile(path)
	if err != nil || !bytes.Equal(got, before) {
		t.Fatalf("existing file changed: %q %v", got, err)
	}
	if _, err := os.Stat(filepath.Join(dir, "missing")); !os.IsNotExist(err) {
		t.Fatalf("parent directory created: %v", err)
	}
	link := filepath.Join(dir, "link.png")
	if err := os.Symlink(path, link); err != nil {
		t.Fatal(err)
	}
	if _, _, err := run(t, dir, "--output", link); err == nil {
		t.Fatal("overwrote symlink")
	}
	if target, err := os.Readlink(link); err != nil || target != path {
		t.Fatalf("symlink changed: %s %v", target, err)
	}
	if os.Geteuid() != 0 {
		locked := filepath.Join(dir, "locked")
		if err := os.Mkdir(locked, 0o500); err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() {
			if err := os.Chmod(locked, 0o700); err != nil {
				t.Error(err)
			}
		})
		if _, _, err := run(t, dir, "--output", filepath.Join(locked, "icon.png")); err == nil {
			t.Fatal("wrote in unwritable directory")
		}
	}
}

func TestExamples(t *testing.T) {
	dir := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR")
	if dir == "" {
		dir = t.TempDir()
	}
	type example struct {
		Arguments   []string `json:"arguments"`
		Seed        string   `json:"seed"`
		PixelSHA256 string   `json:"pixel_sha256"`
	}
	var manifest []example
	for _, tc := range []struct{ name, shape, bg, factor, algorithm, density string }{
		{"squares", "square", "#f5f1e8", "0", "strands", "0.4"},
		{"circles", "circle", "#f5f1e8", "0", "strands", "0.4"},
		{"transparent", "circle", "transparent", "0", "strands", "0.4"},
		{"clustered", "circle", "#f5f1e8", "0.75", "strands", "0.4"},
		{"compact", "circle", "#f5f1e8", "0.75", "compact", "0.4"},
		{"replication", "circle", "#f5f1e8", "0.4", "replication", "0.04"},
		{"replication-growth", "square", "#212121", "0.5", "replication", "0.04"},
		{"replication-regions", "square", "#212121", "0.5", "replication", "0.04"},
	} {
		args := []string{"--width", "256", "--height", "256", "--pixel-size", "16", "--pixel-shape", tc.shape, "--background", tc.bg, "--colors", "#223843,#d77a61,#e3b23c", "--density", tc.density, "--clusterization", tc.factor, "--clusterization-algorithm", tc.algorithm, "--seed", "42", "--output", tc.name + ".png"}
		if tc.name == "replication-growth" || tc.name == "replication-regions" {
			args = append(args, "--pixel-size", "1", "--replication-inheritance", "0.9", "--replication-forward-bias", "0.8", "--replication-crowding", "0.75", "--replication-branching", "0.15")
		}
		if tc.name == "replication-regions" {
			args = append(args, "--replication-field-scale", "64", "--replication-growth-variation", "0.1", "--replication-seeding-variation", "0.8")
		}
		seed := reportedSeed(t, success(t, dir, args...))
		im := decode(t, filepath.Join(dir, tc.name+".png"))
		manifest = append(manifest, example{args, seed, digest(im)})
	}
	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "manifest.json"), append(data, '\n'), 0o600); err != nil {
		t.Fatal(err)
	}
}
