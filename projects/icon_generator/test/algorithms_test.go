package icon_generator_test

import (
	"image"
	"path/filepath"
	"slices"
	"strconv"
	"testing"
)

func population(im image.Image) (count, pairs int) {
	bounds := im.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if pixel(im, x, y).A == 0 {
				continue
			}
			count++
			for _, offset := range []image.Point{{1, 0}, {0, 1}, {1, 1}, {-1, 1}} {
				p := image.Pt(x+offset.X, y+offset.Y)
				if p.In(bounds) && pixel(im, p.X, p.Y).A != 0 {
					pairs++
				}
			}
		}
	}
	return count, pairs
}

func TestAlgorithmSelection(t *testing.T) {
	dir := t.TempDir()
	base := []string{"--width", "64", "--height", "48", "--pixel-size", "1", "--background", "transparent", "--colors", "#123abc,#abcdef", "--density", "0.2", "--seed", "42"}
	success(t, dir, append(slices.Clone(base), "--clusterization", "0", "--output", "independent.png")...)
	success(t, dir, append(slices.Clone(base), "--clusterization", "0.75", "--output", "default-compact.png")...)
	independent := decode(t, filepath.Join(dir, "independent.png"))
	initial, _ := population(independent)
	for _, algorithm := range []string{"compact", "strands", "replication"} {
		t.Run(algorithm, func(t *testing.T) {
			var previousPairs int
			for _, factor := range []string{"0", "0.25", "0.75", "1"} {
				name := algorithm + "-" + factor + ".png"
				args := append(slices.Clone(base), "--clusterization-algorithm", algorithm, "--clusterization", factor)
				success(t, dir, append(slices.Clone(args), "--output", name)...)
				im := decode(t, filepath.Join(dir, name))
				count, pairs := population(im)
				if factor == "0" && digest(im) != digest(independent) {
					t.Fatal("zero factor changed independent placement")
				}
				if algorithm == "replication" {
					if count < initial || (count-initial)%2 != 0 {
						t.Fatalf("offspring must arrive in pairs: initial=%d final=%d", initial, count)
					}
					for y := 0; y < 48; y++ {
						for x := 0; x < 64; x++ {
							if pixel(independent, x, y).A != 0 && pixel(im, x, y) != pixel(independent, x, y) {
								t.Fatal("replication removed or recolored an existing cell")
							}
						}
					}
				} else if count != initial {
					t.Fatalf("sampled population changed: %d -> %d", initial, count)
				}
				if algorithm == "compact" && pairs < previousPairs {
					t.Fatalf("compact adjacency decreased: %d -> %d", previousPairs, pairs)
				}
				previousPairs = pairs
				if algorithm == "compact" && factor == "0.75" && digest(im) != digest(decode(t, filepath.Join(dir, "default-compact.png"))) {
					t.Fatal("explicit compact differs from the default")
				}
				success(t, dir, append(args, "--output", "replay-"+name)...)
				if digest(im) != digest(decode(t, filepath.Join(dir, "replay-"+name))) {
					t.Fatal("algorithm seed replay changed")
				}
			}
		})
	}
}

func TestReplicationBirths(t *testing.T) {
	for _, tc := range []struct {
		name          string
		width, height int
		seed          string
		want          int
	}{
		{"two-adjacent-children", 3, 1, "1", 3},
		{"two-distinct-children", 2, 2, "4", 3},
		{"one-vacancy-is-insufficient", 2, 1, "1", 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			args := []string{"--width", strconv.Itoa(tc.width), "--height", strconv.Itoa(tc.height), "--pixel-size", "1", "--background", "transparent", "--density", "0.2", "--seed", tc.seed, "--clusterization-algorithm", "replication"}
			success(t, dir, append(slices.Clone(args), "--clusterization", "0", "--output", "initial.png")...)
			initial := decode(t, filepath.Join(dir, "initial.png"))
			if count, _ := population(initial); count != 1 {
				t.Fatalf("fixture needs one initial cell, got %d", count)
			}
			args = append(args, "--clusterization", "1")
			success(t, dir, append(slices.Clone(args), "--output", "offspring.png")...)
			if count, _ := population(decode(t, filepath.Join(dir, "offspring.png"))); count != tc.want {
				t.Fatalf("population: got %d want %d", count, tc.want)
			}
			args = append(args, "--replication-inheritance", "1", "--replication-forward-bias", "1", "--replication-crowding", "1", "--replication-branching", "1", "--output", "steered.png")
			success(t, dir, args...)
			if count, _ := population(decode(t, filepath.Join(dir, "steered.png"))); count != tc.want {
				t.Fatalf("preferences prevented paired births: got %d want %d", count, tc.want)
			}
		})
	}
}

