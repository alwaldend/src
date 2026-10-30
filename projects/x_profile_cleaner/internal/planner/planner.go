package planner

import (
	"context"
	"fmt"
	"sort"
	"time"

	"git.alwaldend.com/alwaldend/src/projects/x_profile_cleaner/internal/model"
)

// CallbackID is an adapter-owned callback identifier; zero means no callback.
type CallbackID int

// Edge describes a traversal step between owned posts.
type Edge struct {
	// FromID identifies the traversal origin.
	FromID string
	// ToID identifies the traversal destination.
	ToID string
	// Kind identifies the quote or reply relationship.
	Kind model.ReferenceKind
	// Direction is child for forward traversal or parent for a reply-parent step.
	Direction string
}

// CallbackInput contains plain snapshot values for one callback invocation.
type CallbackInput struct {
	// Post is present for selection and inclusion callbacks.
	Post *model.Item
	// Edge is present for edge callbacks.
	Edge *Edge
	// SeedID identifies the preservation seed; selection predicates leave it empty.
	SeedID string
}

// Evaluator executes an adapter-owned callback with plain input and a boolean result.
// Input records borrow snapshot data and must not be mutated by the evaluator.
type Evaluator func(CallbackID, CallbackInput) (bool, error)

// Rule is a language-independent selection declaration.
type Rule struct {
	// Operation names a supported direct selection operation.
	Operation string
	// Name is an optional user-facing description.
	Name string
	// IDs lists direct selections when Predicate is zero.
	IDs []string
	// Predicate identifies a deferred selection callback, or zero for direct IDs.
	Predicate CallbackID
	// Include identifies optional preservation inclusion logic.
	Include CallbackID
	// Follow identifies optional edge traversal logic.
	Follow CallbackID
}

// deletionLevel orders requests so later ordinary deletion cannot weaken a force.
type deletionLevel uint8

const (
	// noDeletion leaves an activity unchanged unless another request selects it.
	noDeletion deletionLevel = iota
	// ordinaryDeletion is prevented by preservation.
	ordinaryDeletion
	// forcedDeletion overrides preservation for the selected activity.
	forcedDeletion
)

// operation declares one cleanup selection's target and effect.
type operation struct {
	// name is the direct Lua function name, without the predicate suffix.
	name string
	// kind selects post or repost activities.
	kind model.ItemKind
	// preserve requests preservation instead of deletion.
	preserve bool
	// deletion is the level requested by a deletion operation.
	deletion deletionLevel
	// quotes extends traversal to future quotes and both directions of replies.
	quotes bool
}

// request describes a preservation seed and the declaration that selected it.
type request struct {
	// id is the canonical seed identifier.
	id string
	// rule indexes the immutable declaration list.
	rule int
}

// declaration adds resolved operation metadata to a registered rule.
type declaration struct {
	// rule holds the copied, language-independent configuration.
	rule Rule
	// operation supplies the planner's semantics for this rule.
	operation operation
}

// selection accumulates independent requests before resolving precedence.
type selection struct {
	// preserved prevents ordinary deletion.
	preserved bool
	// deletion is the strongest requested deletion level.
	deletion deletionLevel
	// reasons retains one match per declaration, keyed by its stable index.
	reasons map[int]model.Reason
}

// Planner owns one immutable snapshot and its registered cleanup rules.
type Planner struct {
	// ctx bounds collection, callbacks, and Go graph traversal.
	ctx context.Context
	// now anchors every age calculation and report timestamp.
	now time.Time
	// account admits only the authenticated author's records.
	account string
	// load supplies a complete view for each requested kind and window.
	load func(context.Context, model.ItemKind, model.TimeWindow) ([]model.Item, error)
	// window is frozen by the first inventory read.
	window model.TimeWindow
	// loaded identifies admitted kind views.
	loaded map[model.ItemKind]bool
	// loadErr invalidates planning even when the adapter catches a query error.
	loadErr error
	// items indexes account-owned activities.
	items map[string]model.Item
	// canonical resolves edit-history aliases.
	canonical map[string]string
	// ids fixes query and evaluation order.
	ids []string
	// children lazily indexes each requested forward relationship.
	children map[model.ReferenceKind]map[string][]Edge
	// rules holds declarations in registration order.
	rules []declaration
	// requests holds direct and predicate-selected preservation seeds.
	requests []request
	// selections accumulates effects by activity ID.
	selections map[string]*selection
	// sealed prevents declarations once evaluation begins.
	sealed bool
}

// New creates a planner without loading inventory or evaluating callbacks.
func New(ctx context.Context, account string, load func(context.Context, model.ItemKind, model.TimeWindow) ([]model.Item, error), now time.Time) *Planner {
	return &Planner{ctx: ctx, account: account, load: load, now: now, loaded: map[model.ItemKind]bool{}, items: map[string]model.Item{}, canonical: map[string]string{}, children: map[model.ReferenceKind]map[string][]Edge{}, selections: map[string]*selection{}}
}

