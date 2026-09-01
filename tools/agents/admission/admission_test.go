package admission

import (
	"context"
	"reflect"
	"strings"
	"testing"
	"time"

	"git.alwaldend.com/alwaldend/src/tools/agents/api/v1alpha1"
)

// faultContract returns a minimally complete operation contract for a
// workspace-scoped provider. Tests mutate the copy they need.
func faultContract() v1alpha1.OperationContract {
	return v1alpha1.OperationContract{
		APIVersion: v1alpha1.APIVersion,
		Kind:       v1alpha1.KindOperationContract,
		Metadata: v1alpha1.OperationMetadata{
			Name:     "op/test",
			OwnerRef: validRef(v1alpha1.ReferenceWorkspace, "work/ws"),
		},
		Spec: v1alpha1.OperationSpec{
			ProviderRef: validRef(v1alpha1.ReferenceProvider, "work/ws/provider"),
			Effects: []v1alpha1.Effect{
				v1alpha1.EffectSourceRead,
				v1alpha1.EffectCodeExecute,
			},
			InputInformation:  v1alpha1.InformationSet{Public: true},
			OutputInformation: v1alpha1.InformationSet{Public: true},
			Network:           v1alpha1.NetworkNone,
			AuthorityGate:     "authority_gate",
			Preflight:         "preflight",
			Verification:      "verification",
			Cost: v1alpha1.Budget{
				Calls:      10,
				Bytes:      1 << 20,
				DurationMS: 30000,
			},
			Cacheability: v1alpha1.Cacheable,
			Cancellation: "cancel",
		},
	}
}

func validRef(kind v1alpha1.ReferenceKind, id string) v1alpha1.Reference {
	return v1alpha1.Reference{Kind: kind, ID: id, Version: "v1"}
}

// baseRequest returns an admitted-eligible request for faultContract with a
// read/hermetic effect profile and a real cancellation budget.
func baseRequest() v1alpha1.AdmissionRequest {
	return v1alpha1.AdmissionRequest{
		APIVersion: v1alpha1.APIVersion,
		Kind:       "AdmissionRequest",
		Operation:  validRef(v1alpha1.ReferenceOperation, "op/test"),
		Authority: v1alpha1.AuthorityEnvelope{
			ActorRef: validRef(v1alpha1.ReferenceActor, "actor/alice"),
			SubjectRefs: []v1alpha1.Reference{
				validRef(v1alpha1.ReferenceSubject, "work/ws"),
			},
			Effects: []v1alpha1.Effect{
				v1alpha1.EffectSourceRead,
				v1alpha1.EffectCodeExecute,
				v1alpha1.EffectNetworkRead,
				v1alpha1.EffectCredentialConsume,
				v1alpha1.EffectRemoteWrite,
			},
			Budget: v1alpha1.Budget{
				Calls:      100,
				Bytes:      1 << 22,
				DurationMS: 120000,
			},
		},
		SubjectRefs: []v1alpha1.Reference{
			validRef(v1alpha1.ReferenceWorkspace, "work/ws"),
		},
		Environment: []string{"sandbox"},
		Budget: v1alpha1.Budget{
			Calls:      10,
			DurationMS: 60000,
		},
	}
}

func nowValue() time.Time {
	return time.Date(2026, 9, 1, 12, 0, 0, 0, time.UTC)
}

func denyReason(
	t *testing.T,
	code string,
	decision v1alpha1.AdmissionDecision,
) {
	t.Helper()
	if decision.Admitted {
		t.Fatalf(
			"expected denial %q for contract %q, got admission with %v",
			code,
			decision.ReasonCode,
			decision.EffectSet,
		)
	}
	if decision.ReasonCode != code {
		t.Fatalf(
			"expected reason code %q, got %q with detail %v",
			code,
			decision.ReasonCode,
			decision.ReasonDetail,
		)
	}
}