func TestReplicationGrowthControls(t *testing.T) {
	dir := t.TempDir()
	base := []string{"--width", "96", "--height", "64", "--pixel-size", "1", "--background", "transparent", "--colors", "#123abc,#abcdef", "--density", "0.02", "--seed", "42", "--clusterization-algorithm", "replication"}
	success(t, dir, append(slices.Clone(base), "--clusterization", "0", "--output", "initial.png")...)
	initial := decode(t, filepath.Join(dir, "initial.png"))
	initialCount, _ := population(initial)
	base = append(base, "--clusterization", "0.5")
	success(t, dir, append(slices.Clone(base), "--output", "default.png")...)
	for i, controls := range [][]string{
		{"0", "0", "0", "0"},
		{"1", "0", "0", "0"},
		{"0", "1", "0", "0"},
		{"1", "1", "0", "0"},
		{"0", "0", "1", "0"},
		{"0", "0", "0", "1"},
		{"0.9", "0.8", "0.75", "0.15"},
		{"1", "1", "1", "1"},
	} {
		name := strconv.Itoa(i) + ".png"
		args := append(slices.Clone(base), "--replication-inheritance", controls[0], "--replication-forward-bias", controls[1], "--replication-crowding", controls[2], "--replication-branching", controls[3])
		success(t, dir, append(slices.Clone(args), "--output", name)...)
		im := decode(t, filepath.Join(dir, name))
		count, _ := population(im)
		if count < initialCount || (count-initialCount)%2 != 0 {
			t.Fatalf("controls %v violated paired births: %d -> %d", controls, initialCount, count)
		}
		for y := 0; y < 64; y++ {
			for x := 0; x < 96; x++ {
				if pixel(initial, x, y).A != 0 && pixel(initial, x, y) != pixel(im, x, y) {
					t.Fatalf("controls %v removed or recolored a parent", controls)
				}
			}
		}
		if i < 2 && digest(im) != digest(decode(t, filepath.Join(dir, "default.png"))) {
			t.Fatalf("inactive placement controls changed baseline: %v", controls)
		}
		success(t, dir, append(args, "--output", "replay-"+name)...)
		if digest(im) != digest(decode(t, filepath.Join(dir, "replay-"+name))) {
			t.Fatalf("controls %v failed seed replay", controls)
		}
	}
	for _, algorithm := range []string{"compact", "strands", "replication"} {
		args := append(slices.Clone(base), "--clusterization-algorithm", algorithm)
		if algorithm == "replication" {
			args = append(args, "--clusterization", "0")
		}
		success(t, dir, append(slices.Clone(args), "--output", algorithm+".png")...)
		args = append(args, "--replication-inheritance", "1", "--replication-forward-bias", "1", "--replication-crowding", "1", "--replication-branching", "1", "--output", algorithm+"-controls.png")
		success(t, dir, args...)
		if digest(decode(t, filepath.Join(dir, algorithm+".png"))) != digest(decode(t, filepath.Join(dir, algorithm+"-controls.png"))) {
			t.Fatalf("controls affected inactive replication in %s", algorithm)
		}
	}
}

func TestReplicationCrowdingPreference(t *testing.T) {
	dir := t.TempDir()
	var counts, pairs [2]int
	for _, seed := range []string{"0", "42", "123", "999"} {
		for mode, crowding := range []string{"0", "1"} {
			name := "crowding-" + crowding + "-" + seed + ".png"
			success(t, dir, "--width", "128", "--height", "128", "--pixel-size", "1", "--background", "transparent", "--density", "0.005", "--clusterization", "0.4", "--clusterization-algorithm", "replication", "--replication-crowding", crowding, "--seed", seed, "--output", name)
			count, contacts := population(decode(t, filepath.Join(dir, name)))
			counts[mode] += count
			pairs[mode] += contacts
		}
	}
	if float64(pairs[1])/float64(counts[1]) >= float64(pairs[0])/float64(counts[0]) {
		t.Fatalf("sparse crowding preference did not reduce contacts per cell: counts=%v pairs=%v", counts, pairs)
	}
}

func TestReplicationGenerationsAndEndpoints(t *testing.T) {
	for _, density := range []string{"0", "0.02", "1"} {
		t.Run(density, func(t *testing.T) {
			dir := t.TempDir()
			args := []string{"--width", "32", "--height", "24", "--pixel-size", "1", "--background", "transparent", "--density", density, "--seed", "42", "--clusterization-algorithm", "replication"}
			success(t, dir, append(slices.Clone(args), "--clusterization", "0", "--output", "initial.png")...)
			success(t, dir, append(args, "--clusterization", "1", "--output", "grown.png")...)
			initial, _ := population(decode(t, filepath.Join(dir, "initial.png")))
			grown, _ := population(decode(t, filepath.Join(dir, "grown.png")))
			if density == "0.02" {
				if initial == 0 || grown <= 3*initial {
					t.Fatalf("descendants did not reproduce: %d -> %d", initial, grown)
				}
			} else if grown != initial {
				t.Fatalf("density endpoint changed: %d -> %d", initial, grown)
			}
		})
	}
}
