package sources_test

import (
	"testing"

	"github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil"
	"github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func TestSourceResourceHTTP(t *testing.T) {
	testutil.AssertSource(t, "http", []configs.TestConfig{configs.EmptyTestConfig})
}

func TestSourceResourceJavascript(t *testing.T) {
	testutil.AssertSource(t, "javascript", []configs.TestConfig{configs.EmptyTestConfig})
}