// TestWrongEnvironmentDeny asserts wrong_environment when the request
// environment misses the required contract selector.
func TestWrongEnvironmentDeny(t *testing.T) {
	contract := faultContract()
	contract.Spec.EnvironmentSelector = "sandbox"
	request := baseRequest()
	request.Environment = []string{"prod"}
	decision := Admit(request, contract, nowValue())
	denyReason(t, ReasonWrongEnvironment, decision)
}

// TestStaleAuthorityDeny asserts authority_expired when the envelope is past
// ExpiresAt.
func TestStaleAuthorityDeny(t *testing.T) {
	request := baseRequest()
	request.Authority.ExpiresAt = "2026-09-01T11:59:59Z"
	decision := Admit(request, faultContract(), nowValue())
	denyReason(t, ReasonAuthorityExpired, decision)
}

// TestSecretOutputWithoutAuthorityDeny asserts missing_credential_authority
// when the contract's secret output has no credential.consume grant.
func TestSecretOutputWithoutAuthorityDeny(t *testing.T) {
	contract := faultContract()
	contract.Spec.OutputInformation = v1alpha1.InformationSet{Secret: true}
	request := baseRequest()
	request.Authority.Effects = []v1alpha1.Effect{
		v1alpha1.EffectSourceRead,
		v1alpha1.EffectCodeExecute,
	}
	decision := Admit(request, contract, nowValue())
	denyReason(t, ReasonMissingCredentialAuthority, decision)
}

// TestContractOutputPersonalWithoutAuthorityDeny asserts the same gate for
// personal information output.
func TestContractOutputPersonalWithoutAuthorityDeny(t *testing.T) {
	contract := faultContract()
	contract.Spec.OutputInformation = v1alpha1.InformationSet{Personal: true}
	request := baseRequest()
	request.Authority.Effects = []v1alpha1.Effect{
		v1alpha1.EffectSourceRead,
		v1alpha1.EffectCodeExecute,
	}
	decision := Admit(request, contract, nowValue())
	denyReason(t, ReasonMissingCredentialAuthority, decision)
}

// TestUnknownActionDeny asserts unknown_action for an effect the authority
// never granted.
func TestUnknownActionDeny(t *testing.T) {
	contract := faultContract()
	contract.Spec.Effects = []v1alpha1.Effect{v1alpha1.EffectSourceWrite}
	request := baseRequest()
	request.Authority.Effects = []v1alpha1.Effect{
		v1alpha1.EffectSourceRead,
		v1alpha1.EffectCodeExecute,
	}
	decision := Admit(request, contract, nowValue())
	denyReason(t, ReasonUnknownAction, decision)
}

// TestStalePreStateDeny asserts stale_pre_state for a remote write whose
// pre-state digest no longer matches the expected digest.
func TestStalePreStateDeny(t *testing.T) {
	contract := faultContract()
	contract.Spec.Effects = []v1alpha1.Effect{v1alpha1.EffectRemoteWrite}
	contract.Spec.Cost.Calls = 100
	request := baseRequest()
	request.RemoteWrite = true
	request.PrepareVerified = true
	request.ExpectedPreState = "digest-expected"
	request.PreStateDigest = "digest-stale"
	decision := Admit(request, contract, nowValue())
	denyReason(t, ReasonStalePreState, decision)
}

// TestRemoteWriteWithoutPrepareVerifyDeny asserts missing_prepare_verify for
// a remote write whose prepare verification flag is unset.
func TestRemoteWriteWithoutPrepareVerifyDeny(t *testing.T) {
	contract := faultContract()
	contract.Spec.Effects = []v1alpha1.Effect{v1alpha1.EffectRemoteWrite}
	contract.Spec.Cost.Calls = 100
	request := baseRequest()
	request.RemoteWrite = true
	request.ExpectedPreState = "digest-expected"
	request.PreStateDigest = "digest-expected"
	request.PrepareVerified = false
	decision := Admit(request, contract, nowValue())
	denyReason(t, ReasonMissingPrepareVerify, decision)
}

