// Package generator renders random grid patterns and writes them as PNGs.
package generator

import (
	cryptorand "crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"math/rand/v2"
	"os"
	"strconv"
	"strings"
)

// Options describes the canvas, grid, and output. A nil Seed selects a fresh seed.
type Options struct {
	Width, Height, PixelSize    int
	Background, Colors, Shape   string
	Density                     float64
	Clusterization              float64
	ClusterAlgorithm            string
	ReplicationInheritance      float64
	ReplicationForwardBias      float64
	ReplicationCrowding         float64
	ReplicationBranching        float64
	ReplicationFieldScale       int
	ReplicationGrowthVariation  float64
	ReplicationSeedingVariation float64
	Seed                        *uint64
	Output                      string
}

// Generate validates all options before writing a new PNG and returns its seed.
// Existing destinations are never overwritten.
func Generate(options Options) (uint64, error) {
	background, palette, err := validate(options)
	if err != nil {
		return 0, fmt.Errorf("validate image options: %w", err)
	}
	var seed uint64
	if options.Seed != nil {
		seed = *options.Seed
	} else {
		var data [8]byte
		if _, err := cryptorand.Read(data[:]); err != nil {
			return 0, fmt.Errorf("generate random seed: %w", err)
		}
		seed = binary.LittleEndian.Uint64(data[:])
	}
	im := render(options, background, palette, seed)
	if err := writePNG(options.Output, im); err != nil {
		return 0, fmt.Errorf("save generated image: %w", err)
	}
	return seed, nil
}

func validate(options Options) (color.NRGBA, []color.NRGBA, error) {
	for _, dimension := range []struct {
		name  string
		value int
	}{
		{"width", options.Width},
		{"height", options.Height},
		{"pixel-size", options.PixelSize},
		{"replication-field-scale", options.ReplicationFieldScale},
	} {
		if dimension.value < 1 || dimension.value > 4096 {
			return color.NRGBA{}, nil, fmt.Errorf("--%s must be an integer from 1 through 4096", dimension.name)
		}
	}
	if math.IsNaN(options.Density) || math.IsInf(options.Density, 0) || options.Density < 0 || options.Density > 1 {
		return color.NRGBA{}, nil, fmt.Errorf("--density must be a finite number from 0 through 1")
	}
	if math.IsNaN(options.Clusterization) || math.IsInf(options.Clusterization, 0) || options.Clusterization < 0 || options.Clusterization > 1 {
		return color.NRGBA{}, nil, fmt.Errorf("--clusterization must be a finite number from 0 through 1")
	}
	if options.Shape != "square" && options.Shape != "circle" {
		return color.NRGBA{}, nil, fmt.Errorf("--pixel-shape must be square or circle")
	}
	if options.ClusterAlgorithm != "compact" && options.ClusterAlgorithm != "strands" && options.ClusterAlgorithm != "replication" {
		return color.NRGBA{}, nil, fmt.Errorf("--clusterization-algorithm must be compact, strands, or replication")
	}
	for _, control := range []struct {
		name  string
		value float64
	}{
		{"replication-inheritance", options.ReplicationInheritance},
		{"replication-forward-bias", options.ReplicationForwardBias},
		{"replication-crowding", options.ReplicationCrowding},
		{"replication-branching", options.ReplicationBranching},
		{"replication-growth-variation", options.ReplicationGrowthVariation},
		{"replication-seeding-variation", options.ReplicationSeedingVariation},
	} {
		if math.IsNaN(control.value) || math.IsInf(control.value, 0) || control.value < 0 || control.value > 1 {
			return color.NRGBA{}, nil, fmt.Errorf("--%s must be a finite number from 0 through 1", control.name)
		}
	}
	if options.Output == "" {
		return color.NRGBA{}, nil, fmt.Errorf("--output must name a file")
	}
	var background color.NRGBA
	if options.Background != "transparent" {
		parsed, err := parseColor(options.Background)
		if err != nil {
			return color.NRGBA{}, nil, fmt.Errorf("parse --background: %w", err)
		}
		background = parsed
	}
	var palette []color.NRGBA
	for _, entry := range strings.Split(options.Colors, ",") {
		parsed, err := parseColor(strings.TrimSpace(entry))
		if err != nil {
			return color.NRGBA{}, nil, fmt.Errorf("parse --colors: %w", err)
		}
		palette = append(palette, parsed)
	}
	return background, palette, nil
}

