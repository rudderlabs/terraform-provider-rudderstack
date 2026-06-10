package rudderstack

import (
	"testing"

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
