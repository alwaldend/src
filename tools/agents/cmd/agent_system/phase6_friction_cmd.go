package agent_system

import (
	"crypto/sha256"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"time"

	v1alpha1 "git.alwaldend.com/alwaldend/src/tools/agents/api/v1alpha1"
)

type frictionOptions struct {
	inputs   multiFlag
	output   string
	markdown string
}

func parseFrictionFlags(args []string) (frictionOptions, error) {
	var opts frictionOptions
	flags := flag.NewFlagSet("friction", flag.ContinueOnError)
	flags.Var(&opts.inputs, "input",
		"FrictionRecord JSON file (repeatable)")
	flags.StringVar(&opts.output, "output", "",
		"bounded aggregate JSON output path")
	flags.StringVar(&opts.markdown, "markdown", "",
		"bounded aggregate Markdown output path")
	if err := flags.Parse(args); err != nil {
		return frictionOptions{}, err
	}
	if len(opts.inputs) == 0 || opts.output == "" {
		return frictionOptions{}, fmt.Errorf(
			"at least one --input and --output are required")
	}
	return opts, nil
}

type frictionAggregate struct {
	APIVersion             string                   `json:"apiVersion"`
	Kind                   string                   `json:"kind"`
	ID                     string                   `json:"id"`
	TotalRecords           int                      `json:"totalRecords"`
	UniqueSignatures       int                      `json:"uniqueSignatures"`
	RecordInputs           []string                 `json:"recordInputs"`
	ObservedAt             string                   `json:"observedAt"`
	TotalAvoidableReads    int                      `json:"totalAvoidableReads"`
	TotalAvoidableCommands int                      `json:"totalAvoidableCommands"`
	TotalLatencyMS         int64                    `json:"totalLatencyMs"`
	Groups                 []frictionSignatureGroup `json:"groups"`
	Truncated              bool                     `json:"truncated"`
	Digest                 string                   `json:"digest,omitempty"`
}

type frictionSignatureGroup struct {
	DefectSignature        string   `json:"defectSignature"`
	RecordCount            int      `json:"recordCount"`
	TotalAvoidableReads    int      `json:"totalAvoidableReads"`
	TotalAvoidableCommands int      `json:"totalAvoidableCommands"`
	TotalLatencyMS         int64    `json:"totalLatencyMs"`
	RecordIDs              []string `json:"recordIds"`
}

func runFriction(args []string, stdout io.Writer) error {
	opts, err := parseFrictionFlags(args)
	if err != nil {
		return err
	}
	records := make([]v1alpha1.FrictionRecord, 0, len(opts.inputs))
	seen := make(map[string]bool, len(opts.inputs))
	for _, input := range opts.inputs {
		content, err := os.ReadFile(input)
		if err != nil {
			return fmt.Errorf("read friction record: %w", err)
		}
		record, err := v1alpha1.DecodeFrictionRecord(content)
		if err != nil {
			return fmt.Errorf("decode friction record %s: %w", input, err)
		}
		if err := record.Validate(); err != nil {
			return fmt.Errorf("validate friction record %s: %w", input, err)
		}
		if seen[record.ID] {
			return fmt.Errorf("duplicate friction record ID %q", record.ID)
		}
		seen[record.ID] = true
		records = append(records, record)
	}
	aggregate := buildFrictionAggregate(records, opts.inputs)
	content, err := json.MarshalIndent(aggregate, "", "  ")
	if err != nil {
		return fmt.Errorf("encode friction aggregate: %w", err)
	}
	content = append(content, '\n')
	if err := os.WriteFile(opts.output, content, 0o644); err != nil {
		return fmt.Errorf("write friction aggregate: %w", err)
	}
	if opts.markdown != "" {
		markdown := renderFrictionMarkdown(aggregate)
		if err := os.WriteFile(opts.markdown, []byte(markdown), 0o644); err != nil {
			return fmt.Errorf("write friction markdown: %w", err)
		}
	}
	return writeJSONLine(stdout, map[string]any{
		"accepted":   len(records),
		"signatures": aggregate.UniqueSignatures,
		"output":     opts.output,
	})
}

