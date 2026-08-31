package v1alpha1

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"
)

var identifierPattern = regexp.MustCompile(`^[a-z][a-z0-9]*(?:[._/-][a-z0-9]+)*$`)

func (reference Reference) Validate() error {
	if !oneOfReferenceKind(reference.Kind) {
		return fmt.Errorf("unknown reference kind %q", reference.Kind)
	}
	if len(reference.ID) > 253 || !identifierPattern.MatchString(reference.ID) ||
		strings.Contains(reference.ID, "..") {
		return fmt.Errorf("malformed %s identity %q", reference.Kind, reference.ID)
	}
	if reference.Digest != "" && !validDigest(reference.Digest) {
		return fmt.Errorf("malformed digest for %s %q", reference.Kind, reference.ID)
	}
	return nil
}

func (budget Budget) Validate() error {
	if budget.Calls < 0 || budget.Bytes < 0 || budget.DurationMS < 0 ||
		budget.Concurrency < 0 {
		return fmt.Errorf("budget values cannot be negative")
	}
	return nil
}

func (information InformationSet) Validate() error {
	if !information.Public && !information.Secret && !information.Personal {
		return fmt.Errorf("information classification is empty")
	}
	return nil
}

func ValidateInformationFlow(input InformationSet, output InformationSet) error {
	if err := input.Validate(); err != nil {
		return fmt.Errorf("input information: %w", err)
	}
	if err := output.Validate(); err != nil {
		return fmt.Errorf("output information: %w", err)
	}
	if input.Public && !output.Public || input.Secret && !output.Secret ||
		input.Personal && !output.Personal {
		return fmt.Errorf("output classification does not retain every input class")
	}
	return nil
}

func (policy PathPolicy) Validate() error {
	if !oneOf(policy.Disclosure, "public_source", "restricted_source") {
		return fmt.Errorf("unknown disclosure policy %q", policy.Disclosure)
	}
	if !oneOf(policy.BuildConsumer, "production", "repository_internal", "test_only", "forbidden") {
		return fmt.Errorf("unknown build-consumer policy %q", policy.BuildConsumer)
	}
	if !oneOf(policy.Publication, "publishable", "not_publishable") {
		return fmt.Errorf("unknown publication policy %q", policy.Publication)
	}
	return policy.Information.Validate()
}

func (envelope ArtifactEnvelope) Validate() error {
	if envelope.APIVersion != APIVersion {
		return fmt.Errorf("unsupported artifact API version %q", envelope.APIVersion)
	}
	if !identifierPattern.MatchString(envelope.Kind) {
		return fmt.Errorf("malformed artifact kind %q", envelope.Kind)
	}
	if !identifierPattern.MatchString(envelope.ID) {
		return fmt.Errorf("malformed artifact identity %q", envelope.ID)
	}
	if err := envelope.ProducerRef.Validate(); err != nil {
		return fmt.Errorf("producer reference: %w", err)
	}
	if envelope.ProducerRef.Kind != ReferenceProvider {
		return fmt.Errorf("artifact producer must be a provider reference")
	}
	for _, reference := range envelope.AuthorityRefs {
		if err := reference.Validate(); err != nil {
			return fmt.Errorf("authority reference: %w", err)
		}
	}
	for _, input := range envelope.InputRefs {
		if err := input.Reference.Validate(); err != nil {
			return fmt.Errorf("input reference: %w", err)
		}
		if !identifierPattern.MatchString(input.Role) {
			return fmt.Errorf("malformed input role %q", input.Role)
		}
	}
	if envelope.SubjectRef != nil {
		if err := envelope.SubjectRef.Validate(); err != nil {
			return fmt.Errorf("subject reference: %w", err)
		}
		if envelope.SubjectRef.Kind != ReferenceSubject {
			return fmt.Errorf("artifact subject must be a subject reference")
		}
	}
	if !oneOf(
		string(envelope.Completeness),
		string(CompletenessComplete),
		string(CompletenessPartial),
		string(CompletenessTruncated),
		string(CompletenessUnknown),
	) {
		return fmt.Errorf("unknown completeness %q", envelope.Completeness)
	}
	if envelope.Completeness != CompletenessComplete && len(envelope.Limitations) == 0 {
		return fmt.Errorf("incomplete artifact must declare at least one limitation")
	}
	for _, limitation := range envelope.Limitations {
		if strings.TrimSpace(limitation) == "" {
			return fmt.Errorf("artifact limitation cannot be empty")
		}
	}
	if err := envelope.Information.Validate(); err != nil {
		return err
	}
	if !oneOf(
		string(envelope.Retention),
		string(RetentionEphemeral),
		string(RetentionTask),
		string(RetentionDurable),
	) {
		return fmt.Errorf("unknown retention class %q", envelope.Retention)
	}
	if envelope.Digest != "" && !validDigest(envelope.Digest) {
		return fmt.Errorf("malformed artifact digest")
	}
	return nil
}

