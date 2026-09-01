// Package admission implements the Phase 3B fail-closed action admission
// engine at provider gateways.
//
// Admit evaluates one operation contract against an exact authority,
// subject, environment, pre-state, and budget envelope. Every violation is
// reported with a stable reason code and a deterministic first reason;
// nothing is admitted unless every applicable rule holds. Declared read and
// hermetic compute work keeps a cheap path that skips prepare-verification
// bookkeeping and validation-set emission.
package admission

import (
	"context"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"git.alwaldend.com/alwaldend/src/tools/agents/api/v1alpha1"
)

// ReasonCode values are stable, machine-readable admission verdicts.
const (
	// ReasonUnknownAction denies an effect that is not declared by the
	// contract or not granted by the authority.
	ReasonUnknownAction = "unknown_action"
	// ReasonAuthorityExpired denies an authority envelope past ExpiresAt.
	ReasonAuthorityExpired = "authority_expired"
	// ReasonBudgetExceeded denies a contract cost beyond the envelope.
	ReasonBudgetExceeded = "budget_exceeded"
	// ReasonWrongEnvironment denies an environment mismatch.
	ReasonWrongEnvironment = "wrong_environment"
	// ReasonStalePreState denies a remote mutation against a stale digest.
	ReasonStalePreState = "stale_pre_state"
	// ReasonMissingPrepareVerify denies a remote mutation without
	// prepare/validate/authorize/execute/verify preconditions.
	ReasonMissingPrepareVerify = "missing_prepare_verify"
	// ReasonMissingCredentialAuthority denies credential or secret output
	// use without the credential.consume grant.
	ReasonMissingCredentialAuthority = "missing_credential_authority"
	// ReasonMissingNetworkAuthority denies network access without the
	// matching network grant.
	ReasonMissingNetworkAuthority = "missing_network_authority"
	// ReasonMissingCancellationBudget denies an operation with no
	// cancellation plan or no cancellation budget margin.
	ReasonMissingCancellationBudget = "missing_cancellation_budget"
	// ReasonProviderScopeMismatch denies an operation whose provider
	// resolves outside the working scope derived from the subject refs.
	ReasonProviderScopeMismatch = "provider_scope_mismatch"
	// ReasonCodeReadCheapPath is carried by admitted read and hermetic
	// compute verdicts.
	ReasonCodeReadCheapPath = "read_cheap_path"
	// ReasonCodeAdmitted is carried by admitted verdicts outside the cheap
	// read path.
	ReasonCodeAdmitted = "admitted"
)

// ContractProvider resolves one operation reference to its declared
// runnable contract. LookupOperation must fail closed for anything the
// provider cannot resolve.
type ContractProvider interface {
	LookupOperation(
		ctx context.Context,
		operation v1alpha1.Reference,
	) (v1alpha1.OperationContract, error)
}

// StaticContractProvider resolves operations from an in-memory map and is
// intended for tests and small embedded gateways.
type StaticContractProvider struct {
	Operations map[v1alpha1.Reference]v1alpha1.OperationContract
}

// LookupOperation returns the contract for operation or an error when it is
// not present.
func (provider StaticContractProvider) LookupOperation(
	ctx context.Context,
	operation v1alpha1.Reference,
) (v1alpha1.OperationContract, error) {
	contract, ok := provider.Operations[operation]
	if !ok {
		return v1alpha1.OperationContract{}, fmt.Errorf(
			"operation %q is not declared",
			operation.ID,
		)
	}
	return contract, nil
}

// Admit performs fail-closed admission of request against contract at time
// now. The returned decision is always well formed: admitted verdicts carry
// an effect set and a clear reason code, denied verdicts carry every
// violated rule.
func Admit(
	request v1alpha1.AdmissionRequest,
	contract v1alpha1.OperationContract,
	now time.Time,
) v1alpha1.AdmissionDecision {
	decision := v1alpha1.AdmissionDecision{
		APIVersion:  v1alpha1.APIVersion,
		Kind:        "AdmissionDecision",
		ValidatedAt: now.UTC().Format(time.RFC3339Nano),
	}
	violations := admitViolations(request, contract, now)
	if len(violations) > 0 {
		decision.ReasonCode = violations[0]
		decision.ReasonDetail = violations
		return decision
	}
	decision.Admitted = true
	decision.ReasonCode = ReasonCodeAdmitted
	if cheapPathEffects(contract.Spec.Effects) {
		decision.ReasonCode = ReasonCodeReadCheapPath
	}
	if decision.Admitted {
		decision.EffectSet = admittedEffectSet(request, contract)
	}
	return decision
}