// TestCleanupMissingDeny asserts missing_prepare_verify when a write
// contract has no preflight or verification step at all.
func TestCleanupMissingDeny(t *testing.T) {
	contract := faultContract()
	contract.Spec.Effects = []v1alpha1.Effect{
		v1alpha1.EffectHostWrite,
		v1alpha1.EffectCodeExecute,
	}
	contract.Spec.Preflight = ""
	contract.Spec.Verification = ""
	request := baseRequest()
	request.RemoteWrite = true
	request.PrepareVerified = true
	request.ExpectedPreState = "digest-expected"
	request.PreStateDigest = "digest-expected"
	request.Authority.Effects = []v1alpha1.Effect{
		v1alpha1.EffectHostWrite,
		v1alpha1.EffectCodeExecute,
		v1alpha1.EffectSourceRead,
	}
	decision := Admit(request, contract, nowValue())
	denyReason(t, ReasonMissingPrepareVerify, decision)
}

// TestSymlinkEscapeProviderScopeDeny asserts provider_scope_mismatch when a
// provider path traverses outside the working scope derived from the subject
// refs, including a symlink-style escape.
func TestSymlinkEscapeProviderScopeDeny(t *testing.T) {
	request := baseRequest()
	for _, escaped := range []string{
		"work/other/link",
		"../escape",
		"work/ws/../../escape",
	} {
		contract := faultContract()
		contract.Spec.ProviderRef.ID = escaped
		decision := Admit(request, contract, nowValue())
		denyReason(t, ReasonProviderScopeMismatch, decision)
	}
}

// TestRedirectToPrivateResourceDeny asserts missing_network_authority when
// the contract declares network.read but the authority lacks the grant.
func TestRedirectToPrivateResourceDeny(t *testing.T) {
	contract := faultContract()
	contract.Spec.Network = v1alpha1.NetworkRead
	contract.Spec.Effects = []v1alpha1.Effect{v1alpha1.EffectNetworkRead}
	request := baseRequest()
	request.Authority.Effects = []v1alpha1.Effect{
		v1alpha1.EffectSourceRead,
		v1alpha1.EffectCodeExecute,
	}
	decision := Admit(request, contract, nowValue())
	denyReason(t, ReasonMissingNetworkAuthority, decision)
}

// TestRedirectToPrivateResourceExplicitURLDeny asserts the same gate for an
// absolute provider URL outside the working scope.
func TestRedirectToPrivateResourceExplicitURLDeny(t *testing.T) {
	contract := faultContract()
	contract.Spec.ProviderRef.ID = "https://example.com/private-resource"
	request := baseRequest()
	decision := Admit(request, contract, nowValue())
	if decision.ReasonCode != ReasonProviderScopeMismatch &&
		decision.ReasonCode != ReasonMissingNetworkAuthority {
		t.Fatalf(
			"expected scope or network denial, got %q with detail %v",
			decision.ReasonCode,
			decision.ReasonDetail,
		)
	}
}

// TestInfiniteRuntimeCodeBudgetExceeded asserts budget_exceeded when the
// contract cost exceeds the authority budget on any axis.
func TestInfiniteRuntimeCodeBudgetExceeded(t *testing.T) {
	contract := faultContract()
	contract.Spec.Cost.DurationMS = 1 << 30
	decision := Admit(baseRequest(), contract, nowValue())
	denyReason(t, ReasonBudgetExceeded, decision)
}

// TestCancellationMissingDeny asserts missing_cancellation_budget when the
// contract declares no cancellation plan.
func TestCancellationMissingDeny(t *testing.T) {
	contract := faultContract()
	contract.Spec.Cancellation = ""
	decision := Admit(baseRequest(), contract, nowValue())
	denyReason(t, ReasonMissingCancellationBudget, decision)
}

