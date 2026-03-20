package main

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rudderlabs/rudder-api-go/client"
	"github.com/zclconf/go-cty/cty"

	"github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

// IntegrationResource holds parsed information about an onboarded integration
// resource extracted from a .tf file.
type IntegrationResource struct {
	Kind            string // "destination" or "source"
	IntegrationType string // e.g., "webhook", "slack"
	ResourceType    string // full terraform type e.g., "rudderstack_destination_webhook"
	Name            string // resource label
	ConfigState     string // JSON string representing config state
}

// VerifyResult holds the outcome of comparing the onboarded integration's
// Terraform config against the live API response.
type VerifyResult struct {
	Match       bool
	Expected    string
	Actual      string
	Differences []string
}

// ParseTFFile reads a .tf file produced during the /onboard-integration E2E step
// and extracts the IntegrationResource for the target resource.
// If targetResource is empty, the first rudderstack resource block is used.
func ParseTFFile(filePath, targetResource string) (*IntegrationResource, error) {
	parser := hclparse.NewParser()
	file, diags := parser.ParseHCLFile(filePath)
	if diags.HasErrors() {
		return nil, fmt.Errorf("parsing HCL file: %s", diags.Error())
	}

	body, ok := file.Body.(*hclsyntax.Body)
	if !ok {
		return nil, fmt.Errorf("unexpected HCL body type")
	}

	for _, block := range body.Blocks {
		if block.Type != "resource" || len(block.Labels) < 2 {
			continue
		}

		resourceType := block.Labels[0]
		resourceName := block.Labels[1]

		if !strings.HasPrefix(resourceType, "rudderstack_") {
			continue
		}

		if targetResource != "" && resourceName != targetResource {
			continue
		}

		kind, integrationType, err := ExtractResourceType(resourceType)
		if err != nil {
			continue
		}

		// Look up ConfigMeta from the provider's integration registry.
		var cm configs.ConfigMeta
		var found bool
		switch kind {
		case "destination":
			cm, found = configs.Destinations.Entries()[integrationType]
		case "source":
			cm, found = configs.Sources.Entries()[integrationType]
		}
		if !found {
			return nil, fmt.Errorf("integration type %q not found in %s registry — was it onboarded correctly?", integrationType, kind)
		}

		// Find the config block inside the resource block.
		var configBody *hclsyntax.Body
		for _, innerBlock := range block.Body.Blocks {
			if innerBlock.Type == "config" {
				configBody = innerBlock.Body
				break
			}
		}

		var stateJSON string
		if configBody != nil && cm.ConfigSchema != nil {
			stateMap, err := HCLBodyToStateMap(configBody, cm.ConfigSchema)
			if err != nil {
				return nil, fmt.Errorf("converting HCL config to state map: %w", err)
			}
			stateBytes, err := json.Marshal(stateMap)
			if err != nil {
				return nil, fmt.Errorf("marshaling state map: %w", err)
			}
			stateJSON = string(stateBytes)
		} else {
			stateJSON = "{}"
		}

		return &IntegrationResource{
			Kind:            kind,
			IntegrationType: integrationType,
			ResourceType:    resourceType,
			Name:            resourceName,
			ConfigState:     stateJSON,
		}, nil
	}

	return nil, fmt.Errorf("no matching rudderstack resource block found in %s", filePath)
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

// HCLBodyToStateMap converts an HCL body (the config block) into a map[string]interface{}
// that matches what Terraform's d.Get("config.0") returns.
func HCLBodyToStateMap(body *hclsyntax.Body, configSchema map[string]*schema.Schema) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	// Process attributes.
	for name, attr := range body.Attributes {
		sch, ok := configSchema[name]
		if !ok {
			continue
		}

		val, diags := attr.Expr.Value(nil)
		if diags.HasErrors() {
			return nil, fmt.Errorf("evaluating attribute %q: %s", name, diags.Error())
		}

		goVal, err := ctyToGo(val, sch)
		if err != nil {
			return nil, fmt.Errorf("converting attribute %q: %w", name, err)
		}
		result[name] = goVal
	}

	// Process blocks (TypeList with Elem: *schema.Resource, no ConfigModeAttr).
	for _, block := range body.Blocks {
		sch, ok := configSchema[block.Type]
		if !ok {
			continue
		}

		if sch.Type != schema.TypeList {
			continue
		}
		elemResource, ok := sch.Elem.(*schema.Resource)
		if !ok {
			continue
		}

		nestedMap, err := HCLBodyToStateMap(block.Body, elemResource.Schema)
		if err != nil {
			return nil, fmt.Errorf("processing block %q: %w", block.Type, err)
		}

		// Blocks accumulate as arrays. If we already have entries, append.
		if existing, ok := result[block.Type]; ok {
			result[block.Type] = append(existing.([]interface{}), nestedMap)
		} else {
			result[block.Type] = []interface{}{nestedMap}
		}
	}

	return result, nil
}

