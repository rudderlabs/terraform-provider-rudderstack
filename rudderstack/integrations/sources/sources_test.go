package sources_test

import (
	"testing"

	cmt "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/cm"
	"github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func TestSourceResourceHTTP(t *testing.T) {
	cmt.AssertSource(t, "http", []configs.TestConfig{configs.EmptyTestConfig})
}

func TestSourceResourceJavascript(t *testing.T) {
	cmt.AssertSource(t, "javascript", []configs.TestConfig{configs.EmptyTestConfig})
}
