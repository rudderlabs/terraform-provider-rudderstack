package rudderstack

import (
	"testing"

	"github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
	_ "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/integrations"
)

// TestProvider validates the provider schema and that key resources/data sources are registered.
func TestProvider(t *testing.T) {
	p := New()
	if err := p.InternalValidate(); err != nil {
		t.Fatalf("provider InternalValidate: %v", err)
	}
}

// TestProviderAccountRegistration confirms that the account resource and data source
// are registered in the provider after wiring (DEX-380).
func TestProviderAccountRegistration(t *testing.T) {
	p := New()

	if _, ok := p.ResourcesMap["rudderstack_account_source_bigquery"]; !ok {
		t.Error("expected rudderstack_account_source_bigquery in ResourcesMap, but it is missing")
	}
	if _, ok := p.DataSourcesMap["rudderstack_account"]; !ok {
		t.Error("expected rudderstack_account in DataSourcesMap, but it is missing")
	}

	// Spot-check that existing registrations are still present.
	if _, ok := p.ResourcesMap["rudderstack_connection"]; !ok {
		t.Error("expected rudderstack_connection in ResourcesMap, but it is missing")
	}
}

// TestProviderDoesNotRegisterCustomActivation prevents Terraform from creating
// or updating Custom Activation destinations. That integration's authentication
// contract is account-only and is not supported by this provider.
func TestProviderDoesNotRegisterCustomActivation(t *testing.T) {
	if _, ok := configs.Destinations.Entries()["custom_audience"]; ok {
		t.Fatal("custom_audience must not be registered as a Terraform destination")
	}
	for name, config := range configs.Destinations.Entries() {
		if config.APIType == "CUSTOM_AUDIENCE" {
			t.Fatalf("CUSTOM_AUDIENCE must not be registered as Terraform destination %q", name)
		}
	}

	if _, ok := New().ResourcesMap["rudderstack_destination_custom_audience"]; ok {
		t.Fatal("rudderstack_destination_custom_audience must not be exposed by the provider")
	}
}