// TestConcurrentWritersStalePreStateDeny asserts that a second admission on
// the same request with a mutated pre-state digest is denied stale_pre_state.
func TestConcurrentWritersStalePreStateDeny(t *testing.T) {
	contract := faultContract()
	contract.Spec.Effects = []v1alpha1.Effect{v1alpha1.EffectRemoteWrite}
	contract.Spec.Cost.Calls = 100
	first := baseRequest()
	first.RemoteWrite = true
	first.PrepareVerified = true
	first.ExpectedPreState = "digest-expected"
	first.PreStateDigest = "digest-expected"
	if decision := Admit(first, contract, nowValue()); !decision.Admitted {
		t.Fatalf("first writer must admit, got %q", decision.ReasonCode)
	}
	second := first
	second.PreStateDigest = "digest-other"
	decision := Admit(second, contract, nowValue())
	denyReason(t, ReasonStalePreState, decision)
}

// TestBudgetExceededCallsDeny asserts budget_exceeded on the calls axis.
func TestBudgetExceededCallsDeny(t *testing.T) {
	contract := faultContract()
	contract.Spec.Cost.Calls = 1000
	decision := Admit(baseRequest(), contract, nowValue())
	denyReason(t, ReasonBudgetExceeded, decision)
}

// TestBudgetExceededBytesDeny asserts budget_exceeded on the bytes axis.
func TestBudgetExceededBytesDeny(t *testing.T) {
	contract := faultContract()
	contract.Spec.Cost.Bytes = 1 << 40
	decision := Admit(baseRequest(), contract, nowValue())
	denyReason(t, ReasonBudgetExceeded, decision)
}

// TestReadCheapPathAdmitted asserts the fast path reason code and effect set
// for a read-only request within budget.
func TestReadCheapPathAdmitted(t *testing.T) {
	decision := Admit(baseRequest(), faultContract(), nowValue())
	if !decision.Admitted {
		t.Fatalf("expected cheap-path admission, got %q", decision.ReasonCode)
	}
	if decision.ReasonCode != ReasonCodeReadCheapPath {
		t.Fatalf("expected %q, got %q", ReasonCodeReadCheapPath, decision.ReasonCode)
	}
	want := []v1alpha1.Effect{
		v1alpha1.EffectCodeExecute,
		v1alpha1.EffectSourceRead,
	}
	if !reflect.DeepEqual(decision.EffectSet, want) {
		t.Fatalf("effect set = %v, want %v", decision.EffectSet, want)
	}
	if got := AdmittedEffectSet(decision); !reflect.DeepEqual(got, want) {
		t.Fatalf("AdmittedEffectSet() = %v, want %v", got, want)
	}
}

// TestRemoteWritePrepareVerifyAdmitted asserts an effectful write requires
// prepare verification, cancellation budget, and a matching grant and admits
// only then.
func TestRemoteWritePrepareVerifyAdmitted(t *testing.T) {
	contract := faultContract()
	contract.Spec.Effects = []v1alpha1.Effect{v1alpha1.EffectRemoteWrite}
	contract.Spec.Cost.Calls = 100
	request := baseRequest()
	request.RemoteWrite = true
	request.PrepareVerified = true
	request.ExpectedPreState = "digest-expected"
	request.PreStateDigest = "digest-expected"
	request.Authority.Effects = append(
		request.Authority.Effects,
		v1alpha1.EffectRemoteWrite,
	)
	request.Budget.Calls = 100
	request.Budget.DurationMS = 120000
	decision := Admit(request, contract, nowValue())
	if !decision.Admitted {
		t.Fatalf("expected write admission, got %q", decision.ReasonCode)
	}
	if decision.ReasonCode != ReasonCodeAdmitted {
		t.Fatalf("expected %q, got %q", ReasonCodeAdmitted, decision.ReasonCode)
	}
	want := []v1alpha1.Effect{v1alpha1.EffectRemoteWrite}
	if !reflect.DeepEqual(decision.EffectSet, want) {
		t.Fatalf("effect set = %v, want %v", decision.EffectSet, want)
	}
}

