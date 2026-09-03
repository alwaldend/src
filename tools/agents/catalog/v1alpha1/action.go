package catalogv1alpha1

import (
	"encoding/json"
	"fmt"
)

const (
	// KindActionCatalog identifies the action catalog kind.
	KindActionCatalog = "action-catalog"
)

// ActionProvider is one declared operation provider.
type ActionProvider struct {
	ID             string `json:"id"`
	Owner          string `json:"owner"`
	DefinitionPath string `json:"definitionPath"`
}

// ActionRecord is one classified provider action.
type ActionRecord struct {
	ID                  string   `json:"id"`
	ProviderRef         string   `json:"providerRef"`
	Owner               string   `json:"owner"`
	SourcePath          string   `json:"sourcePath"`
	Selector            string   `json:"selector"`
	Classification      string   `json:"classification"`
	Effects             []string `json:"effects"`
	Inputs              []string `json:"inputs"`
	Outputs             []string `json:"outputs"`
	Information         []string `json:"information"`
	CredentialUse       string   `json:"credentialUse"`
	NetworkUse          string   `json:"networkUse"`
	EnvironmentSelector string   `json:"environmentSelector"`
	AuthorityGate       string   `json:"authorityGate"`
	Preflight           string   `json:"preflight"`
	Verification        string   `json:"verification"`
	Cost                string   `json:"cost"`
	Cacheability        string   `json:"cacheability"`
	Cancellation        string   `json:"cancellation"`
}

// ActionAlias is one removed or replaced selector.
type ActionAlias struct {
	ProviderRef    string `json:"providerRef"`
	Selector       string `json:"selector"`
	State          string `json:"state"`
	ReplacementRef string `json:"replacementRef,omitempty"`
	Reason         string `json:"reason,omitempty"`
}

// ActionCatalog is the bounded action projection over the registered
// operation files and their owner-local definitions.
type ActionCatalog struct {
	CatalogEnvelope
	Providers []ActionProvider `json:"providers"`
	Actions   []ActionRecord   `json:"actions"`
	Aliases   []ActionAlias    `json:"aliases"`
}

// Validate checks action-specific invariants and the shared envelope.
func (catalog ActionCatalog) Validate() error {
	if catalog.Providers == nil || catalog.Actions == nil || catalog.Aliases == nil {
		return fmt.Errorf("action catalog arrays must be non-null")
	}
	itemIDs := make([]string, 0,
		len(catalog.Providers)+len(catalog.Actions)+len(catalog.Aliases))
	providers := map[string]bool{}
	for _, provider := range catalog.Providers {
		if provider.ID == "" || provider.Owner == "" || provider.DefinitionPath == "" {
			return fmt.Errorf("action provider %q is incomplete", provider.ID)
		}
		if providers[provider.ID] {
			return fmt.Errorf("duplicate action provider %q", provider.ID)
		}
		providers[provider.ID] = true
		itemIDs = append(itemIDs, "provider."+provider.ID)
	}
	seenActions := map[string]bool{}
	seenSelector := map[string]bool{}
	for _, action := range catalog.Actions {
		if action.ID == "" || action.ProviderRef == "" ||
			action.Selector == "" || action.Classification == "" ||
			action.AuthorityGate == "" || action.Preflight == "" ||
			action.Verification == "" || action.Cost == "" ||
			action.Cancellation == "" {
			return fmt.Errorf("action %q is incomplete", action.ID)
		}
		if !providers[action.ProviderRef] {
			return fmt.Errorf("action %s refers to unknown provider %q",
				action.ID, action.ProviderRef)
		}
		if seenActions[action.ID] {
			return fmt.Errorf("duplicate action identity %q", action.ID)
		}
		seenActions[action.ID] = true
		selectorKey := action.ProviderRef + "." + action.Selector
		if seenSelector[selectorKey] {
			return fmt.Errorf("duplicate provider-local selector %q", selectorKey)
		}
		seenSelector[selectorKey] = true
		if len(action.Effects) == 0 {
			return fmt.Errorf("action %s has no effects", action.ID)
		}
		for _, effect := range action.Effects {
			if !oneOf(effect,
				"source.read", "source.write", "task_state.write", "host.write",
				"history.write", "code.execute", "credential.consume",
				"network.read", "remote.write", "remote.destroy") {
				return fmt.Errorf("action %s has unknown effect %q", action.ID, effect)
			}
		}
		if !oneOf(action.Cacheability, "cacheable", "not_cacheable", "unknown") {
			return fmt.Errorf("action %s has unknown cacheability %q",
				action.ID, action.Cacheability)
		}
		if action.SourcePath == "" || absPath(action.SourcePath) ||
			escapes(action.SourcePath) {
			return fmt.Errorf("action %s has malformed source path %q",
				action.ID, action.SourcePath)
		}
		itemIDs = append(itemIDs, "action."+action.ID)
	}
	for _, alias := range catalog.Aliases {
		if alias.ProviderRef == "" || alias.Selector == "" || alias.State == "" {
			return fmt.Errorf("action alias is incomplete")
		}
		if !oneOf(alias.State, "removed", "replaced") {
			return fmt.Errorf("action alias has unknown state %q", alias.State)
		}
		if alias.State == "replaced" && alias.ReplacementRef == "" {
			return fmt.Errorf("replaced alias %s.%s lacks a replacement",
				alias.ProviderRef, alias.Selector)
		}
		itemIDs = append(itemIDs, "alias."+alias.ProviderRef+"."+alias.Selector)
	}
	return catalog.CatalogEnvelope.Validate(itemIDs)
}

