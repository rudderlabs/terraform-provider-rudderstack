package integrations_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
	_ "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/integrations"
)

func TestIntegrations(t *testing.T) {
	// importing integrations package should add entries to Sources/Destinations registries
	assert.Greater(t, len(configs.Sources.Entries()), 0)
	assert.Greater(t, len(configs.Destinations.Entries()), 0)
}
