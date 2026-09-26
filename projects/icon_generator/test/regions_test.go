package icon_generator_test

import (
	"fmt"
	"image"
	"math"
	"path/filepath"
	"slices"
	"strconv"
	"testing"
)

func TestRegionalCompatibility(t *testing.T) {
	dir := t.TempDir()
	for _, algorithm := range []string{"compact", "strands", "replication"} {
		args := []string{"--width", "96", "--height", "64", "--pixel-size", "1", "--density", "0.03", "--seed", "42", "--clusterization", "0.5", "--clusterization-algorithm", algorithm, "--replication-inheritance", "0.65", "--replication-forward-bias", "0.65", "--replication-crowding", "0.75", "--replication-branching", "0.2"}
		success(t, dir, append(slices.Clone(args), "--output", algorithm+".png")...)
		variation := "1"
		if algorithm == "replication" {
			variation = "0"
		}
		success(t, dir, append(args, "--replication-field-scale", "1", "--replication-growth-variation", variation, "--replication-seeding-variation", variation, "--output", algorithm+"-regional.png")...)
		if digest(decode(t, filepath.Join(dir, algorithm+".png"))) != digest(decode(t, filepath.Join(dir, algorithm+"-regional.png"))) {
			t.Fatalf("inactive regional controls changed %s", algorithm)
		}
	}
	args := []string{"--width", "64", "--height", "32", "--pixel-size", "1", "--density", "0.03", "--seed", "42", "--clusterization-algorithm", "replication", "--clusterization", "0"}
	success(t, dir, append(slices.Clone(args), "--output", "zero.png")...)
	success(t, dir, append(args, "--replication-growth-variation", "1", "--output", "zero-regional.png")...)
	if digest(decode(t, filepath.Join(dir, "zero.png"))) != digest(decode(t, filepath.Join(dir, "zero-regional.png"))) {
		t.Fatal("growth variation introduced births at zero base probability")
	}
}

func regionalContrast(im image.Image, density float64) float64 {
	score := 0.0
	for top := 0; top < 256; top += 32 {
		for left := 0; left < 256; left += 32 {
			count := 0
			for y := top; y < top+32; y++ {
				for x := left; x < left+32; x++ {
					if pixel(im, x, y).A != 0 {
						count++
					}
				}
			}
			delta := float64(count) - density*32*32
			score += delta * delta
		}
	}
	return score
}

func TestRegionalSeedingDensityAndContrast(t *testing.T) {
	dir := t.TempDir()
	for _, seed := range []string{"0", "42"} {
		for _, density := range []float64{0.03, 0.4, 0.9} {
			base := []string{"--width", "256", "--height", "256", "--pixel-size", "1", "--background", "transparent", "--density", strconv.FormatFloat(density, 'f', -1, 64), "--seed", seed, "--clusterization", "0", "--clusterization-algorithm", "replication", "--replication-field-scale", "64"}
			var scores [2]float64
			for i, strength := range []string{"0", "1"} {
				name := fmt.Sprintf("seed-%s-density-%g-variation-%s.png", seed, density, strength)
				success(t, dir, append(slices.Clone(base), "--replication-seeding-variation", strength, "--output", name)...)
				im := decode(t, filepath.Join(dir, name))
				count, _ := population(im)
				expected := density * 256 * 256
				tolerance := 6*math.Sqrt(expected*(1-density)) + 2
				if math.Abs(float64(count)-expected) > tolerance {
					t.Fatalf("average seeding probability drifted: %s count=%d expected=%g tolerance=%g", name, count, expected, tolerance)
				}
				scores[i] = regionalContrast(im, density)
			}
			if scores[1] <= 2*scores[0] {
				t.Fatalf("regional seeding lacks broad contrast: seed=%s density=%g scores=%v", seed, density, scores)
			}
		}
	}
}

func TestRegionalReproductionAndBoundaries(t *testing.T) {
	for _, tc := range []struct {
		width, height, size int
		scale               string
	}{
		{1, 1, 1, "4096"},
		{1, 7, 1, "1"},
		{7, 1, 1, "4096"},
		{31, 19, 4, "8"},
		{96, 64, 1, "16"},
	} {
		for _, density := range []string{"0", "0.1", "0.9", "1"} {
			t.Run(fmt.Sprintf("%dx%d-%s", tc.width, tc.height, density), func(t *testing.T) {
				dir := t.TempDir()
				args := []string{"--width", strconv.Itoa(tc.width), "--height", strconv.Itoa(tc.height), "--pixel-size", strconv.Itoa(tc.size), "--background", "transparent", "--colors", "#123abc,#abcdef", "--density", density, "--seed", "42", "--clusterization-algorithm", "replication", "--replication-field-scale", tc.scale, "--replication-seeding-variation", "1", "--replication-growth-variation", "1", "--replication-inheritance", "0.65", "--replication-forward-bias", "0.65", "--replication-crowding", "0.75", "--replication-branching", "0.2"}
				success(t, dir, append(slices.Clone(args), "--clusterization", "0", "--output", "initial.png")...)
				success(t, dir, append(slices.Clone(args), "--clusterization", "0.5", "--output", "grown.png")...)
				initial, grown := decode(t, filepath.Join(dir, "initial.png")), decode(t, filepath.Join(dir, "grown.png"))
				before, after := 0, 0
				for y := 0; y < tc.height; y += tc.size {
					for x := 0; x < tc.width; x += tc.size {
						if pixel(initial, x, y).A != 0 {
							before++
							if pixel(initial, x, y) != pixel(grown, x, y) {
								t.Fatal("regional growth removed or recolored a parent")
							}
						}
						if pixel(grown, x, y).A != 0 {
							after++
						}
					}
				}
				if after < before || (after-before)%2 != 0 {
					t.Fatalf("regional growth violated paired births: %d -> %d", before, after)
				}
				if (density == "0" || density == "1") && digest(initial) != digest(grown) {
					t.Fatal("regional growth changed a density endpoint")
				}
				success(t, dir, append(args, "--clusterization", "0.5", "--output", "replay.png")...)
				if digest(grown) != digest(decode(t, filepath.Join(dir, "replay.png"))) {
					t.Fatal("regional seed replay changed")
				}
			})
		}
	}
}
