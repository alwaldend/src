package v1alpha1

import (
	"bytes"
	"strings"
	"testing"
)

func TestOperationContractRoundTripsDeterministically(t *testing.T) {
	contract := validOperationContract()
	first, err := CanonicalOperationJSON(contract)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := DecodeOperationContract(first)
	if err != nil {
		t.Fatal(err)
	}
	second, err := CanonicalOperationJSON(decoded)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(first, second) {
		t.Fatalf("canonical round trip changed bytes:\nfirst: %s\nsecond: %s", first, second)
	}
	withUnknown := bytes.Replace(first, []byte(`"kind":"AgentOperation"`),
		[]byte(`"kind":"AgentOperation","unknown":true`), 1)
	if _, err := DecodeOperationContract(withUnknown); err == nil ||
		!strings.Contains(err.Error(), "unknown field") {
		t.Fatalf("DecodeOperationContract() error = %v, want unknown-field rejection", err)
	}
}

func TestArtifactEnvelopeRoundTripsDeterministically(t *testing.T) {
	envelope := ArtifactEnvelope{
		APIVersion:    APIVersion,
		Kind:          "validation.set",
		ID:            "validation.phase1",
		ProducerRef:   Reference{Kind: ReferenceProvider, ID: "bazel.agent"},
		AuthorityRefs: []Reference{{Kind: ReferenceArtifact, ID: "task.authority"}},
		InputRefs: []InputReference{{
			Reference: Reference{Kind: ReferenceSubject, ID: "git.tree"},
			Role:      "candidate",
		}},
		SubjectRef:   &Reference{Kind: ReferenceSubject, ID: "git.tree"},
		Completeness: CompletenessComplete,
		Information:  InformationSet{Public: true},
		Retention:    RetentionTask,
	}
	first, err := CanonicalArtifactJSON(envelope)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := DecodeArtifactEnvelope(first)
	if err != nil {
		t.Fatal(err)
	}
	second, err := CanonicalArtifactJSON(decoded)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(first, second) {
		t.Fatalf("canonical round trip changed bytes:\nfirst: %s\nsecond: %s", first, second)
	}
}

func TestArtifactEnvelopeRequiresLimitationsWhenIncomplete(t *testing.T) {
	for _, completeness := range []Completeness{
		CompletenessPartial,
		CompletenessTruncated,
		CompletenessUnknown,
	} {
		t.Run(string(completeness), func(t *testing.T) {
			envelope := ArtifactEnvelope{
				APIVersion:   APIVersion,
				Kind:         "validation.set",
				ID:           "validation.phase1",
				ProducerRef:  Reference{Kind: ReferenceProvider, ID: "bazel.agent"},
				Completeness: completeness,
				Information:  InformationSet{Public: true},
				Retention:    RetentionTask,
			}
			if _, err := CanonicalArtifactJSON(envelope); err == nil ||
				!strings.Contains(err.Error(), "must declare at least one limitation") {
				t.Fatalf("CanonicalArtifactJSON() error = %v, want limitation rejection", err)
			}
			envelope.Limitations = []string{"field unavailable"}
			if _, err := CanonicalArtifactJSON(envelope); err != nil {
				t.Fatalf("incomplete envelope with limitation rejected: %v", err)
			}
		})
	}
}

func TestOperationContractRejectsMalformedIdentityAndUnknownEffect(t *testing.T) {
	for _, test := range []struct {
		name   string
		mutate func(*OperationContract)
		want   string
	}{
		{
			name: "malformed identity",
			mutate: func(contract *OperationContract) {
				contract.Metadata.Name = "../apply"
			},
			want: "malformed operation identity",
		},
		{
			name: "unknown effect",
			mutate: func(contract *OperationContract) {
				contract.Spec.Effects = []Effect{"remote.superuser"}
			},
			want: "unknown effectful operation",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			contract := validOperationContract()
			test.mutate(&contract)
			if _, err := CanonicalOperationJSON(contract); err == nil ||
				!strings.Contains(err.Error(), test.want) {
				t.Fatalf("CanonicalOperationJSON() error = %v, want %q", err, test.want)
			}
		})
	}
}

func TestOperationContractRejectsUnderdeclaredEffects(t *testing.T) {
	for _, test := range []struct {
		name   string
		mutate func(*OperationContract)
		want   string
	}{
		{
			name: "credentials",
			mutate: func(contract *OperationContract) {
				contract.Spec.UsesCredentials = true
			},
			want: string(EffectCredentialConsume),
		},
		{
			name: "network read",
			mutate: func(contract *OperationContract) {
				contract.Spec.Network = NetworkRead
			},
			want: string(EffectNetworkRead),
		},
		{
			name: "network write",
			mutate: func(contract *OperationContract) {
				contract.Spec.Network = NetworkWrite
			},
			want: string(EffectRemoteWrite),
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			contract := validOperationContract()
			test.mutate(&contract)
			if _, err := CanonicalOperationJSON(contract); err == nil ||
				!strings.Contains(err.Error(), test.want) {
				t.Fatalf("CanonicalOperationJSON() error = %v, want %q", err, test.want)
			}
		})
	}
}