// TestReasonDetailListsEveryViolation asserts the denied decision carries
// every violated rule, not only the first.
func TestReasonDetailListsEveryViolation(t *testing.T) {
	contract := faultContract()
	contract.Spec.EnvironmentSelector = "sandbox"
	contract.Spec.Effects = []v1alpha1.Effect{v1alpha1.EffectSourceWrite}
	request := baseRequest()
	request.Authority.ExpiresAt = "2026-09-01T11:59:59Z"
	request.Environment = []string{"prod"}
	request.Budget = v1alpha1.Budget{}
	decision := Admit(request, contract, nowValue())
	denyReason(t, ReasonUnknownAction, decision)
	joined := strings.Join(decision.ReasonDetail, "|")
	for _, code := range []string{
		ReasonUnknownAction,
		ReasonAuthorityExpired,
		ReasonWrongEnvironment,
		ReasonMissingCancellationBudget,
	} {
		if !strings.Contains(joined, code) {
			t.Fatalf("reason detail %v misses %q", decision.ReasonDetail, code)
		}
	}
}

// TestStaticContractProviderLookup asserts the provider interface resolves
// declared operations and fails closed otherwise.
func TestStaticContractProviderLookup(t *testing.T) {
	operation := validRef(v1alpha1.ReferenceOperation, "op/test")
	provider := StaticContractProvider{
		Operations: map[v1alpha1.Reference]v1alpha1.OperationContract{
			operation: faultContract(),
		},
	}
	contract, err := provider.LookupOperation(context.Background(), operation)
	if err != nil {
		t.Fatalf("LookupOperation() error = %v", err)
	}
	if contract.Metadata.Name != "op/test" {
		t.Fatalf("contract = %+v", contract)
	}
	if _, err := provider.LookupOperation(
		context.Background(),
		validRef(v1alpha1.ReferenceOperation, "op/missing"),
	); err == nil {
		t.Fatal("LookupOperation() must fail closed for undeclared operations")
	}
}

// TestAdmittedEffectSetDeniedNil asserts a denied decision exposes no effect
// set through the accessor.
func TestAdmittedEffectSetDeniedNil(t *testing.T) {
	contract := faultContract()
	contract.Spec.Cancellation = ""
	decision := Admit(baseRequest(), contract, nowValue())
	if got := AdmittedEffectSet(decision); got != nil {
		t.Fatalf("denied AdmittedEffectSet() = %v, want nil", got)
	}
}

// TestProviderInterfaceStaticContractProvider asserts StaticContractProvider
// satisfies ContractProvider at compile time.
func TestProviderInterfaceStaticContractProvider(t *testing.T) {
	var _ ContractProvider = StaticContractProvider{
		Operations: map[v1alpha1.Reference]v1alpha1.OperationContract{},
	}
}

// TestEnvironmentSelectorAuthorityLists asserts the environment selector is
// accepted only when both request and authority carry it.
func TestEnvironmentSelectorAuthorityLists(t *testing.T) {
	contract := faultContract()
	contract.Spec.EnvironmentSelector = "sandbox"
	request := baseRequest()
	request.Environment = []string{"sandbox"}
	request.Authority.EnvironmentSelectors = []string{"prod"}
	decision := Admit(request, contract, nowValue())
	denyReason(t, ReasonWrongEnvironment, decision)

	request.Authority.EnvironmentSelectors = []string{"sandbox"}
	if decision := Admit(request, contract, nowValue()); !decision.Admitted {
		t.Fatalf(
			"expected admission for matching environment, got %q",
			decision.ReasonCode,
		)
	}
}

// TestCheapPathNeedsNoPrepareVerification asserts the read/hermetic fast
// path admits without prepare verification, preflight, or pre-state binding.
func TestCheapPathNeedsNoPrepareVerification(t *testing.T) {
	contract := faultContract()
	contract.Spec.Preflight = ""
	contract.Spec.Verification = ""
	request := baseRequest()
	request.PrepareVerified = false
	decision := Admit(request, contract, nowValue())
	if !decision.Admitted {
		t.Fatalf(
			"expected cheap-path admission, got %q with detail %v",
			decision.ReasonCode,
			decision.ReasonDetail,
		)
	}
	if decision.ReasonCode != ReasonCodeReadCheapPath {
		t.Fatalf(
			"expected %q, got %q",
			ReasonCodeReadCheapPath,
			decision.ReasonCode,
		)
	}
}

