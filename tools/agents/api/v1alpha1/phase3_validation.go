package v1alpha1

import (
	"fmt"
	"sort"
	"strings"
)

func oneOfImpactProfile(value ImpactProfile) bool {
	switch value {
	case ImpactProfileChangedFast, ImpactProfileWorkspace,
		ImpactProfileFreshEvidence, ImpactProfileFullAudit,
		ImpactProfileDiagnose:
		return true
	}
	return false
}

func oneOfCostClass(value CostClass) bool {
	switch value {
	case CostClassLow, CostClassMedium, CostClassHigh:
		return true
	}
	return false
}

func (target PlanTarget) Validate() error {
	if strings.TrimSpace(target.Label) == "" {
		return fmt.Errorf("plan target label is required")
	}
	if target.Path != "" && target.AffectedBy == "" {
		return fmt.Errorf("plan target %q has a path but no affected-by binding", target.Label)
	}
	return nil
}

func (check PlanCheck) Validate() error {
	if !identifierPattern.MatchString(check.Identifier) {
		return fmt.Errorf("malformed plan check identifier %q", check.Identifier)
	}
	return nil
}

func (plan ImpactPlan) Validate() error {
	if plan.APIVersion != APIVersion || plan.Kind != "ImpactPlan" {
		return fmt.Errorf(
			"unsupported impact plan type %q %q",
			plan.APIVersion,
			plan.Kind,
		)
	}
	if !identifierPattern.MatchString(plan.ID) {
		return fmt.Errorf("malformed impact plan identity %q", plan.ID)
	}
	if err := plan.IntentRef.Validate(); err != nil {
		return fmt.Errorf("intent reference: %w", err)
	}
	if !oneOfImpactProfile(plan.Profile) {
		return fmt.Errorf("unknown impact profile %q", plan.Profile)
	}
	if len(plan.Capabilities) == 0 {
		return fmt.Errorf("impact plan must select at least one capability")
	}
	seen := map[Reference]bool{}
	for _, reference := range plan.Capabilities {
		if err := reference.Validate(); err != nil {
			return fmt.Errorf("capability reference: %w", err)
		}
		if reference.Kind != ReferenceProvider &&
			reference.Kind != ReferenceArtifact {
			return fmt.Errorf(
				"capability %q must be a provider or artifact reference",
				reference.ID,
			)
		}
		if seen[reference] {
			return fmt.Errorf("duplicate capability %q", reference.ID)
		}
		seen[reference] = true
	}
	if len(plan.Effects) == 0 {
		return fmt.Errorf("impact plan must declare at least one effect")
	}
	effectSeen := map[Effect]bool{}
	for _, effect := range plan.Effects {
		if !knownEffect(effect) {
			return fmt.Errorf("unknown effect %q", effect)
		}
		if effectSeen[effect] {
			return fmt.Errorf("duplicate effect %q", effect)
		}
		effectSeen[effect] = true
	}
	if !sort.SliceIsSorted(plan.Effects, func(i, j int) bool {
		return plan.Effects[i] < plan.Effects[j]
	}) {
		return fmt.Errorf("effects must use stable lexical order")
	}
	forbidden := map[Effect]bool{}
	for _, effect := range plan.ForbiddenEffects {
		if !knownEffect(effect) {
			return fmt.Errorf("unknown forbidden effect %q", effect)
		}
		if forbidden[effect] {
			return fmt.Errorf("duplicate forbidden effect %q", effect)
		}
		if effectSeen[effect] {
			return fmt.Errorf("effect %q is both required and forbidden", effect)
		}
		forbidden[effect] = true
	}
	if len(plan.Targets) == 0 {
		return fmt.Errorf("impact plan must declare at least one target")
	}
	for _, target := range plan.Targets {
		if err := target.Validate(); err != nil {
			return err
		}
	}
	if len(plan.Checks) == 0 {
		return fmt.Errorf("impact plan must declare at least one check")
	}
	for _, check := range plan.Checks {
		if err := check.Validate(); err != nil {
			return err
		}
	}
	if !oneOfCostClass(plan.CostClass) {
		return fmt.Errorf("unknown cost class %q", plan.CostClass)
	}
	if err := plan.ExpectedMaxCost.Validate(); err != nil {
		return fmt.Errorf("expected max cost: %w", err)
	}
	if !oneOf(
		string(plan.Cacheability),
		string(Cacheable),
		string(NotCacheable),
		string(CacheabilityUnknown),
	) {
		return fmt.Errorf("unknown cacheability %q", plan.Cacheability)
	}
	if plan.Digest != "" && !validDigest(plan.Digest) {
		return fmt.Errorf("malformed impact plan digest")
	}
	return nil
}

