// Package v1alpha1 defines shared repository-agent control contracts.
//
// Payload schemas and persistence policy remain with their owning components.
package v1alpha1

const (
	APIVersion = "agents.alwaldend.com/v1alpha1"

	KindOperationContract = "AgentOperation"
	KindTaskRunManifest   = "TaskRunManifest"
)

type ReferenceKind string

const (
	ReferenceRepository ReferenceKind = "repository"
	ReferenceWorkspace  ReferenceKind = "workspace"
	ReferenceTask       ReferenceKind = "task"
	ReferenceSession    ReferenceKind = "session"
	ReferenceActor      ReferenceKind = "actor"
	ReferenceProvider   ReferenceKind = "provider"
	ReferenceSubject    ReferenceKind = "subject"
	ReferenceArtifact   ReferenceKind = "artifact"
	ReferenceOperation  ReferenceKind = "operation"
)

type Reference struct {
	Kind    ReferenceKind `json:"kind"`
	ID      string        `json:"id"`
	Version string        `json:"version,omitempty"`
	Digest  string        `json:"digest,omitempty"`
}

type InputReference struct {
	Reference Reference `json:"reference"`
	Role      string    `json:"role"`
}

type Effect string

const (
	EffectSourceRead        Effect = "source.read"
	EffectSourceWrite       Effect = "source.write"
	EffectTaskStateWrite    Effect = "task_state.write"
	EffectHostWrite         Effect = "host.write"
	EffectHistoryWrite      Effect = "history.write"
	EffectCodeExecute       Effect = "code.execute"
	EffectCredentialConsume Effect = "credential.consume"
	EffectNetworkRead       Effect = "network.read"
	EffectRemoteWrite       Effect = "remote.write"
	EffectRemoteDestroy     Effect = "remote.destroy"
)

type Budget struct {
	Calls       int64 `json:"calls,omitempty"`
	Bytes       int64 `json:"bytes,omitempty"`
	DurationMS  int64 `json:"durationMs,omitempty"`
	Concurrency int64 `json:"concurrency,omitempty"`
}

type AuthorityEnvelope struct {
	ActorRef             Reference   `json:"actorRef"`
	SubjectRefs          []Reference `json:"subjectRefs"`
	Effects              []Effect    `json:"effects"`
	EnvironmentSelectors []string    `json:"environmentSelectors"`
	Budget               Budget      `json:"budget"`
	ExpiresAt            string      `json:"expiresAt,omitempty"`
}

type InformationSet struct {
	Public   bool `json:"public,omitempty"`
	Secret   bool `json:"secret,omitempty"`
	Personal bool `json:"personal,omitempty"`
}

type PathPolicy struct {
	Disclosure    string         `json:"disclosure"`
	BuildConsumer string         `json:"buildConsumer"`
	Publication   string         `json:"publication"`
	Information   InformationSet `json:"information"`
}

type Completeness string

const (
	CompletenessComplete  Completeness = "complete"
	CompletenessPartial   Completeness = "partial"
	CompletenessTruncated Completeness = "truncated"
	CompletenessUnknown   Completeness = "unknown"
)

type RetentionClass string

const (
	RetentionEphemeral RetentionClass = "ephemeral"
	RetentionTask      RetentionClass = "task"
	RetentionDurable   RetentionClass = "durable"
)

type AvailabilityReason string

const (
	ReasonUnavailable       AvailabilityReason = "unavailable"
	ReasonStale             AvailabilityReason = "stale"
	ReasonPartial           AvailabilityReason = "partial"
	ReasonTruncated         AvailabilityReason = "truncated"
	ReasonDenied            AvailabilityReason = "denied"
	ReasonConflict          AvailabilityReason = "conflict"
	ReasonRequiresMigration AvailabilityReason = "requires_migration"
)

type EvidenceApplicability string

const (
	EvidenceApplicable   EvidenceApplicability = "applicable"
	EvidenceInapplicable EvidenceApplicability = "inapplicable"
	EvidenceStale        EvidenceApplicability = "stale"
	EvidenceUnknown      EvidenceApplicability = "unknown"
)

type ArtifactEnvelope struct {
	APIVersion    string           `json:"apiVersion"`
	Kind          string           `json:"kind"`
	ID            string           `json:"id"`
	ProducerRef   Reference        `json:"producerRef"`
	AuthorityRefs []Reference      `json:"authorityRefs,omitempty"`
	InputRefs     []InputReference `json:"inputRefs,omitempty"`
	SubjectRef    *Reference       `json:"subjectRef,omitempty"`
	Completeness  Completeness     `json:"completeness"`
	Limitations   []string         `json:"limitations,omitempty"`
	Information   InformationSet   `json:"information"`
	Retention     RetentionClass   `json:"retention"`
	Digest        string           `json:"digest,omitempty"`
}

type NetworkAccess string

const (
	NetworkNone  NetworkAccess = "none"
	NetworkRead  NetworkAccess = "read"
	NetworkWrite NetworkAccess = "write"
)

type Cacheability string

const (
	Cacheable           Cacheability = "cacheable"
	NotCacheable        Cacheability = "not_cacheable"
	CacheabilityUnknown Cacheability = "unknown"
)

type OperationMetadata struct {
	Name     string    `json:"name"`
	OwnerRef Reference `json:"ownerRef"`
}

type OperationSpec struct {
	ProviderRef         Reference      `json:"providerRef"`
	Effects             []Effect       `json:"effects"`
	InputInformation    InformationSet `json:"inputInformation"`
	OutputInformation   InformationSet `json:"outputInformation"`
	UsesCredentials     bool           `json:"usesCredentials"`
	Network             NetworkAccess  `json:"network"`
	EnvironmentSelector string         `json:"environmentSelector,omitempty"`
	AuthorityGate       string         `json:"authorityGate"`
	Preflight           string         `json:"preflight"`
	Verification        string         `json:"verification"`
	Cost                Budget         `json:"cost"`
	Cacheability        Cacheability   `json:"cacheability"`
	Cancellation        string         `json:"cancellation"`
}

type OperationContract struct {
	APIVersion string            `json:"apiVersion"`
	Kind       string            `json:"kind"`
	Metadata   OperationMetadata `json:"metadata"`
	Spec       OperationSpec     `json:"spec"`
}

type TaskRunManifest struct {
	APIVersion   string         `json:"apiVersion"`
	Kind         string         `json:"kind"`
	TaskID       string         `json:"taskId"`
	RunID        string         `json:"runId"`
	WorkerID     string         `json:"workerId"`
	Information  InformationSet `json:"information"`
	Budget       Budget         `json:"budget"`
	Retention    RetentionClass `json:"retention"`
	LockScope    string         `json:"lockScope"`
	CleanupOwner string         `json:"cleanupOwner"`
}