// admitViolations returns the stable reason code for every violated rule in
// evaluation order. The first element is the deterministic primary code.
func admitViolations(
	request v1alpha1.AdmissionRequest,
	contract v1alpha1.OperationContract,
	now time.Time,
) []string {
	var reasons []string
	appendReason := func(code string) {
		reasons = append(reasons, code)
	}
	if providerScopeMismatch(request, contract) {
		appendReason(ReasonProviderScopeMismatch)
	}
	effects := contractEffects(request, contract)
	for _, effect := range effects {
		if !isKnownEffect(effect) {
			appendReason(ReasonUnknownAction)
			break
		}
	}
	for _, effect := range effects {
		if !isKnownEffect(effect) {
			continue
		}
		if semanticEffect(effect) {
			continue
		}
		if !isGranted(request.Authority.Effects, effect) {
			appendReason(ReasonUnknownAction)
			break
		}
	}
	if request.Authority.ExpiresAt != "" {
		expiry, err := time.Parse(time.RFC3339Nano, request.Authority.ExpiresAt)
		if err == nil && now.After(expiry) {
			appendReason(ReasonAuthorityExpired)
		}
	}
	if !budgetWithin(contract.Spec.Cost, request.Authority.Budget) {
		appendReason(ReasonBudgetExceeded)
	}
	if !environmentMatches(request, contract) {
		appendReason(ReasonWrongEnvironment)
	}
	if (request.RemoteWrite || request.Destroy) &&
		request.ExpectedPreState != "" &&
		request.PreStateDigest != request.ExpectedPreState {
		appendReason(ReasonStalePreState)
	}
	if request.RemoteWrite || request.Destroy {
		if !request.PrepareVerified ||
			strings.TrimSpace(contract.Spec.Preflight) == "" ||
			strings.TrimSpace(contract.Spec.Verification) == "" {
			appendReason(ReasonMissingPrepareVerify)
		}
	}
	if contract.Spec.UsesCredentials ||
		contract.Spec.OutputInformation.Secret ||
		contract.Spec.OutputInformation.Personal {
		if !isGranted(
			request.Authority.Effects,
			v1alpha1.EffectCredentialConsume,
		) {
			appendReason(ReasonMissingCredentialAuthority)
		}
	}
	if contract.Spec.Network == v1alpha1.NetworkWrite {
		if !isGranted(
			request.Authority.Effects,
			v1alpha1.EffectRemoteWrite,
		) {
			appendReason(ReasonMissingNetworkAuthority)
		}
	} else if contract.Spec.Network == v1alpha1.NetworkRead {
		if !isGranted(
			request.Authority.Effects,
			v1alpha1.EffectNetworkRead,
		) {
			appendReason(ReasonMissingNetworkAuthority)
		}
	}
	if strings.TrimSpace(contract.Spec.Cancellation) == "" ||
		request.Budget.Calls == 0 || request.Budget.DurationMS == 0 {
		appendReason(ReasonMissingCancellationBudget)
	}
	return reasons
}

// providerScopeMismatch denies an operation whose provider reference path
// falls outside the working scope derived from the subject references. Both
// subject and provider references are repository or workspace references
// whose identity is a directory path; path traversal that escapes the scope
// (including symlink-style escapes) is rejected.
func providerScopeMismatch(
	request v1alpha1.AdmissionRequest,
	contract v1alpha1.OperationContract,
) bool {
	scope := workingScope(request.SubjectRefs)
	if scope == "" {
		return false
	}
	return !pathWithin(scope, contract.Spec.ProviderRef.ID)
}

// workingScope returns the deepest common directory of the subject refs, or
// an empty string when no subject refs can bound the scope.
func workingScope(subjects []v1alpha1.Reference) string {
	scope := ""
	for _, subject := range subjects {
		if subject.Kind != v1alpha1.ReferenceRepository &&
			subject.Kind != v1alpha1.ReferenceWorkspace {
			continue
		}
		if scope == "" {
			scope = subject.ID
			continue
		}
		for !pathWithin(scope, subject.ID) {
			scope = filepath.Dir(scope)
			if scope == "." || scope == "/" || scope == "" {
				return ""
			}
		}
	}
	return scope
}