func (inputs ValidationInputs) Validate() error {
	for name, digest := range inputs.ConfigDigests {
		if !validDigest(digest) {
			return fmt.Errorf("malformed config digest for %q", name)
		}
	}
	for name, digest := range inputs.ToolchainDigests {
		if !validDigest(digest) {
			return fmt.Errorf("malformed toolchain digest for %q", name)
		}
	}
	for name, digest := range inputs.PolicyDigests {
		if !validDigest(digest) {
			return fmt.Errorf("malformed policy digest for %q", name)
		}
	}
	for _, reference := range inputs.ProviderRefs {
		if err := reference.Validate(); err != nil {
			return fmt.Errorf("provider reference: %w", err)
		}
	}
	return nil
}

func (result ValidationResult) Validate() error {
	if !identifierPattern.MatchString(result.CheckID) {
		return fmt.Errorf("malformed validation check identifier %q", result.CheckID)
	}
	if !oneOf(result.Status, "pass", "fail", "skipped", "blocked") {
		return fmt.Errorf("unknown validation result status %q", result.Status)
	}
	if result.DurationMS < 0 {
		return fmt.Errorf("validation duration cannot be negative")
	}
	return nil
}

func (set ValidationSet) Validate() error {
	if set.APIVersion != APIVersion || set.Kind != "ValidationSet" {
		return fmt.Errorf(
			"unsupported validation set type %q %q",
			set.APIVersion,
			set.Kind,
		)
	}
	if !identifierPattern.MatchString(set.ID) {
		return fmt.Errorf("malformed validation set identity %q", set.ID)
	}
	if !oneOfImpactProfile(set.Profile) {
		return fmt.Errorf("unknown validation profile %q", set.Profile)
	}
	if err := set.Candidate.Validate(); err != nil {
		return fmt.Errorf("candidate reference: %w", err)
	}
	if set.Candidate.Kind != ReferenceSubject &&
		set.Candidate.Kind != ReferenceRepository {
		return fmt.Errorf("validation candidate must be a subject or repository reference")
	}
	if err := set.Inputs.Validate(); err != nil {
		return fmt.Errorf("validation inputs: %w", err)
	}
	for _, argument := range set.SanitizedArgs {
		if strings.Contains(argument, "@") && strings.Contains(argument, ":") {
			return fmt.Errorf("sanitized argument %q resembles a credential", argument)
		}
	}
	if strings.TrimSpace(set.WorkingScope) == "" {
		return fmt.Errorf("validation set working scope is required")
	}
	if len(set.Results) == 0 {
		return fmt.Errorf("validation set must contain at least one result")
	}
	seen := map[string]bool{}
	for _, result := range set.Results {
		if err := result.Validate(); err != nil {
			return err
		}
		if seen[result.CheckID] {
			return fmt.Errorf("duplicate validation result for check %q", result.CheckID)
		}
		seen[result.CheckID] = true
	}
	if set.TotalDurationMS < 0 {
		return fmt.Errorf("validation total duration cannot be negative")
	}
	if err := set.OutputBounds.Validate(); err != nil {
		return fmt.Errorf("output bounds: %w", err)
	}
	for _, reference := range set.RawLogRefs {
		if err := reference.Validate(); err != nil {
			return fmt.Errorf("raw-log reference: %w", err)
		}
		if reference.Kind != ReferenceArtifact {
			return fmt.Errorf("raw-log reference must be an artifact reference")
		}
	}
	if set.Digest != "" && !validDigest(set.Digest) {
		return fmt.Errorf("malformed validation set digest")
	}
	return nil
}

