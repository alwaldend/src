package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"git.alwaldend.com/alwaldend/src/projects/icon_generator/internal/generator"
)

func main() {
	if err := run(os.Args[1:], os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, "icon_generator:", err)
		os.Exit(1)
	}
}

func run(args []string, stdout, stderr io.Writer) error {
	var options generator.Options
	var seed uint64
	var help bytes.Buffer
	flags := flag.NewFlagSet("icon_generator", flag.ContinueOnError)
	flags.SetOutput(&help)
	flags.IntVar(&options.Width, "width", 256, "canvas width in output pixels (1..4096)")
	flags.IntVar(&options.Height, "height", 256, "canvas height in output pixels (1..4096)")
	flags.StringVar(&options.Background, "background", "#ffffff", "background #RRGGBB or transparent")
	flags.StringVar(&options.Colors, "colors", "#000000", "comma-separated #RRGGBB palette")
	flags.IntVar(&options.PixelSize, "pixel-size", 16, "grid cell size: square side or circle diameter (1..4096)")
	flags.StringVar(&options.Shape, "pixel-shape", "square", "shape: square or circle")
	flags.Float64Var(&options.Density, "density", 0.4, "average initial cell occupancy probability (0..1)")
	flags.Float64Var(&options.Clusterization, "clusterization", 0, "clustering strength, or reproduction probability in replication mode (0..1)")
	flags.StringVar(&options.ClusterAlgorithm, "clusterization-algorithm", "compact", "algorithm: compact, strands, or replication")
	flags.Float64Var(&options.ReplicationInheritance, "replication-inheritance", 0, "probability of retaining parent-to-child direction; used with forward bias (0..1)")
	flags.Float64Var(&options.ReplicationForwardBias, "replication-forward-bias", 0, "prefer forward movement and gentle turns in replication (0..1)")
	flags.Float64Var(&options.ReplicationCrowding, "replication-crowding", 0, "prefer replication destinations with fewer occupied neighbors (0..1)")
	flags.Float64Var(&options.ReplicationBranching, "replication-branching", 0, "probability the second child favors a sideways fork in replication (0..1)")
	flags.IntVar(&options.ReplicationFieldScale, "replication-field-scale", 128, "approximate regional scale in logical cells (1..4096)")
	flags.Float64Var(&options.ReplicationGrowthVariation, "replication-growth-variation", 0, "maximum regional offset to reproduction probability, clamped to 0..1 (0..1)")
	flags.Float64Var(&options.ReplicationSeedingVariation, "replication-seeding-variation", 0, "regional initial-population contrast, preserving average density (0..1)")
	flags.Uint64Var(&seed, "seed", 0, "reproducible unsigned seed; omitted means a fresh seed, zero is valid")
	flags.StringVar(&options.Output, "output", "icon.png", "PNG output path; existing files are preserved")
	flags.Usage = func() {
		fmt.Fprintln(&help, "Usage: icon_generator [options]")
		fmt.Fprintln(&help, "Generate a random grid of colored squares or circles as a PNG.")
		flags.PrintDefaults()
		fmt.Fprintln(&help, "\nWidth and height are output pixels. Logical pixels occupy pixel-size grid cells.")
		fmt.Fprintln(&help, "Density is a probability per cell, not an exact count. Edge shapes are clipped.")
		fmt.Fprintln(&help, "Clustering uses all eight neighbors; zero keeps independent placement.")
		fmt.Fprintln(&help, "Compact and strands preserve cell count. Replication can grow beyond the initial density.")
		fmt.Fprintln(&help, "Each replication cell adds two adjacent children or stops; parents remain visible.")
		fmt.Fprintln(&help, "Replication placement and variation strengths default to zero and affect replication only.")
		fmt.Fprintln(&help, "Regional seeding keeps average initial density and can vary placement even with zero reproduction.")
		fmt.Fprintln(&help, "Replay the reported seed with the same settings to reproduce the image.")
	}
	if err := flags.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			if _, err := io.Copy(stdout, &help); err != nil {
				return fmt.Errorf("write help: %w", err)
			}
			return nil
		}
		return fmt.Errorf("parse arguments: %w", err)
	}
	if flags.NArg() != 0 {
		return fmt.Errorf("unexpected positional arguments; use --help for usage")
	}
	flags.Visit(func(f *flag.Flag) {
		if f.Name == "seed" {
			options.Seed = &seed
		}
	})
	effective, err := generator.Generate(options)
	if err != nil {
		return fmt.Errorf("generate icon: %w", err)
	}
	if _, err := fmt.Fprintf(stderr, "seed: %d\n", effective); err != nil {
		return fmt.Errorf("report seed: %w", err)
	}
	return nil
}