func TestDelegationRejectsEveryAuthorityWideningAxis(t *testing.T) {
	parent := AuthorityEnvelope{
		ActorRef:             Reference{Kind: ReferenceActor, ID: "agent.primary"},
		SubjectRefs:          []Reference{{Kind: ReferenceSubject, ID: "branch.feature"}},
		Effects:              []Effect{EffectSourceRead, EffectSourceWrite},
		EnvironmentSelectors: []string{"workspace.local"},
		Budget:               Budget{Calls: 10, Bytes: 1024, DurationMS: 5000, Concurrency: 2},
		ExpiresAt:            "2026-08-31T22:00:00Z",
	}
	child := parent
	child.ActorRef = Reference{Kind: ReferenceActor, ID: "agent.worker"}
	child.Effects = []Effect{EffectSourceRead}
	child.Budget = Budget{Calls: 5, Bytes: 512, DurationMS: 1000, Concurrency: 1}
	child.ExpiresAt = "2026-08-31T21:00:00Z"
	if err := ValidateDelegation(parent, child); err != nil {
		t.Fatalf("narrow delegation rejected: %v", err)
	}

	for _, test := range []struct {
		name   string
		mutate func(*AuthorityEnvelope)
		want   string
	}{
		{
			name: "effect",
			mutate: func(value *AuthorityEnvelope) {
				value.Effects = append(value.Effects, EffectRemoteWrite)
			},
			want: "widens effects",
		},
		{
			name: "subject",
			mutate: func(value *AuthorityEnvelope) {
				value.SubjectRefs = append(value.SubjectRefs,
					Reference{Kind: ReferenceSubject, ID: "branch.other"})
			},
			want: "widens subjects",
		},
		{
			name: "environment",
			mutate: func(value *AuthorityEnvelope) {
				value.EnvironmentSelectors = append(value.EnvironmentSelectors, "production")
			},
			want: "widens environments",
		},
		{
			name: "budget",
			mutate: func(value *AuthorityEnvelope) {
				value.Budget.Calls = 11
			},
			want: "widens budget",
		},
		{
			name: "expiry",
			mutate: func(value *AuthorityEnvelope) {
				value.ExpiresAt = "2026-08-31T23:00:00Z"
			},
			want: "widens expiry",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			widened := child
			test.mutate(&widened)
			if err := ValidateDelegation(parent, widened); err == nil ||
				!strings.Contains(err.Error(), test.want) {
				t.Fatalf("ValidateDelegation() error = %v, want %q", err, test.want)
			}
		})
	}
}

func TestInformationClassesRemainIndependentFromPathPolicyAxes(t *testing.T) {
	policy := PathPolicy{
		Disclosure:    "public_source",
		BuildConsumer: "repository_internal",
		Publication:   "not_publishable",
		Information:   InformationSet{Public: true, Personal: true},
	}
	if err := policy.Validate(); err != nil {
		t.Fatalf("orthogonal policy rejected: %v", err)
	}
	if err := ValidateInformationFlow(
		InformationSet{Public: true, Secret: true},
		InformationSet{Public: true},
	); err == nil || !strings.Contains(err.Error(), "does not retain") {
		t.Fatalf("ValidateInformationFlow() error = %v, want incompatible-flow rejection", err)
	}
}

func TestTaskRunManifestRoundTripsAndRejectsSharedOwnership(t *testing.T) {
	manifest := TaskRunManifest{
		APIVersion:  APIVersion,
		Kind:        KindTaskRunManifest,
		TaskID:      "task-one",
		RunID:       "run-one",
		WorkerID:    "worker-one",
		Information: InformationSet{Public: true, Secret: true, Personal: true},
		Budget: Budget{
			Calls: 1, Bytes: 1_048_576, DurationMS: 30_000, Concurrency: 1,
		},
		Retention:    RetentionTask,
		LockScope:    "task-one/run-one",
		CleanupOwner: "task-owner",
	}
	first, err := CanonicalTaskRunJSON(manifest)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := DecodeTaskRunManifest(first)
	if err != nil {
		t.Fatal(err)
	}
	second, err := CanonicalTaskRunJSON(decoded)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(first, second) {
		t.Fatalf("task-run round trip changed bytes")
	}
	manifest.LockScope = "shared"
	if _, err := CanonicalTaskRunJSON(manifest); err == nil ||
		!strings.Contains(err.Error(), "lock scope") {
		t.Fatalf("shared lock scope error = %v", err)
	}
}

func validOperationContract() OperationContract {
	return OperationContract{
		APIVersion: APIVersion,
		Kind:       KindOperationContract,
		Metadata: OperationMetadata{
			Name:     "goal.validate",
			OwnerRef: Reference{Kind: ReferenceRepository, ID: "alwaldend.src"},
		},
		Spec: OperationSpec{
			ProviderRef:       Reference{Kind: ReferenceProvider, ID: "goal.cli"},
			Effects:           []Effect{EffectSourceRead},
			InputInformation:  InformationSet{Public: true},
			OutputInformation: InformationSet{Public: true},
			Network:           NetworkNone,
			AuthorityGate:     "source-read",
			Preflight:         "validate-record-path",
			Verification:      "validated-record",
			Cost:              Budget{Calls: 1, Bytes: 65536, DurationMS: 5000, Concurrency: 1},
			Cacheability:      NotCacheable,
			Cancellation:      "process-context",
		},
	}
}