func (assertion EvidenceAssertion) Validate() error {
	if assertion.APIVersion != APIVersion || assertion.Kind != "EvidenceAssertion" {
		return fmt.Errorf(
			"unsupported evidence assertion type %q %q",
			assertion.APIVersion,
			assertion.Kind,
		)
	}
	if !identifierPattern.MatchString(assertion.ID) {
		return fmt.Errorf("malformed evidence assertion identity %q", assertion.ID)
	}
	if err := assertion.CriterionRef.Validate(); err != nil {
		return fmt.Errorf("criterion reference: %w", err)
	}
	if assertion.CriterionRef.Kind != ReferenceGoal &&
		assertion.CriterionRef.Kind != ReferenceArtifact {
		return fmt.Errorf("evidence assertion criterion must be a goal or artifact reference")
	}
	if err := assertion.CriterionRef.Validate(); err != nil {
		return fmt.Errorf("criterion reference: %w", err)
	}
	if strings.TrimSpace(assertion.CriterionRev) == "" {
		return fmt.Errorf("evidence assertion criterion revision is required")
	}
	if len(assertion.ValidationRefs) == 0 {
		return fmt.Errorf("evidence assertion must reference at least one validation set")
	}
	seen := map[Reference]bool{}
	for _, reference := range assertion.ValidationRefs {
		if err := reference.Validate(); err != nil {
			return fmt.Errorf("validation reference: %w", err)
		}
		if reference.Kind != ReferenceArtifact {
			return fmt.Errorf("validation reference must be an artifact reference")
		}
		if seen[reference] {
			return fmt.Errorf("duplicate validation reference %q", reference.ID)
		}
		seen[reference] = true
	}
	if !oneOf(assertion.Verdict, "satisfied", "unsatisfied", "blocked", "unknown") {
		return fmt.Errorf("unknown evidence verdict %q", assertion.Verdict)
	}
	if assertion.Digest != "" && !validDigest(assertion.Digest) {
		return fmt.Errorf("malformed evidence assertion digest")
	}
	return nil
}

func (request AdmissionRequest) Validate() error {
	if request.APIVersion != APIVersion || request.Kind != "AdmissionRequest" {
		return fmt.Errorf(
			"unsupported admission request type %q %q",
			request.APIVersion,
			request.Kind,
		)
	}
	if err := request.Operation.Validate(); err != nil {
		return fmt.Errorf("operation reference: %w", err)
	}
	if request.Operation.Kind != ReferenceOperation {
		return fmt.Errorf("admission operation must be an operation reference")
	}
	if err := request.Authority.Validate(); err != nil {
		return fmt.Errorf("authority envelope: %w", err)
	}
	for _, reference := range request.SubjectRefs {
		if err := reference.Validate(); err != nil {
			return fmt.Errorf("subject reference: %w", err)
		}
	}
	if err := request.Budget.Validate(); err != nil {
		return fmt.Errorf("budget: %w", err)
	}
	if request.RemoteWrite || request.Destroy {
		if !request.PrepareVerified {
			return fmt.Errorf("remote write or destroy requires prepare verification")
		}
	}
	return nil
}

func (decision AdmissionDecision) Validate() error {
	if decision.APIVersion != APIVersion || decision.Kind != "AdmissionDecision" {
		return fmt.Errorf(
			"unsupported admission decision type %q %q",
			decision.APIVersion,
			decision.Kind,
		)
	}
	if !decision.Admitted && decision.ReasonCode == "" {
		return fmt.Errorf("a denied admission must carry a reason code")
	}
	for _, effect := range decision.EffectSet {
		if !knownEffect(effect) {
			return fmt.Errorf("unknown decision effect %q", effect)
		}
	}
	return nil
}

func (authority AuthorityEnvelope) Validate() error {
	if err := authority.ActorRef.Validate(); err != nil {
		return fmt.Errorf("actor reference: %w", err)
	}
	for _, reference := range authority.SubjectRefs {
		if err := reference.Validate(); err != nil {
			return fmt.Errorf("subject reference: %w", err)
		}
	}
	if len(authority.Effects) == 0 {
		return fmt.Errorf("authority must grant at least one effect")
	}
	seen := map[Effect]bool{}
	for _, effect := range authority.Effects {
		if !knownEffect(effect) {
			return fmt.Errorf("authority grants unknown effect %q", effect)
		}
		if seen[effect] {
			return fmt.Errorf("authority grants duplicate effect %q", effect)
		}
		seen[effect] = true
	}
	return authority.Budget.Validate()
}