// operations defines each supported operation once for registration and evaluation.
var operations = []operation{
	{name: "delete_post", kind: model.ItemPost, deletion: ordinaryDeletion},
	{name: "force_delete_post", kind: model.ItemPost, deletion: forcedDeletion},
	{name: "delete_repost", kind: model.ItemRepost, deletion: ordinaryDeletion},
	{name: "force_delete_repost", kind: model.ItemRepost, deletion: forcedDeletion},
	{name: "preserve_post", kind: model.ItemPost, preserve: true},
	{name: "preserve_thread", kind: model.ItemPost, preserve: true, quotes: true},
	{name: "preserve_repost", kind: model.ItemRepost, preserve: true},
}

// Operations returns supported direct operation names for language bindings.
func Operations() []string {
	names := make([]string, len(operations))
	for i, op := range operations {
		names[i] = op.name
	}
	return names
}

// Add validates and copies a rule before registering any of its selections.
func (e *Planner) Add(rule Rule) error {
	if e.sealed {
		return fmt.Errorf("cannot register selections during evaluation")
	}
	var op operation
	for _, candidate := range operations {
		if candidate.name == rule.Operation {
			op = candidate
			break
		}
	}
	if op.name == "" {
		return fmt.Errorf("unknown selection operation %q", rule.Operation)
	}
	if (!op.preserve && (rule.Include != 0 || rule.Follow != 0)) || (op.kind == model.ItemRepost && rule.Follow != 0) {
		return fmt.Errorf("unsupported callbacks for %s", op.name)
	}
	if rule.Predicate != 0 && len(rule.IDs) > 0 {
		return fmt.Errorf("predicate and direct selections are mutually exclusive")
	}
	rule.IDs = append([]string(nil), rule.IDs...)
	if rule.Predicate == 0 && len(rule.IDs) > 0 {
		if err := e.ensureKind(op.kind); err != nil {
			return fmt.Errorf("load %s inventory: %w", op.kind, err)
		}
		for i, id := range rule.IDs {
			item, ok := e.items[e.canonical[id]]
			if !ok || item.Kind != op.kind {
				return fmt.Errorf("%s: owned %s %q is absent from the collected inventory", op.name, op.kind, id)
			}
			rule.IDs[i] = item.ID
		}
	}
	index := len(e.rules)
	e.rules = append(e.rules, declaration{rule: rule, operation: op})
	for _, id := range rule.IDs {
		e.selectItem(index, id)
	}
	return nil
}

// SetTimeWindow validates and atomically changes bounds before loading or evaluation.
func (e *Planner) SetTimeWindow(window model.TimeWindow) error {
	if e.sealed || len(e.loaded) > 0 || e.loadErr != nil {
		return fmt.Errorf("set_time_window must precede inventory loading and evaluation")
	}
	if err := window.Validate(); err != nil {
		return fmt.Errorf("validate discovery window: %w", err)
	}
	e.window = window
	return nil
}

// Get returns a collected owned record, resolving edit aliases; false means absent.
func (e *Planner) Get(kind model.ItemKind, id string) (model.Item, bool, error) {
	if err := e.ensureKind(kind); err != nil {
		return model.Item{}, false, fmt.Errorf("load %s: %w", kind, err)
	}
	item, ok := e.items[e.canonical[id]]
	return item, ok && item.Kind == kind, nil
}

// List returns owned records of one kind in stable ID order.
func (e *Planner) List(kind model.ItemKind) ([]model.Item, error) {
	if err := e.ensureKind(kind); err != nil {
		return nil, fmt.Errorf("load %s: %w", kind, err)
	}
	result := []model.Item{}
	for _, id := range e.kindIDs(kind) {
		result = append(result, e.items[id])
	}
	return result, nil
}

// Children returns direct own reply or quote children of an own graph node.
func (e *Planner) Children(kind model.ReferenceKind, id string) ([]model.Item, error) {
	if err := e.ensureGraph(kind); err != nil {
		return nil, fmt.Errorf("load %s graph: %w", kind, err)
	}
	result := []model.Item{}
	for _, child := range e.children[kind][e.canonical[id]] {
		result = append(result, e.items[child.ToID])
	}
	return result, nil
}