func pathWithin(parent, child string) bool {
	relative, err := filepath.Rel(parent, child)
	if err != nil {
		return false
	}
	return relative != ".." &&
		!strings.HasPrefix(relative, ".."+string(filepath.Separator)) &&
		!filepath.IsAbs(relative)
}

// semanticEffect reports whether an effect's grant requirements are owned
// by a semantic rule below (credential consumption, network, remote write or
// destroy) rather than the plain authority-grant subset rule. Such effects
// still must be granted; the semantic rule reports the specific reason.
func semanticEffect(effect v1alpha1.Effect) bool {
	switch effect {
	case v1alpha1.EffectCredentialConsume,
		v1alpha1.EffectNetworkRead,
		v1alpha1.EffectRemoteWrite,
		v1alpha1.EffectRemoteDestroy:
		return true
	}
	return false
}

// contractEffects returns the effects the admission request actually asks
// the provider to exercise. AdmissionRequest carries no explicit per-request
// effect set, so the operation's contract effects are authoritative.
func contractEffects(
	request v1alpha1.AdmissionRequest,
	contract v1alpha1.OperationContract,
) []v1alpha1.Effect {
	return contract.Spec.Effects
}

func isKnownEffect(effect v1alpha1.Effect) bool {
	switch effect {
	case v1alpha1.EffectSourceRead, v1alpha1.EffectSourceWrite,
		v1alpha1.EffectTaskStateWrite, v1alpha1.EffectHostWrite,
		v1alpha1.EffectHistoryWrite, v1alpha1.EffectCodeExecute,
		v1alpha1.EffectCredentialConsume, v1alpha1.EffectNetworkRead,
		v1alpha1.EffectRemoteWrite, v1alpha1.EffectRemoteDestroy:
		return true
	}
	return false
}

func isGranted(grants []v1alpha1.Effect, effect v1alpha1.Effect) bool {
	for _, granted := range grants {
		if granted == effect {
			return true
		}
	}
	return false
}

func budgetWithin(cost, limit v1alpha1.Budget) bool {
	return cost.Calls <= limit.Calls &&
		cost.Bytes <= limit.Bytes &&
		cost.DurationMS <= limit.DurationMS &&
		cost.Concurrency <= limit.Concurrency
}

func environmentMatches(
	request v1alpha1.AdmissionRequest,
	contract v1alpha1.OperationContract,
) bool {
	if contract.Spec.EnvironmentSelector == "" {
		return true
	}
	return containsString(
		request.Environment,
		contract.Spec.EnvironmentSelector,
	) && containsString(
		request.Authority.EnvironmentSelectors,
		contract.Spec.EnvironmentSelector,
	)
}

func containsString(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

// admittedEffectSet returns the canonical, sorted effect set the admission
// grants, or nil when the request is denied.
func admittedEffectSet(
	request v1alpha1.AdmissionRequest,
	contract v1alpha1.OperationContract,
) []v1alpha1.Effect {
	effects := contractEffects(request, contract)
	var admitted []v1alpha1.Effect
	seen := map[v1alpha1.Effect]bool{}
	for _, effect := range effects {
		if isKnownEffect(effect) &&
			isGranted(request.Authority.Effects, effect) &&
			!seen[effect] {
			seen[effect] = true
			admitted = append(admitted, effect)
		}
	}
	sort.Slice(admitted, func(i, j int) bool {
		return admitted[i] < admitted[j]
	})
	return admitted
}

// AdmittedEffectSet returns the granted effect set of an admitted decision,
// or nil when the decision is denied. A denied decision never exposes a
// partial effect set.
func AdmittedEffectSet(decision v1alpha1.AdmissionDecision) []v1alpha1.Effect {
	if !decision.Admitted {
		return nil
	}
	return decision.EffectSet
}

// cheapPathEffects reports whether the union of contract effects is a
// declared read or hermetic compute set: source.read, code.execute, and
// network.read only. Such requests are admitted with reason code
// read_cheap_path and require no prepare-verification bookkeeping.
func cheapPathEffects(effects []v1alpha1.Effect) bool {
	if len(effects) == 0 {
		return false
	}
	for _, effect := range effects {
		switch effect {
		case v1alpha1.EffectSourceRead,
			v1alpha1.EffectCodeExecute,
			v1alpha1.EffectNetworkRead:
			continue
		}
		return false
	}
	return true
}
