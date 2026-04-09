package sources_test

import (
	"testing"

	acc "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/acc"
	cmt "github.com/rudderlabs/terraform-provider-rudderstack/internal/testutil/cm"
	"github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

var emptyTestConfigs = []configs.TestConfig{configs.EmptyTestConfig}

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

func TestSourceResourcePHP(t *testing.T) {
	cmt.AssertSource(t, "php", []configs.TestConfig{configs.EmptyTestConfig})
}

func TestSourceResourceDotNet(t *testing.T) {
	cmt.AssertSource(t, "dotnet", []configs.TestConfig{configs.EmptyTestConfig})
}

func TestSourceResourceFlutter(t *testing.T) {
	cmt.AssertSource(t, "flutter", []configs.TestConfig{configs.EmptyTestConfig})
}

func TestSourceResourceCustomerIO(t *testing.T) {
	cmt.AssertSource(t, "customerio", []configs.TestConfig{configs.EmptyTestConfig})
}

func TestSourceResourceFacebookLeadAds(t *testing.T) {
	cmt.AssertSource(t, "facebook_lead_ads", []configs.TestConfig{configs.EmptyTestConfig})
}

func TestSourceResourceAdjust(t *testing.T) {
	cmt.AssertSource(t, "adjust", []configs.TestConfig{configs.EmptyTestConfig})
}

func TestSourceResourceRust(t *testing.T) {
	cmt.AssertSource(t, "rust", []configs.TestConfig{configs.EmptyTestConfig})
}

// E2E acceptance tests — reuse the same empty test configs from unit tests.

func TestAccSourceBraze(t *testing.T) {
	acc.AccAssertSource(t, "braze", emptyTestConfigs)
}

func TestAccSourceCordova(t *testing.T) {
	acc.AccAssertSource(t, "cordova", emptyTestConfigs)
}

func TestAccSourceGo(t *testing.T) {
	acc.AccAssertSource(t, "go", emptyTestConfigs)
}

func TestAccSourceHTTP(t *testing.T) {
	acc.AccAssertSource(t, "http", emptyTestConfigs)
}

func TestAccSourceAndroid(t *testing.T) {
	acc.AccAssertSource(t, "android", emptyTestConfigs)
}

func TestAccSourceIOS(t *testing.T) {
	acc.AccAssertSource(t, "ios", emptyTestConfigs)
}

func TestAccSourceJava(t *testing.T) {
	acc.AccAssertSource(t, "java", emptyTestConfigs)
}

func TestAccSourceJavascript(t *testing.T) {
	acc.AccAssertSource(t, "javascript", emptyTestConfigs)
}

func TestAccSourceNode(t *testing.T) {
	acc.AccAssertSource(t, "node", emptyTestConfigs)
}

func TestAccSourceReactnative(t *testing.T) {
	acc.AccAssertSource(t, "reactnative", emptyTestConfigs)
}

func TestAccSourceRuby(t *testing.T) {
	acc.AccAssertSource(t, "ruby", emptyTestConfigs)
}

func TestAccSourceWebhook(t *testing.T) {
	acc.AccAssertSource(t, "webhook", emptyTestConfigs)
}

func TestAccSourceWebhookShopify(t *testing.T) {
	acc.AccAssertSource(t, "webhook_shopify", emptyTestConfigs)
}

func TestAccSourcePython(t *testing.T) {
	acc.AccAssertSource(t, "python", emptyTestConfigs)
}

func TestAccSourcePhp(t *testing.T) {
	acc.AccAssertSource(t, "php", emptyTestConfigs)
}

func TestAccSourceDotnet(t *testing.T) {
	acc.AccAssertSource(t, "dotnet", emptyTestConfigs)
}

func TestAccSourceFlutter(t *testing.T) {
	acc.AccAssertSource(t, "flutter", emptyTestConfigs)
}

func TestAccSourceCustomerio(t *testing.T) {
	acc.AccAssertSource(t, "customerio", emptyTestConfigs)
}

func TestAccSourceFacebookLeadAds(t *testing.T) {
	acc.AccAssertSource(t, "facebook_lead_ads", emptyTestConfigs)
}

func TestAccSourceAdjust(t *testing.T) {
	acc.AccAssertSource(t, "adjust", emptyTestConfigs)
}

func TestAccSourceRust(t *testing.T) {
	acc.AccAssertSource(t, "rust", emptyTestConfigs)
}
