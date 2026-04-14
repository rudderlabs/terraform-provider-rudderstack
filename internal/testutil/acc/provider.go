package acc

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/rudderlabs/terraform-provider-rudderstack/rudderstack"

	// Blank-import integrations so every registered source/destination is available.
	_ "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/integrations"
)

// TestAccProviderFactories returns provider factories that use a real API client.
// The provider reads RUDDERSTACK_ACCESS_TOKEN and RUDDERSTACK_API_URL from env.
var TestAccProviderFactories = map[string]func() (*schema.Provider, error){
	"rudderstack": func() (*schema.Provider, error) {
		return rudderstack.New(), nil
	},
}

// TestAccPreCheck validates required environment variables are set.
func TestAccPreCheck(t *testing.T) {
	t.Helper()
	if os.Getenv("RUDDERSTACK_ACCESS_TOKEN") == "" {
		t.Fatal("RUDDERSTACK_ACCESS_TOKEN must be set for acceptance tests")
	}
}

// PlanOnly returns true when TF_ACC_PLAN_ONLY=1 is set, indicating tests should
// only validate the Terraform plan (zero API calls).
func PlanOnly() bool {
	return os.Getenv("TF_ACC_PLAN_ONLY") == "1"
}

// ensureDummyToken sets a placeholder access token for plan-only tests so the
// provider can configure without error. No real API calls are made in plan-only mode.
// Uses os.Setenv (not t.Setenv) to stay compatible with t.Parallel().
func ensureDummyToken(t *testing.T) {
	t.Helper()
	if os.Getenv("RUDDERSTACK_ACCESS_TOKEN") == "" {
		os.Setenv("RUDDERSTACK_ACCESS_TOKEN", "plan-only-dummy-token")
	}
}