func (contract OperationContract) Validate() error {
	if contract.APIVersion != APIVersion || contract.Kind != KindOperationContract {
		return fmt.Errorf("unsupported operation type %q %q", contract.APIVersion, contract.Kind)
	}
	if !identifierPattern.MatchString(contract.Metadata.Name) {
		return fmt.Errorf("malformed operation identity %q", contract.Metadata.Name)
	}
	if err := contract.Metadata.OwnerRef.Validate(); err != nil {
		return fmt.Errorf("owner reference: %w", err)
	}
	if contract.Metadata.OwnerRef.Kind != ReferenceRepository &&
		contract.Metadata.OwnerRef.Kind != ReferenceWorkspace {
		return fmt.Errorf("operation owner must be a repository or workspace reference")
	}
	if err := contract.Spec.ProviderRef.Validate(); err != nil {
		return fmt.Errorf("provider reference: %w", err)
	}
	if contract.Spec.ProviderRef.Kind != ReferenceProvider {
		return fmt.Errorf("operation provider must be a provider reference")
	}
	if len(contract.Spec.Effects) == 0 {
		return fmt.Errorf("operation must declare at least one effect")
	}
	seen := map[Effect]bool{}
	for _, effect := range contract.Spec.Effects {
		if !knownEffect(effect) {
			return fmt.Errorf("unknown effectful operation %q", effect)
		}
		if seen[effect] {
			return fmt.Errorf("duplicate effect %q", effect)
		}
		seen[effect] = true
	}
	if !sort.SliceIsSorted(contract.Spec.Effects, func(i, j int) bool {
		return contract.Spec.Effects[i] < contract.Spec.Effects[j]
	}) {
		return fmt.Errorf("effects must use stable lexical order")
	}
	if err := ValidateInformationFlow(
		contract.Spec.InputInformation,
		contract.Spec.OutputInformation,
	); err != nil {
		return err
	}
	if !oneOf(string(contract.Spec.Network), string(NetworkNone), string(NetworkRead), string(NetworkWrite)) {
		return fmt.Errorf("unknown network access %q", contract.Spec.Network)
	}
	if contract.Spec.UsesCredentials && !seen[EffectCredentialConsume] {
		return fmt.Errorf("credential use requires effect %q", EffectCredentialConsume)
	}
	if contract.Spec.Network == NetworkRead && !seen[EffectNetworkRead] {
		return fmt.Errorf("network read requires effect %q", EffectNetworkRead)
	}
	if contract.Spec.Network == NetworkWrite && !seen[EffectRemoteWrite] {
		return fmt.Errorf("network write requires effect %q", EffectRemoteWrite)
	}
	if !oneOf(
		string(contract.Spec.Cacheability),
		string(Cacheable),
		string(NotCacheable),
		string(CacheabilityUnknown),
	) {
		return fmt.Errorf("unknown cacheability %q", contract.Spec.Cacheability)
	}
	required := []struct {
		name  string
		value string
	}{
		{name: "authority gate", value: contract.Spec.AuthorityGate},
		{name: "preflight", value: contract.Spec.Preflight},
		{name: "verification", value: contract.Spec.Verification},
		{name: "cancellation", value: contract.Spec.Cancellation},
	}
	for _, field := range required {
		if strings.TrimSpace(field.value) == "" {
			return fmt.Errorf("%s is required", field.name)
		}
	}
	return contract.Spec.Cost.Validate()
}

func (manifest TaskRunManifest) Validate() error {
	if manifest.APIVersion != APIVersion || manifest.Kind != KindTaskRunManifest {
		return fmt.Errorf(
			"unsupported task-run manifest type %q %q",
			manifest.APIVersion,
			manifest.Kind,
		)
	}
	for name, value := range map[string]string{
		"task":   manifest.TaskID,
		"run":    manifest.RunID,
		"worker": manifest.WorkerID,
	} {
		if !identifierPattern.MatchString(value) || strings.Contains(value, "..") {
			return fmt.Errorf("malformed %s identity %q", name, value)
		}
	}
	if err := manifest.Information.Validate(); err != nil {
		return err
	}
	if err := manifest.Budget.Validate(); err != nil {
		return err
	}
	if manifest.Budget.Calls == 0 || manifest.Budget.Bytes == 0 ||
		manifest.Budget.DurationMS == 0 || manifest.Budget.Concurrency == 0 {
		return fmt.Errorf("task-run budget must bound every resource axis")
	}
	if manifest.Retention != RetentionTask &&
		manifest.Retention != RetentionEphemeral {
		return fmt.Errorf("task-run retention must be task or ephemeral")
	}
	if manifest.LockScope != manifest.TaskID+"/"+manifest.RunID {
		return fmt.Errorf("task-run lock scope does not match its ownership")
	}
	if manifest.CleanupOwner != "task-owner" {
		return fmt.Errorf("task-run cleanup owner must be task-owner")
	}
	return nil
}