// ctyToGo converts a cty.Value to a Go interface{} using the schema for guidance.
func ctyToGo(val cty.Value, sch *schema.Schema) (interface{}, error) {
	if val.IsNull() {
		return nil, nil
	}

	switch sch.Type {
	case schema.TypeString:
		return val.AsString(), nil
	case schema.TypeBool:
		return val.True(), nil
	case schema.TypeInt:
		bf := val.AsBigFloat()
		i, _ := bf.Int64()
		return i, nil
	case schema.TypeFloat:
		bf := val.AsBigFloat()
		f, _ := bf.Float64()
		return f, nil
	case schema.TypeList:
		return ctyListToGo(val, sch)
	default:
		return nil, fmt.Errorf("unsupported schema type: %v", sch.Type)
	}
}

// ctyListToGo converts a cty list/tuple value to a Go slice.
func ctyListToGo(val cty.Value, sch *schema.Schema) (interface{}, error) {
	if val.IsNull() || !val.IsKnown() || !val.CanIterateElements() {
		return []interface{}{}, nil
	}

	switch elem := sch.Elem.(type) {
	case *schema.Resource:
		// Array of objects (ConfigModeAttr style).
		var items []interface{}
		it := val.ElementIterator()
		for it.Next() {
			_, v := it.Element()
			if v.Type().IsObjectType() || v.Type().IsMapType() {
				objMap := make(map[string]interface{})
				for key, nestedSch := range elem.Schema {
					attrVal := v.GetAttr(key)
					goVal, err := ctyToGo(attrVal, nestedSch)
					if err != nil {
						return nil, fmt.Errorf("converting nested field %q: %w", key, err)
					}
					objMap[key] = goVal
				}
				items = append(items, objMap)
			}
		}
		return items, nil

	case *schema.Schema:
		// Array of primitives.
		var items []interface{}
		it := val.ElementIterator()
		for it.Next() {
			_, v := it.Element()
			goVal, err := ctyToGo(v, elem)
			if err != nil {
				return nil, err
			}
			items = append(items, goVal)
		}
		return items, nil

	default:
		return nil, fmt.Errorf("unsupported list element type: %T", sch.Elem)
	}
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
//  1. Converts the Terraform .tf state into expected API JSON via StateToAPI.
//  2. Fetches the actual config from the live RudderStack API.
//  3. Runs a subset comparison — every key from the .tf config must exist and
//     match in the API response.
//
// This is the core validation step in the /onboard-integration E2E workflow.
func Verify(ctx context.Context, cl *client.Client, info *IntegrationResource, resourceID string) (*VerifyResult, error) {
	// Look up ConfigMeta from the provider's integration registry.
	var cm configs.ConfigMeta
	switch info.Kind {
	case "destination":
		cm = configs.Destinations.Entries()[info.IntegrationType]
	case "source":
		cm = configs.Sources.Entries()[info.IntegrationType]
	}

	// Convert state to expected API format using the integration's property mappings.
	expectedAPI, err := cm.StateToAPI(info.ConfigState)
	if err != nil {
		return nil, fmt.Errorf("StateToAPI conversion: %w", err)
	}

	// Fetch actual config from the live API.
	actualConfig, err := FetchResourceConfig(ctx, cl, info.Kind, resourceID)
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

	// Subset comparison: every key in the .tf-derived config must exist and
	// match in the API response. Extra keys in the API response are ignored.
	diffs := SubsetDiff(expectedMap, actualMap, "")

	expectedPretty, _ := json.MarshalIndent(expectedMap, "", "  ")
	actualPretty, _ := json.MarshalIndent(actualMap, "", "  ")

	return &VerifyResult{
		Match:       len(diffs) == 0,
		Expected:    string(expectedPretty),
		Actual:      string(actualPretty),
		Differences: diffs,
	}, nil
}

// SubsetDiff checks that every key in expected exists in actual with the same value.
// Returns a list of human-readable differences.
func SubsetDiff(expected, actual map[string]interface{}, prefix string) []string {
	var diffs []string

	keys := make([]string, 0, len(expected))
	for k := range expected {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		expVal := expected[key]
		fullKey := key
		if prefix != "" {
			fullKey = prefix + "." + key
		}

		actVal, exists := actual[key]
		if !exists {
			diffs = append(diffs, fmt.Sprintf("  %s: expected %v, but key is missing in API response", fullKey, expVal))
			continue
		}

		// Recurse into nested maps.
		expMap, expIsMap := expVal.(map[string]interface{})
		actMap, actIsMap := actVal.(map[string]interface{})
		if expIsMap && actIsMap {
			diffs = append(diffs, SubsetDiff(expMap, actMap, fullKey)...)
			continue
		}

		// Compare arrays element by element.
		expArr, expIsArr := expVal.([]interface{})
		actArr, actIsArr := actVal.([]interface{})
		if expIsArr && actIsArr {
			diffs = append(diffs, arrayDiff(expArr, actArr, fullKey)...)
			continue
		}

		// Direct comparison.
		if !reflect.DeepEqual(expVal, actVal) {
			diffs = append(diffs, fmt.Sprintf("  %s: expected %v, got %v", fullKey, expVal, actVal))
		}
	}

	return diffs
}

// arrayDiff compares two arrays element by element.
func arrayDiff(expected, actual []interface{}, prefix string) []string {
	var diffs []string

	if len(expected) != len(actual) {
		diffs = append(diffs, fmt.Sprintf("  %s: expected array length %d, got %d", prefix, len(expected), len(actual)))
		return diffs
	}

	for i := range expected {
		elemKey := fmt.Sprintf("%s[%d]", prefix, i)
		expMap, expIsMap := expected[i].(map[string]interface{})
		actMap, actIsMap := actual[i].(map[string]interface{})
		if expIsMap && actIsMap {
			diffs = append(diffs, SubsetDiff(expMap, actMap, elemKey)...)
		} else if !reflect.DeepEqual(expected[i], actual[i]) {
			diffs = append(diffs, fmt.Sprintf("  %s: expected %v, got %v", elemKey, expected[i], actual[i]))
		}
	}

	return diffs
}

// FormatResult formats the verification result for display.
func FormatResult(info *IntegrationResource, resourceID string, result *VerifyResult) string {
	var sb strings.Builder

	shortID := resourceID
	if len(shortID) > 8 {
		shortID = shortID[:8] + "..."
	}

	if result.Match {
		sb.WriteString(fmt.Sprintf("PASS: %s (ID: %s) — onboarded integration config matches .tf file\n", info.ResourceType, shortID))
	} else {
		sb.WriteString(fmt.Sprintf("FAIL: %s (ID: %s) — onboarded integration config mismatch\n", info.ResourceType, shortID))
	}

	sb.WriteString(fmt.Sprintf("\nExpected (from .tf → StateToAPI):\n%s\n", result.Expected))
	sb.WriteString(fmt.Sprintf("\nActual (from API):\n%s\n", result.Actual))

	if !result.Match {
		sb.WriteString("\nDifferences (fix these in the onboarded integration's .go file):\n")
		for _, d := range result.Differences {
			sb.WriteString(d + "\n")
		}
	}

	return sb.String()
}
