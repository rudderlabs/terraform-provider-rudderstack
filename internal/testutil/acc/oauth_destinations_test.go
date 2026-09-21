package acc

import (
	"reflect"
	"testing"
	"time"

	"github.com/rudderlabs/rudder-iac/api/client"
	"github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func TestSelectOAuthAccount(t *testing.T) {
	oldest := time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)
	newer := oldest.Add(time.Hour)
	accounts := []client.Account{
		newAccount("nil-created-at", "expected", "destination", nil),
		newAccount("wrong-category", "expected", "source", &oldest),
		newAccount("wrong-type", "other", "destination", &oldest),
		newAccount("newer", "expected", "destination", &newer),
		newAccount("oldest-b", "expected", "destination", &oldest),
		newAccount("oldest-a", "expected", "destination", &oldest),
	}

	got, ok := selectOAuthAccount(accounts, "expected")
	if !ok {
		t.Fatal("expected a matching OAuth account")
	}
	if got.ID != "oldest-a" {
		t.Fatalf("selected account ID = %q, want %q", got.ID, "oldest-a")
	}
}

func TestSelectOAuthAccountOrdersNilCreatedAtByID(t *testing.T) {
	accounts := []client.Account{
		newAccount("account-b", "expected", "destination", nil),
		newAccount("account-a", "expected", "destination", nil),
	}

	got, ok := selectOAuthAccount(accounts, "expected")
	if !ok {
		t.Fatal("expected a matching OAuth account")
	}
	if got.ID != "account-a" {
		t.Fatalf("selected account ID = %q, want %q", got.ID, "account-a")
	}
}

func TestSelectOAuthAccountNoMatch(t *testing.T) {
	accounts := []client.Account{
		newAccount("wrong-category", "expected", "source", nil),
		newAccount("wrong-type", "other", "destination", nil),
	}

	if _, ok := selectOAuthAccount(accounts, "expected"); ok {
		t.Fatal("expected no matching OAuth account")
	}
}

func TestSubstituteAccountIDCopiesAllConfigFields(t *testing.T) {
	original := []configs.TestConfig{{
		TerraformCreate: oauthAccountIDPlaceholder,
		TerraformUpdate: "before-" + oauthAccountIDPlaceholder,
		APICreate:       oauthAccountIDPlaceholder + "-after",
		APIUpdate:       oauthAccountIDPlaceholder + oauthAccountIDPlaceholder,
	}}

	got := substituteAccountID(original, "real-account-id")

	if original[0].TerraformCreate != oauthAccountIDPlaceholder || original[0].APIUpdate != oauthAccountIDPlaceholder+oauthAccountIDPlaceholder {
		t.Fatal("substituteAccountID mutated the caller's fixtures")
	}
	want := configs.TestConfig{
		TerraformCreate: "real-account-id",
		TerraformUpdate: "before-real-account-id",
		APICreate:       "real-account-id-after",
		APIUpdate:       "real-account-idreal-account-id",
	}
	if !reflect.DeepEqual(got[0], want) {
		t.Fatalf("substituted config = %#v, want %#v", got[0], want)
	}
}

func newAccount(id, accountType, category string, createdAt *time.Time) client.Account {
	account := client.Account{ID: id, CreatedAt: createdAt}
	account.Definition.Type = accountType
	account.Definition.Category = category
	return account
}