func parseColor(value string) (color.NRGBA, error) {
	if len(value) != 7 || value[0] != '#' {
		return color.NRGBA{}, fmt.Errorf("color %q must use #RRGGBB", value)
	}
	n, err := strconv.ParseUint(value[1:], 16, 24)
	if err != nil {
		return color.NRGBA{}, fmt.Errorf("parse hex color %q: %w", value, err)
	}
	return color.NRGBA{R: uint8(n >> 16), G: uint8(n >> 8), B: uint8(n), A: 255}, nil
}

func render(options Options, background color.NRGBA, palette []color.NRGBA, seed uint64) *image.NRGBA {
	im := image.NewNRGBA(image.Rect(0, 0, options.Width, options.Height))
	draw.Draw(im, im.Bounds(), image.NewUniform(background), image.Point{}, draw.Src)
	random := rand.New(rand.NewPCG(seed, 0))
	size := options.PixelSize
	columns := (options.Width + size - 1) / size
	clustered := clusterCells(options, seed, len(palette))
	radiusSquared := int64(size) * int64(size)
	for top := 0; top < options.Height; top += size {
		for left := 0; left < options.Width; left += size {
			occupied := random.Float64() < options.Density
			foreground := palette[random.IntN(len(palette))]
			if clustered != nil {
				occupied = clustered[(top/size)*columns+left/size]
			}
			if !occupied {
				continue
			}
			for y := top; y < min(top+size, options.Height); y++ {
				for x := left; x < min(left+size, options.Width); x++ {
					if options.Shape == "circle" {
						dx, dy := int64(2*(x-left)+1-size), int64(2*(y-top)+1-size)
						if dx*dx+dy*dy > radiusSquared {
							continue
						}
					}
					im.SetNRGBA(x, y, foreground)
				}
			}
		}
	}
	return im
}

func clusterCells(options Options, seed uint64, paletteSize int) []bool {
	regionalSeeding := options.ClusterAlgorithm == "replication" && options.ReplicationSeedingVariation > 0
	if (options.Clusterization == 0 && !regionalSeeding) || options.Density == 0 || options.Density == 1 {
		return nil
	}
	columns := (options.Width + options.PixelSize - 1) / options.PixelSize
	rows := (options.Height + options.PixelSize - 1) / options.PixelSize
	cellCount := columns * rows
	var field *regionalField
	if options.ClusterAlgorithm == "replication" && (regionalSeeding || options.ReplicationGrowthVariation > 0) {
		field = newRegionalField(options, seed, columns, rows)
	}
	cells := make([]bool, cellCount)
	target := 0
	random := rand.New(rand.NewPCG(seed, 0))
	for index := range cellCount {
		probability := options.Density
		if regionalSeeding {
			probability = field.seedingProbability(index%columns, index/columns)
		}
		if random.Float64() < probability {
			cells[index] = true
			target++
		}
		// Match the renderer's stream, including the palette draw per cell.
		random.IntN(paletteSize)
	}
	if target == 0 || target == cellCount {
		if regionalSeeding {
			return cells
		}
		return nil
	}
	if options.Clusterization == 0 {
		return cells
	}
	random = rand.New(rand.NewPCG(seed, 1))
	switch options.ClusterAlgorithm {
	case "compact":
		compactCells(cells, columns, options.Clusterization, random)
	case "replication":
		replicateCells(cells, columns, options, field, random, rand.New(rand.NewPCG(seed, 2)))
	case "strands":
		clear(cells)
		strandCells(cells, columns, target, options.Clusterization, random)
	}
	return cells
}

func strandCells(cells []bool, columns, target int, factor float64, random *rand.Rand) {
	maxLength := 1 + int(64*factor*factor)
	placed, cursor := 0, 0
	for placed < target {
		start := strandStart(cells, columns, &cursor, random)
		cells[start] = true
		budget := min(1+random.IntN(maxLength), target-placed)
		placed += growStrand(cells, columns, start, budget, random)
	}
}

// strandStart prefers space around a new strand. The fallback cursor visits
// every cell at most once, so dense grids cannot stall on random retries.
func strandStart(cells []bool, columns int, cursor *int, random *rand.Rand) int {
	best, bestNeighbors := -1, 9
	for range 8 {
		index := random.IntN(len(cells))
		if cells[index] {
			continue
		}
		neighbors := occupiedNeighbors(cells, columns, index)
		if neighbors < bestNeighbors {
			best, bestNeighbors = index, neighbors
		}
		if neighbors == 0 {
			return index
		}
	}
	if best >= 0 {
		return best
	}
	for cells[*cursor] {
		*cursor++
	}
	return *cursor
}