func buildFrictionAggregate(
	records []v1alpha1.FrictionRecord,
	inputPaths []string,
) *frictionAggregate {
	groups := make(map[string]*frictionSignatureGroup)
	aggregate := &frictionAggregate{
		APIVersion:   "agents.alwaldend.com/v1alpha1",
		Kind:         "FrictionAggregate",
		ID:           "agent-system.friction-aggregate",
		TotalRecords: len(records),
	}
	now := time.Now().UTC().Format(time.RFC3339Nano)
	aggregate.ObservedAt = now
	for _, record := range records {
		group, ok := groups[record.DefectSignature]
		if !ok {
			group = &frictionSignatureGroup{
				DefectSignature: record.DefectSignature,
			}
			groups[record.DefectSignature] = group
		}
		group.RecordCount++
		readSum, ok := checkedAddInt(group.TotalAvoidableReads,
			record.AvoidableReads)
		if ok {
			group.TotalAvoidableReads = readSum
		}
		commandSum, ok := checkedAddInt(group.TotalAvoidableCommands,
			record.AvoidableCommands)
		if ok {
			group.TotalAvoidableCommands = commandSum
		}
		latencySum, ok := checkedAddInt64(group.TotalLatencyMS,
			record.LatencyMS)
		if ok {
			group.TotalLatencyMS = latencySum
		}
		group.RecordIDs = append(group.RecordIDs, record.ID)
		readSum, ok = checkedAddInt(aggregate.TotalAvoidableReads,
			record.AvoidableReads)
		if ok {
			aggregate.TotalAvoidableReads = readSum
		}
		commandSum, ok = checkedAddInt(aggregate.TotalAvoidableCommands,
			record.AvoidableCommands)
		if ok {
			aggregate.TotalAvoidableCommands = commandSum
		}
		latencySum, ok = checkedAddInt64(aggregate.TotalLatencyMS,
			record.LatencyMS)
		if ok {
			aggregate.TotalLatencyMS = latencySum
		}
	}
	aggregate.UniqueSignatures = len(groups)
	for _, group := range groups {
		aggregate.Groups = append(aggregate.Groups, *group)
	}
	sort.Slice(aggregate.Groups, func(left, right int) bool {
		leftCost := aggregate.Groups[left].TotalAvoidableReads +
			aggregate.Groups[left].TotalAvoidableCommands
		rightCost := aggregate.Groups[right].TotalAvoidableReads +
			aggregate.Groups[right].TotalAvoidableCommands
		if leftCost != rightCost {
			return leftCost > rightCost
		}
		return aggregate.Groups[left].DefectSignature <
			aggregate.Groups[right].DefectSignature
	})
	aggregate.RecordInputs = append([]string(nil), inputPaths...)
	aggregate.Digest = frictionAggregateDigest(aggregate)
	return aggregate
}

func frictionAggregateDigest(aggregate *frictionAggregate) string {
	binding := struct {
		TotalRecords           int                      `json:"totalRecords"`
		UniqueSignatures       int                      `json:"uniqueSignatures"`
		TotalAvoidableReads    int                      `json:"totalAvoidableReads"`
		TotalAvoidableCommands int                      `json:"totalAvoidableCommands"`
		TotalLatencyMS         int64                    `json:"totalLatencyMs"`
		Groups                 []frictionSignatureGroup `json:"groups"`
		RecordInputs           []string                 `json:"recordInputs"`
	}{
		TotalRecords:           aggregate.TotalRecords,
		UniqueSignatures:       aggregate.UniqueSignatures,
		TotalAvoidableReads:    aggregate.TotalAvoidableReads,
		TotalAvoidableCommands: aggregate.TotalAvoidableCommands,
		TotalLatencyMS:         aggregate.TotalLatencyMS,
		Groups:                 aggregate.Groups,
		RecordInputs:           aggregate.RecordInputs,
	}
	content, err := json.Marshal(binding)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("sha256:%x", sha256.Sum256(content))
}

func checkedAddInt(left, right int) (int, bool) {
	sum := left + right
	if (left > 0 && right > 0 && sum < 0) ||
		(left < 0 && right < 0 && sum >= 0) {
		return 0, false
	}
	return sum, true
}

func checkedAddInt64(left, right int64) (int64, bool) {
	sum := left + right
	if (left > 0 && right > 0 && sum < 0) ||
		(left < 0 && right < 0 && sum >= 0) {
		return 0, false
	}
	return sum, true
}

func renderFrictionMarkdown(aggregate *frictionAggregate) string {
	result := "# Friction aggregate\n\n"
	result += fmt.Sprintf("- Records: %d\n", aggregate.TotalRecords)
	result += fmt.Sprintf("- Unique defect signatures: %d\n",
		aggregate.UniqueSignatures)
	result += fmt.Sprintf("- Total avoidable reads: %d\n",
		aggregate.TotalAvoidableReads)
	result += fmt.Sprintf("- Total avoidable commands: %d\n",
		aggregate.TotalAvoidableCommands)
	result += fmt.Sprintf("- Total latency: %dms\n", aggregate.TotalLatencyMS)
	result += "\n| Defect signature | Records | Avoidable reads | Avoidable commands | Latency (ms) |\n"
	result += "|---|---:|---:|---:|---:|\n"
	for _, group := range aggregate.Groups {
		result += fmt.Sprintf(
			"| `%s` | %d | %d | %d | %d |\n",
			group.DefectSignature,
			group.RecordCount,
			group.TotalAvoidableReads,
			group.TotalAvoidableCommands,
			group.TotalLatencyMS,
		)
	}
	return result
}
