package sources_test

import (
	"testing"

	cmt "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/cm"
	"github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func TestSourceResourceBraze(t *testing.T) {
	cmt.AssertSource(t, "braze", []configs.TestConfig{configs.EmptyTestConfig})
}

func TestSourceResourceCordova(t *testing.T) {
	cmt.AssertSource(t, "cordova", []configs.TestConfig{configs.EmptyTestConfig})
}

func TestSourceResourceGo(t *testing.T) {
	cmt.AssertSource(t, "go", []configs.TestConfig{configs.EmptyTestConfig})
}

func TestSourceResourceHTTP(t *testing.T) {
	cmt.AssertSource(t, "http", []configs.TestConfig{configs.EmptyTestConfig})
}

func TestSourceResourceAndroid(t *testing.T) {
	cmt.AssertSource(t, "android", []configs.TestConfig{configs.EmptyTestConfig})
}

func TestSourceResourceIOS(t *testing.T) {
	cmt.AssertSource(t, "ios", []configs.TestConfig{configs.EmptyTestConfig})
}

func TestSourceResourceJava(t *testing.T) {
	cmt.AssertSource(t, "java", []configs.TestConfig{configs.EmptyTestConfig})
}

func TestSourceResourceJavascript(t *testing.T) {
	cmt.AssertSource(t, "javascript", []configs.TestConfig{configs.EmptyTestConfig})
}

func TestSourceResourceNode(t *testing.T) {
	cmt.AssertSource(t, "node", []configs.TestConfig{configs.EmptyTestConfig})
}

func TestSourceResourceReactNative(t *testing.T) {
	cmt.AssertSource(t, "reactnative", []configs.TestConfig{configs.EmptyTestConfig})
}

func TestSourceResourceRuby(t *testing.T) {
	cmt.AssertSource(t, "ruby", []configs.TestConfig{configs.EmptyTestConfig})
}

func TestSourceResourceWebhook(t *testing.T) {
	cmt.AssertSource(t, "webhook", []configs.TestConfig{configs.EmptyTestConfig})
}

func TestSourceResourceWebhookShopify(t *testing.T) {
	cmt.AssertSource(t, "webhook_shopify", []configs.TestConfig{configs.EmptyTestConfig})
}

func TestSourceResourcePython(t *testing.T) {
	cmt.AssertSource(t, "python", []configs.TestConfig{configs.EmptyTestConfig})
}