type strandPoint struct {
	index, direction int
}

var neighborDirections = [...]image.Point{
	{X: 1, Y: 0},
	{X: 1, Y: 1},
	{X: 0, Y: 1},
	{X: -1, Y: 1},
	{X: -1, Y: 0},
	{X: -1, Y: -1},
	{X: 0, Y: -1},
	{X: 1, Y: -1},
}

func growStrand(cells []bool, columns, start, budget int, random *rand.Rand) int {
	var history [65]strandPoint
	current := strandPoint{index: start, direction: random.IntN(8)}
	history[0] = current
	placed, failedSteps := 1, 0
	for attempt := 0; attempt < 4*budget && placed < budget; attempt++ {
		if placed >= 4 && random.IntN(10) == 0 {
			// Resume from an earlier point with a turned heading to form a branch.
			current = history[random.IntN(placed-2)]
			current.direction = (current.direction + 6 + 4*random.IntN(2)) % 8
		}
		next, ok := strandStep(cells, columns, current, random)
		if !ok {
			failedSteps++
			if failedSteps == 4 {
				break
			}
			current = history[random.IntN(placed)]
			current.direction = random.IntN(8)
			continue
		}
		failedSteps = 0
		cells[next.index] = true
		history[placed] = next
		placed++
		current = next
	}
	return placed
}

func strandStep(cells []bool, columns int, current strandPoint, random *rand.Rand) (strandPoint, bool) {
	x, y := current.index%columns, current.index/columns
	rows := len(cells) / columns
	var weights [8]int
	total := 0
	for direction, offset := range neighborDirections {
		nx, ny := x+offset.X, y+offset.Y
		if nx < 0 || nx >= columns || ny < 0 || ny >= rows || cells[ny*columns+nx] {
			continue
		}
		turn := (direction - current.direction + 8) % 8
		turn = min(turn, 8-turn)
		weight := [...]int{12, 8, 3, 1, 0}[turn]
		neighbors := occupiedNeighbors(cells, columns, ny*columns+nx)
		if neighbors > 2 {
			continue
		}
		if neighbors == 2 {
			weight = (weight + 3) / 4
		}
		weights[direction] = weight
		total += weight
	}
	if total == 0 {
		return strandPoint{}, false
	}
	choice := random.IntN(total)
	for direction, weight := range weights {
		choice -= weight
		if choice < 0 {
			offset := neighborDirections[direction]
			return strandPoint{index: (y+offset.Y)*columns + x + offset.X, direction: direction}, true
		}
	}
	return strandPoint{}, false
}

func occupiedNeighbors(cells []bool, columns, index int) int {
	count := 0
	x, y := index%columns, index/columns
	rows := len(cells) / columns
	for ny := max(0, y-1); ny <= min(rows-1, y+1); ny++ {
		for nx := max(0, x-1); nx <= min(columns-1, x+1); nx++ {
			if (nx != x || ny != y) && cells[ny*columns+nx] {
				count++
			}
		}
	}
	return count
}

func writePNG(path string, im image.Image) (result error) {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o666)
	if err != nil {
		return fmt.Errorf("create output %q: %w", path, err)
	}
	defer func() {
		// Check the open file's identity against the destination before cleanup.
		identity, statErr := f.Stat()
		if err := f.Close(); err != nil {
			result = errors.Join(result, fmt.Errorf("close output %q: %w", path, err))
		}
		if result == nil {
			return
		}
		if statErr != nil {
			result = errors.Join(result, fmt.Errorf("inspect incomplete output %q: %w", path, statErr))
			return
		}
		current, err := os.Lstat(path)
		if errors.Is(err, os.ErrNotExist) {
			return
		}
		if err != nil {
			result = errors.Join(result, fmt.Errorf("inspect output path %q for cleanup: %w", path, err))
			return
		}
		if !os.SameFile(identity, current) {
			result = errors.Join(result, fmt.Errorf("output path %q changed before cleanup", path))
			return
		}
		if err := os.Remove(path); err != nil {
			result = errors.Join(result, fmt.Errorf("remove incomplete output %q: %w", path, err))
		}
	}()
	if err := png.Encode(f, im); err != nil {
		return fmt.Errorf("encode PNG %q: %w", path, err)
	}
	return nil
}