// ensureKind loads and validates one kind once, retaining any collection failure.
func (e *Planner) ensureKind(kind model.ItemKind) (err error) {
	if e.loadErr != nil {
		return fmt.Errorf("previous inventory failure: %w", e.loadErr)
	}
	if e.loaded[kind] {
		return nil
	}
	defer func() {
		if err != nil {
			e.loadErr = err
		}
	}()
	items, err := e.load(e.ctx, kind, e.window)
	if err != nil {
		return fmt.Errorf("collect %s: %w", kind, err)
	}
	for _, item := range items {
		if item.AuthorID != e.account || !model.ValidID(item.ID) || item.CreatedAt.IsZero() || item.Kind != kind || !e.window.Contains(item.CreatedAt) {
			return fmt.Errorf("invalid owned %s activity %q", kind, item.ID)
		}
		if _, exists := e.items[item.ID]; exists {
			return fmt.Errorf("duplicate activity %q", item.ID)
		}
		for _, version := range append([]string{item.ID}, item.EditHistoryIDs...) {
			if !model.ValidID(version) {
				return fmt.Errorf("invalid edit-history identifier on activity %s", item.ID)
			}
			if previous, exists := e.canonical[version]; exists && previous != item.ID {
				return fmt.Errorf("edit-history identifier %s belongs to both %s and %s", version, previous, item.ID)
			}
			e.canonical[version] = item.ID
		}
		e.items[item.ID] = item
		e.ids = append(e.ids, item.ID)
		e.selections[item.ID] = &selection{reasons: map[int]model.Reason{}}
	}
	sort.Strings(e.ids)
	e.loaded[kind] = true
	return nil
}

// kindIDs snapshots iteration so callbacks may safely load another kind.
func (e *Planner) kindIDs(kind model.ItemKind) []string {
	ids := []string{}
	for _, id := range e.ids {
		if e.items[id].Kind == kind {
			ids = append(ids, id)
		}
	}
	return ids
}

// ensureGraph builds only the requested direct own-child relationship index.
func (e *Planner) ensureGraph(kind model.ReferenceKind) error {
	if err := e.ensureKind(model.ItemPost); err != nil {
		return fmt.Errorf("load graph posts: %w", err)
	}
	if e.children[kind] != nil {
		return nil
	}
	children := map[string][]Edge{}
	for _, id := range e.kindIDs(model.ItemPost) {
		for _, ref := range e.items[id].References {
			parent, exists := e.items[e.canonical[ref.ID]]
			if ref.Kind == kind && exists && parent.Kind == model.ItemPost {
				children[parent.ID] = append(children[parent.ID], Edge{FromID: parent.ID, ToID: id, Kind: kind, Direction: "child"})
			}
		}
	}
	e.children[kind] = children
	return nil
}

// selectItem accumulates direct effects or registers a preservation seed.
func (e *Planner) selectItem(index int, id string) {
	d := e.rules[index]
	if d.operation.preserve {
		e.requests = append(e.requests, request{id: id, rule: index})
		return
	}
	e.selections[id].deletion = max(e.selections[id].deletion, d.operation.deletion)
	e.recordReason(index, id, "")
}

// recordReason retains the first preservation witness for each matching declaration.
func (e *Planner) recordReason(index int, id, seed string) {
	chosen := e.selections[id]
	if _, exists := chosen.reasons[index]; exists {
		return
	}
	rule := e.rules[index].rule
	name := rule.Operation
	if rule.Predicate != 0 {
		name += "_when"
	}
	chosen.reasons[index] = model.Reason{RuleID: index + 1, Name: rule.Name, Operation: name, SeedID: seed}
}

