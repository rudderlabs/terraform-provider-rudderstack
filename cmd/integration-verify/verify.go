package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/go-cmp/cmp"
	"github.com/rudderlabs/rudder-api-go/client"
	"github.com/tidwall/gjson"

	"github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

// IntegrationResource holds parsed information about an onboarded integration
// resource extracted from the terraform state.
type IntegrationResource struct {
	Kind            string // "destination" or "source"
	IntegrationType string // e.g., "webhook", "slack"
	ResourceType    string // full terraform type e.g., "rudderstack_destination_webhook"
	Name            string // resource label
	ConfigState     string // JSON string from terraform state's config.0
	ResourceID      string // resource ID from terraform state
}

// VerifyResult holds the outcome of comparing the onboarded integration's
// Terraform config against the live API response.
type VerifyResult struct {
	Match bool
	Diff  string // go-cmp diff output, empty when Match is true
}

// ParseTerraformState parses `terraform show -json` output and extracts
// IntegrationResource entries for rudderstack resources.
// If targetResource is empty, all rudderstack resources are returned.
func ParseTerraformState(stateJSON []byte, targetResource string) ([]*IntegrationResource, error) {
	resources := gjson.GetBytes(stateJSON, "values.root_module.resources")
	if !resources.Exists() || !resources.IsArray() {
		return nil, fmt.Errorf("no resources found in terraform state")
	}

	var result []*IntegrationResource
	for _, r := range resources.Array() {
		rType := r.Get("type").String()
		if !strings.HasPrefix(rType, "rudderstack_") {
			continue
		}

		rName := r.Get("name").String()
		if targetResource != "" && rName != targetResource {
			continue
		}

		kind, integrationType, err := ExtractResourceType(rType)
		if err != nil {
			continue
		}

		configState := r.Get("values.config.0").Raw
		if configState == "" {
			configState = "{}"
		}

		resourceID := r.Get("values.id").String()

		result = append(result, &IntegrationResource{
			Kind:            kind,
			IntegrationType: integrationType,
			ResourceType:    rType,
			Name:            rName,
			ConfigState:     configState,
			ResourceID:      resourceID,
		})
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("no matching rudderstack resource found in terraform state")
	}

	return result, nil
}

// ExtractResourceType parses "rudderstack_destination_webhook" into ("destination", "webhook").
func ExtractResourceType(resourceType string) (kind, integrationType string, err error) {
	trimmed := strings.TrimPrefix(resourceType, "rudderstack_")
	if strings.HasPrefix(trimmed, "destination_") {
		return "destination", strings.TrimPrefix(trimmed, "destination_"), nil
	}
	if strings.HasPrefix(trimmed, "source_") {
		return "source", strings.TrimPrefix(trimmed, "source_"), nil
	}
	return "", "", fmt.Errorf("unrecognized resource type: %s", resourceType)
}

// FetchResourceConfig fetches the live config for an onboarded integration from
// the RudderStack API.
func FetchResourceConfig(ctx context.Context, cl *client.Client, kind, id string) (json.RawMessage, error) {
	switch kind {
	case "destination":
		dest, err := cl.Destinations.Get(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("fetching destination %s: %w", id, err)
		}
		return dest.Config, nil
	case "source":
		src, err := cl.Sources.Get(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("fetching source %s: %w", id, err)
		}
		return src.Config, nil
	default:
		return nil, fmt.Errorf("unknown resource kind: %s", kind)
	}
}

// Verify performs the full verification of an onboarded integration:
//  1. Converts the Terraform state config into expected API JSON via StateToAPI.
//  2. Fetches the actual config from the live RudderStack API.
//  3. Runs a subset comparison — every key from the state config must exist and
//     match in the API response.
func Verify(ctx context.Context, cl *client.Client, info *IntegrationResource) (*VerifyResult, error) {
	// Look up ConfigMeta from the provider's integration registry.
	var cm configs.ConfigMeta
	var found bool
	switch info.Kind {
	case "destination":
		cm, found = configs.Destinations.Entries()[info.IntegrationType]
	case "source":
		cm, found = configs.Sources.Entries()[info.IntegrationType]
	}
	if !found {
		return nil, fmt.Errorf("integration type %q not found in %s registry — was it onboarded correctly?", info.IntegrationType, info.Kind)
	}

	// Convert state to expected API format using the integration's property mappings.
	expectedAPI, err := cm.StateToAPI(info.ConfigState)
	if err != nil {
		return nil, fmt.Errorf("StateToAPI conversion: %w", err)
	}

	// Fetch actual config from the live API.
	actualConfig, err := FetchResourceConfig(ctx, cl, info.Kind, info.ResourceID)
	if err != nil {
		return nil, err
	}

	// Unmarshal both for comparison.
	var expectedMap, actualMap map[string]interface{}
	if err := json.Unmarshal([]byte(expectedAPI), &expectedMap); err != nil {
		return nil, fmt.Errorf("unmarshaling expected API config: %w", err)
	}
	if err := json.Unmarshal(actualConfig, &actualMap); err != nil {
		return nil, fmt.Errorf("unmarshaling actual API config: %w", err)
	}

	diff := cmp.Diff(expectedMap, actualMap)

	return &VerifyResult{
		Match: diff == "",
		Diff:  diff,
	}, nil
}