func ValidateDelegation(parent AuthorityEnvelope, child AuthorityEnvelope) error {
	if err := validateAuthority(parent); err != nil {
		return fmt.Errorf("parent authority: %w", err)
	}
	if err := validateAuthority(child); err != nil {
		return fmt.Errorf("child authority: %w", err)
	}
	if !referenceSubset(parent.SubjectRefs, child.SubjectRefs) {
		return fmt.Errorf("delegation widens subjects")
	}
	if !effectSubset(parent.Effects, child.Effects) {
		return fmt.Errorf("delegation widens effects")
	}
	if !stringSubset(parent.EnvironmentSelectors, child.EnvironmentSelectors) {
		return fmt.Errorf("delegation widens environments")
	}
	if budgetWider(parent.Budget, child.Budget) {
		return fmt.Errorf("delegation widens budget")
	}
	parentExpiry, parentSet, err := parseExpiry(parent.ExpiresAt)
	if err != nil {
		return fmt.Errorf("parent authority: %w", err)
	}
	childExpiry, childSet, err := parseExpiry(child.ExpiresAt)
	if err != nil {
		return fmt.Errorf("child authority: %w", err)
	}
	if parentSet && (!childSet || childExpiry.After(parentExpiry)) {
		return fmt.Errorf("delegation widens expiry")
	}
	return nil
}

func validateAuthority(authority AuthorityEnvelope) error {
	if err := authority.ActorRef.Validate(); err != nil {
		return err
	}
	if authority.ActorRef.Kind != ReferenceActor {
		return fmt.Errorf("authority actor must be an actor reference")
	}
	for _, subject := range authority.SubjectRefs {
		if err := subject.Validate(); err != nil {
			return err
		}
		if subject.Kind != ReferenceSubject {
			return fmt.Errorf("authority subject must be a subject reference")
		}
	}
	for _, effect := range authority.Effects {
		if !knownEffect(effect) {
			return fmt.Errorf("unknown effect %q", effect)
		}
	}
	if err := authority.Budget.Validate(); err != nil {
		return err
	}
	_, _, err := parseExpiry(authority.ExpiresAt)
	return err
}

func knownEffect(effect Effect) bool {
	return oneOf(
		string(effect),
		string(EffectSourceRead),
		string(EffectSourceWrite),
		string(EffectTaskStateWrite),
		string(EffectHostWrite),
		string(EffectHistoryWrite),
		string(EffectCodeExecute),
		string(EffectCredentialConsume),
		string(EffectNetworkRead),
		string(EffectRemoteWrite),
		string(EffectRemoteDestroy),
	)
}

func oneOfReferenceKind(kind ReferenceKind) bool {
	return oneOf(
		string(kind),
		string(ReferenceRepository),
		string(ReferenceWorkspace),
		string(ReferenceTask),
		string(ReferenceSession),
		string(ReferenceActor),
		string(ReferenceProvider),
		string(ReferenceSubject),
		string(ReferenceArtifact),
		string(ReferenceOperation),
	)
}

func oneOf(value string, allowed ...string) bool {
	for _, candidate := range allowed {
		if value == candidate {
			return true
		}
	}
	return false
}

func validDigest(value string) bool {
	if !strings.HasPrefix(value, "sha256:") || len(value) != len("sha256:")+64 {
		return false
	}
	for _, character := range value[len("sha256:"):] {
		if !strings.ContainsRune("0123456789abcdef", character) {
			return false
		}
	}
	return true
}

func referenceSubset(parent []Reference, child []Reference) bool {
	allowed := map[Reference]bool{}
	for _, value := range parent {
		allowed[value] = true
	}
	for _, value := range child {
		if !allowed[value] {
			return false
		}
	}
	return true
}

func effectSubset(parent []Effect, child []Effect) bool {
	allowed := map[Effect]bool{}
	for _, value := range parent {
		allowed[value] = true
	}
	for _, value := range child {
		if !allowed[value] {
			return false
		}
	}
	return true
}

func stringSubset(parent []string, child []string) bool {
	allowed := map[string]bool{}
	for _, value := range parent {
		allowed[value] = true
	}
	for _, value := range child {
		if !allowed[value] {
			return false
		}
	}
	return true
}

func budgetWider(parent Budget, child Budget) bool {
	return child.Calls > parent.Calls || child.Bytes > parent.Bytes ||
		child.DurationMS > parent.DurationMS || child.Concurrency > parent.Concurrency
}

func parseExpiry(value string) (time.Time, bool, error) {
	if value == "" {
		return time.Time{}, false, nil
	}
	parsed, err := time.Parse(time.RFC3339Nano, value)
	if err != nil {
		return time.Time{}, false, fmt.Errorf("expiry must be RFC3339: %w", err)
	}
	return parsed, true, nil
}
