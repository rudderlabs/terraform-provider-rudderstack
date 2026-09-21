package acc

import (
	"context"
	"os"
	"sort"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/rudderlabs/rudder-iac/api/client"
	"github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

// oauthAccountIDPlaceholder is the literal used in shared destination test configs.
// Unit tests assert those fixtures exactly; only the acceptance path substitutes it.
const oauthAccountIDPlaceholder = "account-id-1"

// AccAssertOAuthDestination resolves a real workspace account for full CRUD runs.
// Plan-only runs retain the placeholder and make no account API calls.
func AccAssertOAuthDestination(t *testing.T, destination, accountType string, testConfigs []configs.TestConfig) {
	t.Helper()
	if os.Getenv(resource.EnvTfAcc) == "" || PlanOnly() {
		// Delegate so the SDK retains its normal TF_ACC-unset skip and plan-only behavior.
		AccAssertDestination(t, destination, testConfigs)
		return
	}

	TestAccPreCheck(t)
	accountID := resolveOAuthAccountID(t, accountType)
	AccAssertDestination(t, destination, substituteAccountID(testConfigs, accountID))
}

// resolveOAuthAccountID returns the oldest destination-category account whose
// definition type matches accountType. Lookup is the only source of the ID.
func resolveOAuthAccountID(t *testing.T, accountType string) string {
	t.Helper()

	cl, err := newTestAPIClient()
	if err != nil {
		t.Fatalf("create API client for OAuth account lookup: %v", err)
	}

	accounts, err := cl.Accounts.ListAll(context.Background())
	if err != nil {
		t.Fatalf("list OAuth accounts for type %q: %v", accountType, err)
	}

	account, ok := selectOAuthAccount(accounts, accountType)
	if !ok {
		t.Fatalf("no OAuth account found for type %q (category \"destination\") in the test workspace:\nconnect a %s account in the workspace that RUDDERSTACK_ACCESS_TOKEN belongs to,\nor run with TF_ACC_PLAN_ONLY=1 to validate the plan only", accountType, accountType)
	}

	t.Logf("using OAuth account id %q, name %q for type %q", account.ID, account.Name, accountType)
	return account.ID
}

func selectOAuthAccount(accounts []client.Account, accountType string) (client.Account, bool) {
	matches := make([]client.Account, 0)
	for _, account := range accounts {
		if account.Definition.Type == accountType && account.Definition.Category == "destination" {
			matches = append(matches, account)
		}
	}
	if len(matches) == 0 {
		return client.Account{}, false
	}

	sort.Slice(matches, func(i, j int) bool {
		left, right := matches[i], matches[j]
		switch {
		case left.CreatedAt == nil && right.CreatedAt != nil:
			return false
		case left.CreatedAt != nil && right.CreatedAt == nil:
			return true
		case left.CreatedAt != nil && right.CreatedAt != nil && !left.CreatedAt.Equal(*right.CreatedAt):
			return left.CreatedAt.Before(*right.CreatedAt)
		default:
			return left.ID < right.ID
		}
	})
	return matches[0], true
}

func substituteAccountID(testConfigs []configs.TestConfig, accountID string) []configs.TestConfig {
	substituted := make([]configs.TestConfig, len(testConfigs))
	copy(substituted, testConfigs)

	for i := range substituted {
		substituted[i].TerraformCreate = strings.ReplaceAll(substituted[i].TerraformCreate, oauthAccountIDPlaceholder, accountID)
		substituted[i].TerraformUpdate = strings.ReplaceAll(substituted[i].TerraformUpdate, oauthAccountIDPlaceholder, accountID)
		substituted[i].APICreate = strings.ReplaceAll(substituted[i].APICreate, oauthAccountIDPlaceholder, accountID)
		substituted[i].APIUpdate = strings.ReplaceAll(substituted[i].APIUpdate, oauthAccountIDPlaceholder, accountID)
	}
	return substituted
}
