package icon_generator_test

import (
	"image"
	"path/filepath"
	"slices"
	"strconv"
	"testing"
)

func TestTipsGrowth(t *testing.T) {
	dir := t.TempDir()
	base := []string{"--width", "384", "--height", "384", "--pixel-size", "1", "--background", "transparent", "--colors", "#123abc,#abcdef", "--density", "0.0005", "--seed", "42", "--clusterization-algorithm", "tips"}
	success(t, dir, append(slices.Clone(base), "--clusterization", "0", "--output", "initial.png")...)
	initial := decode(t, filepath.Join(dir, "initial.png"))
	success(t, dir, append(slices.Clone(base), "--clusterization-algorithm", "compact", "--output", "independent.png")...)
	if digest(initial) != digest(decode(t, filepath.Join(dir, "independent.png"))) {
		t.Fatal("zero tip growth changed initial sampling")
	}
	base = append(base, "--clusterization", "0.8")
	success(t, dir, append(slices.Clone(base), "--output", "grown.png")...)
	grown := decode(t, filepath.Join(dir, "grown.png"))
	startCount, _ := population(initial)
	endCount, _ := population(grown)
	if startCount == 0 || endCount < 5*startCount || endCount > 384*384/10 {
		t.Fatalf("sparse tips need substantial growth with open space: %d -> %d", startCount, endCount)
	}
	for y := 0; y < 384; y++ {
		for x := 0; x < 384; x++ {
			if pixel(initial, x, y).A != 0 && pixel(grown, x, y) != pixel(initial, x, y) {
				t.Fatal("tip growth removed or recolored an initial cell")
			}
		}
	}
	// Every visible connected component must descend from a sampled root.
	seen := make(map[image.Point]bool)
	longest := 0
	for y := 0; y < 384; y++ {
		for x := 0; x < 384; x++ {
			start := image.Pt(x, y)
			if seen[start] || pixel(grown, x, y).A == 0 {
				continue
			}
			queue := []image.Point{start}
			seen[start] = true
			rooted := false
			bounds := image.Rect(x, y, x+1, y+1)
			for head := 0; head < len(queue); head++ {
				p := queue[head]
				rooted = rooted || pixel(initial, p.X, p.Y).A != 0
				bounds = bounds.Union(image.Rect(p.X, p.Y, p.X+1, p.Y+1))
				for dy := -1; dy <= 1; dy++ {
					for dx := -1; dx <= 1; dx++ {
						next := p.Add(image.Pt(dx, dy))
						if next.In(grown.Bounds()) && !seen[next] && pixel(grown, next.X, next.Y).A != 0 {
							seen[next] = true
							queue = append(queue, next)
						}
					}
				}
			}
			if !rooted {
				t.Fatal("growth produced an unconnected new component")
			}
			longest = max(longest, bounds.Dx(), bounds.Dy())
		}
	}
	if longest < 32 {
		t.Fatalf("no extended strand: largest extent %d", longest)
	}
	for i, controls := range [][]string{
		nil,
		{"--replication-inheritance", "1", "--replication-forward-bias", "1", "--replication-crowding", "1", "--replication-branching", "1", "--replication-field-scale", "1", "--replication-growth-variation", "1", "--replication-seeding-variation", "1"},
	} {
		name := "replay-" + strconv.Itoa(i) + ".png"
		args := append(slices.Clone(base), controls...)
		success(t, dir, append(args, "--output", name)...)
		if digest(grown) != digest(decode(t, filepath.Join(dir, name))) {
			t.Fatal("replay or replication-only options changed tips output")
		}
	}
}

func TestTipsConstrainedCanvases(t *testing.T) {
	for _, size := range []image.Point{{1, 1}, {1, 47}, {53, 1}, {19, 23}} {
		for _, density := range []string{"0", "0.2", "1"} {
			for _, shape := range []string{"square", "circle"} {
				dir := t.TempDir()
				base := []string{"--width", strconv.Itoa(size.X), "--height", strconv.Itoa(size.Y), "--pixel-size", "3", "--pixel-shape", shape, "--background", "transparent", "--density", density, "--seed", "0", "--clusterization-algorithm", "tips"}
				success(t, dir, append(slices.Clone(base), "--clusterization", "0", "--output", "initial.png")...)
				success(t, dir, append(base, "--clusterization", "1", "--output", "grown.png")...)
				initial := decode(t, filepath.Join(dir, "initial.png"))
				grown := decode(t, filepath.Join(dir, "grown.png"))
				if grown.Bounds().Size() != size {
					t.Fatalf("changed bounds: %v", grown.Bounds())
				}
				if density != "0.2" && digest(initial) != digest(grown) {
					t.Fatal("growth changed empty or full density")
				}
			}
		}
	}
}

func TestTipsCanExtendOneChild(t *testing.T) {
	// A two-cell canvas cannot support paired reproduction, but a growing tip
	// pointing into its sole vacancy must be able to extend once.
	dir := t.TempDir()
	for seed := range 32 {
		base := []string{"--width", "2", "--height", "1", "--pixel-size", "1", "--background", "transparent", "--density", "0.5", "--seed", strconv.Itoa(seed), "--clusterization-algorithm", "tips"}
		initialName, grownName := "initial-"+strconv.Itoa(seed)+".png", "grown-"+strconv.Itoa(seed)+".png"
		success(t, dir, append(slices.Clone(base), "--output", initialName)...)
		initial, _ := population(decode(t, filepath.Join(dir, initialName)))
		if initial != 1 {
			continue
		}
		success(t, dir, append(base, "--clusterization", "1", "--output", grownName)...)
		grown, _ := population(decode(t, filepath.Join(dir, grownName)))
		if grown == 2 {
			return
		}
	}
	t.Fatal("tips never extended into a sole adjacent vacancy")
}