// Evaluate seals declarations, runs sequential callbacks, and constructs an exact plan.
func (e *Planner) Evaluate(evaluate Evaluator) (model.Plan, error) {
	if e.sealed {
		return model.Plan{}, fmt.Errorf("planner already evaluated")
	}
	e.sealed = true
	for index, d := range e.rules {
		if d.rule.Predicate == 0 {
			continue
		}
		if err := e.ensureKind(d.operation.kind); err != nil {
			return model.Plan{}, fmt.Errorf("load %s_when candidates: %w", d.operation.name, err)
		}
		for _, id := range e.kindIDs(d.operation.kind) {
			if err := e.ctx.Err(); err != nil {
				return model.Plan{}, fmt.Errorf("evaluate policy: %w", err)
			}
			item := e.items[id]
			matches, err := evaluate(d.rule.Predicate, CallbackInput{Post: &item})
			if err != nil {
				return model.Plan{}, fmt.Errorf("%s_when for %s: %w", d.operation.name, id, err)
			}
			if matches {
				e.selectItem(index, id)
			}
		}
	}
	shared := map[int]map[string]bool{}
	for _, r := range e.requests {
		seen := map[string]bool{}
		rule := e.rules[r.rule].rule
		if rule.Include == 0 && rule.Follow == 0 {
			if shared[r.rule] == nil {
				shared[r.rule] = seen
			}
			seen = shared[r.rule]
		}
		if err := e.preserve(r, seen, evaluate); err != nil {
			return model.Plan{}, fmt.Errorf("%s for %s: %w", rule.Operation, r.id, err)
		}
	}
	if e.loadErr != nil {
		return model.Plan{}, fmt.Errorf("load policy inventory: %w", e.loadErr)
	}
	if err := e.ctx.Err(); err != nil {
		return model.Plan{}, fmt.Errorf("finish policy: %w", err)
	}
	plan := model.Plan{Version: 1, AccountID: e.account, EvaluatedAt: e.now.UTC(), Actions: []model.Action{}}
	for _, id := range e.ids {
		item, chosen := e.items[id], e.selections[id]
		outcome := "unchanged"
		if chosen.preserved {
			outcome = "preserved"
		}
		selected := chosen.deletion == forcedDeletion || (chosen.deletion == ordinaryDeletion && !chosen.preserved)
		if selected {
			outcome = "delete"
			if chosen.deletion == forcedDeletion {
				outcome = "force_delete"
			}
			action := model.Action{Kind: model.ActionDeletePost, TargetID: item.ID, Forced: chosen.deletion == forcedDeletion}
			if item.Kind == model.ItemRepost {
				action.Kind = model.ActionUnrepost
				action.TargetID = item.SourceID
			}
			plan.Actions = append(plan.Actions, action)
		}
		reasons := []model.Reason{}
		for _, reason := range chosen.reasons {
			reasons = append(reasons, reason)
		}
		sort.Slice(reasons, func(i, j int) bool { return reasons[i].RuleID < reasons[j].RuleID })
		plan.Decisions = append(plan.Decisions, model.Decision{ItemID: id, SourceID: item.SourceID, Kind: item.Kind, CreatedAt: item.CreatedAt, Text: item.Text, Outcome: outcome, Reasons: reasons})
	}
	sort.Slice(plan.Actions, func(i, j int) bool {
		a, b := plan.Actions[i], plan.Actions[j]
		if a.Kind != b.Kind {
			return a.Kind == model.ActionUnrepost
		}
		return a.TargetID < b.TargetID
	})
	if err := plan.Validate(); err != nil {
		return model.Plan{}, fmt.Errorf("validate generated plan: %w", err)
	}
	return plan, nil
}

// preserve walks a seed with per-request callbacks or shared callback-free visitation.
func (e *Planner) preserve(r request, seen map[string]bool, evaluate Evaluator) error {
	d := e.rules[r.rule]
	if d.operation.kind == model.ItemPost {
		if err := e.ensureGraph(model.ReferenceReply); err != nil {
			return fmt.Errorf("prepare reply graph: %w", err)
		}
		if d.operation.quotes {
			if err := e.ensureGraph(model.ReferenceQuote); err != nil {
				return fmt.Errorf("prepare quote graph: %w", err)
			}
		}
	}
	queue := []string{r.id}
	for len(queue) > 0 {
		if err := e.ctx.Err(); err != nil {
			return fmt.Errorf("traverse preservation graph: %w", err)
		}
		id := queue[0]
		queue = queue[1:]
		if seen[id] {
			continue
		}
		seen[id] = true
		include := true
		if d.rule.Include != 0 {
			item := e.items[id]
			var err error
			include, err = evaluate(d.rule.Include, CallbackInput{Post: &item, SeedID: r.id})
			if err != nil {
				return fmt.Errorf("include post %s: %w", id, err)
			}
		}
		if include {
			e.selections[id].preserved = true
			e.recordReason(r.rule, id, r.id)
		}
		for _, edge := range e.edges(id, d.operation.quotes) {
			follow := true
			if d.rule.Follow != 0 {
				var err error
				follow, err = evaluate(d.rule.Follow, CallbackInput{Edge: &edge, SeedID: r.id})
				if err != nil {
					return fmt.Errorf("follow edge %s to %s: %w", edge.FromID, edge.ToID, err)
				}
			}
			if follow {
				queue = append(queue, edge.ToID)
			}
		}
	}
	return nil
}

// edges returns traversal steps without following quoted ancestors or foreign posts.
// Threads also follow reply parents to preserve the source of a quoting reply.
func (e *Planner) edges(id string, thread bool) []Edge {
	edges := append([]Edge{}, e.children[model.ReferenceReply][id]...)
	if thread {
		edges = append(edges, e.children[model.ReferenceQuote][id]...)
		for _, ref := range e.items[id].References {
			parent, exists := e.items[e.canonical[ref.ID]]
			if ref.Kind == model.ReferenceReply && exists && parent.Kind == model.ItemPost {
				edges = append(edges, Edge{FromID: id, ToID: parent.ID, Kind: model.ReferenceReply, Direction: "parent"})
			}
		}
	}
	sort.Slice(edges, func(i, j int) bool {
		if edges[i].ToID != edges[j].ToID {
			return edges[i].ToID < edges[j].ToID
		}
		return edges[i].Kind < edges[j].Kind
	})
	return edges
}
