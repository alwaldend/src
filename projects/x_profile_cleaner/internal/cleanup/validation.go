package cleanup

import (
	"fmt"
	"unicode/utf8"

	"git.alwaldend.com/alwaldend/src/projects/x_profile_cleaner/api/workerpb"
	"git.alwaldend.com/alwaldend/src/projects/x_profile_cleaner/internal/model"
	"git.alwaldend.com/alwaldend/src/projects/x_profile_cleaner/internal/plan"
)

// Validate checks exact cleanup operation/decision agreement before backend publication.
// Archive membership belongs to the source-owning cleanup helper.
func Validate(value *workerpb.Plan) error {
	if err := plan.Validate(value); err != nil {
		return fmt.Errorf("validate worker plan: %w", err)
	}
	if value.CleanupReport == nil {
		return fmt.Errorf("missing cleanup report")
	}
	actions := make(map[string]*workerpb.Operation, len(value.Operations))
	for _, operation := range value.Operations {
		if operation.GetDeletePost() == nil && operation.GetUnrepost() == nil {
			return fmt.Errorf("unsupported cleanup operation")
		}
		key := plan.OperationName(operation) + ":" + plan.TargetID(operation)
		if actions[key] != nil {
			return fmt.Errorf("duplicate cleanup operation")
		}
		actions[key] = operation
	}
	seen := make(map[string]bool, len(value.CleanupReport.Decisions))
	explained := make(map[string]bool, len(value.Operations))
	for _, decision := range value.CleanupReport.Decisions {
		if !model.ValidID(decision.ItemId) || seen[decision.ItemId] || decision.CreatedAt.AsTime().IsZero() || decision.CreatedAt.AsTime().Year() < 1 || decision.CreatedAt.AsTime().Year() > 9999 || !utf8.ValidString(decision.Text) {
			return fmt.Errorf("invalid or duplicate worker decision")
		}
		for _, reason := range decision.Reasons {
			if reason.RuleId <= 0 {
				return fmt.Errorf("invalid cleanup reason identifier")
			}
		}
		seen[decision.ItemId] = true
		kind, target := "delete_post", decision.ItemId
		switch decision.Kind {
		case string(model.ItemPost):
			if decision.SourceId != "" {
				return fmt.Errorf("invalid worker post source")
			}
		case string(model.ItemRepost):
			if !model.ValidID(decision.SourceId) {
				return fmt.Errorf("invalid worker repost source")
			}
			kind, target = "unrepost", decision.SourceId
		default:
			return fmt.Errorf("invalid worker decision kind")
		}
		key := kind + ":" + target
		action, selected := actions[key]
		switch decision.Outcome {
		case "unchanged", "preserved":
			if selected {
				return fmt.Errorf("worker action contradicts unselected decision")
			}
		case "delete", "force_delete":
			if !selected || (action.GetDeletePost().GetForced() || action.GetUnrepost().GetForced()) != (decision.Outcome == "force_delete") || explained[key] {
				return fmt.Errorf("worker selected decision disagrees with action")
			}
			explained[key] = true
		default:
			return fmt.Errorf("invalid worker decision outcome")
		}
	}
	if len(explained) != len(actions) {
		return fmt.Errorf("worker action has no selected decision")
	}
	return nil
}