// ItemIDs returns the sorted set-like identity list for an action catalog.
func (catalog ActionCatalog) ItemIDs() []string {
	ids := make([]string, 0,
		len(catalog.Providers)+len(catalog.Actions)+len(catalog.Aliases))
	for _, provider := range catalog.Providers {
		ids = append(ids, "provider."+provider.ID)
	}
	for _, action := range catalog.Actions {
		ids = append(ids, "action."+action.ID)
	}
	for _, alias := range catalog.Aliases {
		ids = append(ids, "alias."+alias.ProviderRef+"."+alias.Selector)
	}
	return ids
}

// VerifySelfDigest recomputes the content-addressed digest over the canonical
// bytes of the complete action document and compares it to the stored digest,
// rejecting an empty or mismatched digest.
func (catalog ActionCatalog) VerifySelfDigest() error {
	withoutDigest := catalog
	withoutDigest.Digest = ""
	return verifySelfDigest(catalog.Digest, withoutDigest)
}

// CanonicalJSONAction encodes a complete action catalog with a digest over
// the canonical bytes with the digest field omitted.
func CanonicalJSONAction(catalog ActionCatalog) ([]byte, error) {
	if err := catalog.Validate(); err != nil {
		return nil, err
	}
	withoutDigest := catalog
	withoutDigest.Digest = ""
	content, err := json.Marshal(withoutDigest)
	if err != nil {
		return nil, fmt.Errorf("encode action catalog: %w", err)
	}
	withoutDigest.Digest = digest(content)
	contentWithDigest, err := json.MarshalIndent(withoutDigest, "", "    ")
	if err != nil {
		return nil, fmt.Errorf("encode action catalog with digest: %w", err)
	}
	return append(contentWithDigest, '\n'), nil
}

// DecodeActionStrict decodes and validates a complete action catalog
// document with unknown-field and trailing-JSON rejection.
func DecodeActionStrict(content []byte) (ActionCatalog, error) {
	var catalog ActionCatalog
	if err := DecodeStrict(content, &catalog); err != nil {
		return ActionCatalog{}, err
	}
	if err := catalog.Validate(); err != nil {
		return ActionCatalog{}, err
	}
	if err := catalog.VerifySelfDigest(); err != nil {
		return ActionCatalog{}, err
	}
	return catalog, nil
}
