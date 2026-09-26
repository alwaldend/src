package generator

import "math"

// regionalField stores only normalization statistics; lattice values are hashed
// on demand, so field memory does not grow with the canvas or shrinkage of scale.
type regionalField struct {
	seed                               uint64
	scale                              int
	density, reproduction              float64
	growthVariation, seedingVariation  float64
	mean, spread, meanWeight, seedGain float64
}

func newRegionalField(options Options, seed uint64, columns, rows int) *regionalField {
	field := &regionalField{
		seed: seed, scale: options.ReplicationFieldScale,
		density: options.Density, reproduction: options.Clusterization,
		growthVariation:  options.ReplicationGrowthVariation,
		seedingVariation: options.ReplicationSeedingVariation,
	}
	minimum, maximum := math.Inf(1), math.Inf(-1)
	minWeight, maxWeight := math.Inf(1), math.Inf(-1)
	for y := range rows {
		for x := range columns {
			value := field.at(x, y)
			weight := regionalSeedWeight(value)
			field.mean += value
			field.meanWeight += weight
			minimum, maximum = min(minimum, value), max(maximum, value)
			minWeight, maxWeight = min(minWeight, weight), max(maxWeight, weight)
		}
	}
	count := float64(columns * rows)
	field.mean /= count
	field.meanWeight /= count
	field.spread = max(field.mean-minimum, maximum-field.mean)
	if field.spread < 1e-12 {
		field.spread = 0
	}
	below, above := field.meanWeight-minWeight, maxWeight-field.meanWeight
	if below > 1e-12 && above > 1e-12 {
		// Centered weights have mean zero. This gain preserves that mean while
		// fitting every probability into [0,1], without clipping its contrast.
		field.seedGain = min(field.density/below, (1-field.density)/above)
	}
	return field
}

func (f *regionalField) seedingProbability(x, y int) float64 {
	if f.seedGain == 0 || f.seedingVariation == 0 {
		return f.density
	}
	contrast := regionalSeedWeight(f.at(x, y)) - f.meanWeight
	return max(0, min(1, f.density+f.seedingVariation*f.seedGain*contrast))
}

func (f *regionalField) growthProbability(x, y int) float64 {
	if f.reproduction == 0 || f.spread == 0 || f.growthVariation == 0 {
		return f.reproduction
	}
	contrast := (f.at(x, y) - f.mean) / f.spread
	return max(0, min(1, f.reproduction+f.growthVariation*contrast))
}

func regionalSeedWeight(value float64) float64 {
	weight := (value + 1) / 2
	weight *= weight
	return weight * weight
}

func (f *regionalField) at(x, y int) float64 {
	return 0.75*regionalNoise(f.seed, x, y, f.scale) +
		0.25*regionalNoise(f.seed^0xd1b54a32d192ed03, 2*x, 2*y, f.scale)
}

func regionalNoise(seed uint64, x, y, scale int) float64 {
	x0, y0 := x/scale, y/scale
	sx, sy := float64(x%scale)/float64(scale), float64(y%scale)/float64(scale)
	sx, sy = sx*sx*(3-2*sx), sy*sy*(3-2*sy)
	a, b := regionalLattice(seed, x0, y0), regionalLattice(seed, x0+1, y0)
	c, d := regionalLattice(seed, x0, y0+1), regionalLattice(seed, x0+1, y0+1)
	top, bottom := a+(b-a)*sx, c+(d-c)*sx
	return top + (bottom-top)*sy
}

func regionalLattice(seed uint64, x, y int) float64 {
	value := (seed + 0x9e3779b97f4a7c15) ^ (uint64(x) * 0x9e3779b185ebca87) ^ (uint64(y) * 0xc2b2ae3d27d4eb4f)
	value = (value ^ (value >> 30)) * 0xbf58476d1ce4e5b9
	value = (value ^ (value >> 27)) * 0x94d049bb133111eb
	value ^= value >> 31
	return 2*float64(value>>11)/float64(uint64(1)<<53) - 1
}