// TestDestroyRequiresPrepareVerify asserts destroy demands the same
// prepare-verify gate as remote write.
func TestDestroyRequiresPrepareVerify(t *testing.T) {
	contract := faultContract()
	contract.Spec.Effects = []v1alpha1.Effect{v1alpha1.EffectRemoteDestroy}
	contract.Spec.Cost.Calls = 100
	request := baseRequest()
	request.Destroy = true
	request.PrepareVerified = false
	decision := Admit(request, contract, nowValue())
	denyReason(t, ReasonMissingPrepareVerify, decision)
}

// TestCredentialConsumeRequiresGrant asserts missing_credential_authority
// when a credential-consuming contract has no credential grant.
func TestCredentialConsumeRequiresGrant(t *testing.T) {
	contract := faultContract()
	contract.Spec.UsesCredentials = true
	contract.Spec.Effects = append(
		contract.Spec.Effects,
		v1alpha1.EffectCredentialConsume,
	)
	request := baseRequest()
	request.Authority.Effects = []v1alpha1.Effect{
		v1alpha1.EffectSourceRead,
		v1alpha1.EffectCodeExecute,
	}
	decision := Admit(request, contract, nowValue())
	denyReason(t, ReasonMissingCredentialAuthority, decision)
}

// TestNetworkWriteRequiresRemoteWriteGrant asserts missing_network_authority
// when network write has no remote.write grant.
func TestNetworkWriteRequiresRemoteWriteGrant(t *testing.T) {
	contract := faultContract()
	contract.Spec.Network = v1alpha1.NetworkWrite
	contract.Spec.Effects = []v1alpha1.Effect{v1alpha1.EffectRemoteWrite}
	contract.Spec.Cost.Calls = 100
	contract.Spec.Cancellation = "cancel"
	request := baseRequest()
	request.Authority.Effects = []v1alpha1.Effect{v1alpha1.EffectSourceRead}
	request.Budget.Calls = 100
	request.Budget.DurationMS = 60000
	request.RemoteWrite = true
	request.PrepareVerified = true
	decision := Admit(request, contract, nowValue())
	denyReason(t, ReasonMissingNetworkAuthority, decision)
}

// TestBudgetConcurrencyExceededDeny asserts budget_exceeded on the
// concurrency axis.
func TestBudgetConcurrencyExceededDeny(t *testing.T) {
	contract := faultContract()
	contract.Spec.Cost.Concurrency = 4
	decision := Admit(baseRequest(), contract, nowValue())
	denyReason(t, ReasonBudgetExceeded, decision)
}

// TestDeniedDecisionCarriesNoEffectSet asserts denyReason decisions expose
// an empty effect set.
func TestDeniedDecisionCarriesNoEffectSet(t *testing.T) {
	contract := faultContract()
	contract.Spec.Effects = []v1alpha1.Effect{v1alpha1.EffectSourceWrite}
	request := baseRequest()
	request.Authority.Effects = []v1alpha1.Effect{
		v1alpha1.EffectSourceRead,
		v1alpha1.EffectCodeExecute,
	}
	decision := Admit(request, contract, nowValue())
	denyReason(t, ReasonUnknownAction, decision)
	if decision.EffectSet != nil {
		t.Fatalf("denied EffectSet = %v, want nil", decision.EffectSet)
	}
}

// TestAdmittedEffectSetAccessor asserts the admitted-only effect set
// accessor behavior on the API decision type without redefining it.
func TestAdmittedEffectSetAccessor(t *testing.T) {
	admitted := Admit(baseRequest(), faultContract(), nowValue())
	if got := AdmittedEffectSet(admitted); got == nil {
		t.Fatal("admitted AdmittedEffectSet() = nil, want effect set")
	}
	contract := faultContract()
	contract.Spec.Cancellation = ""
	denied := Admit(baseRequest(), contract, nowValue())
	if got := AdmittedEffectSet(denied); got != nil {
		t.Fatalf("denied AdmittedEffectSet() = %v, want nil", got)
	}
}
