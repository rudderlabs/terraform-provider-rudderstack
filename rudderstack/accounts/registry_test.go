package accounts_test

import (
	"testing"

	"github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/accounts"
	"github.com/stretchr/testify/assert"
)

func TestAccountsRegistry(t *testing.T) {
	err := accounts.Accounts.Register("test", accounts.AccountConfigMeta{
		Category: accounts.CategorySource,
	})
	assert.NoError(t, err, "first registration should succeed")

	err = accounts.Accounts.Register("test", accounts.AccountConfigMeta{
		Category: accounts.CategorySource,
	})
	assert.Error(t, err, "duplicate registration should fail for the same name and category")

	err = accounts.Accounts.Register("test", accounts.AccountConfigMeta{
		Category: accounts.CategoryDestination,
	})
	assert.NoError(t, err, "duplicate registration should not fail for the same name and different category")

	categories := accounts.Accounts.Entries()
	assert.Len(t, categories, 2, "should have one entry for each category")
	_, ok := categories[accounts.CategorySource]["test"]
	assert.True(t, ok, "should have an entry for the source")
	_, ok = categories[accounts.CategoryDestination]["test"]
	assert.True(t, ok, "should have an entry for the destination")
}
