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

func TestSourceResourceAMP(t *testing.T) {
	cmt.AssertSource(t, "amp", []configs.TestConfig{configs.EmptyTestConfig})
}

func TestSourceResourceAndroidKotlin(t *testing.T) {
	cmt.AssertSource(t, "android_kotlin", []configs.TestConfig{configs.EmptyTestConfig})
}

func TestSourceResourceIOSSwift(t *testing.T) {
	cmt.AssertSource(t, "ios_swift", []configs.TestConfig{configs.EmptyTestConfig})
}

func TestSourceResourceUnity(t *testing.T) {
	cmt.AssertSource(t, "unity", []configs.TestConfig{configs.EmptyTestConfig})
}

func TestSourceResourceAppcenter(t *testing.T) {
	cmt.AssertSource(t, "appcenter", []configs.TestConfig{configs.EmptyTestConfig})
}

func TestSourceResourceAppsflyer(t *testing.T) {
	cmt.AssertSource(t, "appsflyer", []configs.TestConfig{configs.EmptyTestConfig})
}

func TestSourceResourceAuth0(t *testing.T) {
	cmt.AssertSource(t, "auth0", []configs.TestConfig{configs.EmptyTestConfig})
}

func TestSourceResourceCanny(t *testing.T) {
	cmt.AssertSource(t, "canny", []configs.TestConfig{configs.EmptyTestConfig})
}

func TestSourceResourceCloseCRM(t *testing.T) {
	cmt.AssertSource(t, "close_crm", []configs.TestConfig{configs.EmptyTestConfig})
}

func TestSourceResourceCordial(t *testing.T) {
	cmt.AssertSource(t, "cordial", []configs.TestConfig{configs.EmptyTestConfig})
}

func TestSourceResourceExtole(t *testing.T) {
	cmt.AssertSource(t, "extole", []configs.TestConfig{configs.EmptyTestConfig})
}

func TestSourceResourceFormsort(t *testing.T) {
	cmt.AssertSource(t, "formsort", []configs.TestConfig{configs.EmptyTestConfig})
}

func TestSourceResourceGainsightPX(t *testing.T) {
	cmt.AssertSource(t, "gainsightpx", []configs.TestConfig{configs.EmptyTestConfig})
}

func TestSourceResourceIterable(t *testing.T) {
	cmt.AssertSource(t, "iterable", []configs.TestConfig{configs.EmptyTestConfig})
}

func TestSourceResourceLooker(t *testing.T) {
	cmt.AssertSource(t, "looker", []configs.TestConfig{configs.EmptyTestConfig})
}

func TestSourceResourceMailjet(t *testing.T) {
	cmt.AssertSource(t, "mailjet", []configs.TestConfig{configs.EmptyTestConfig})
}

func TestSourceResourceMailmodo(t *testing.T) {
	cmt.AssertSource(t, "mailmodo", []configs.TestConfig{configs.EmptyTestConfig})
}

func TestSourceResourceMoEngage(t *testing.T) {
	cmt.AssertSource(t, "moengage", []configs.TestConfig{configs.EmptyTestConfig})
}

func TestSourceResourceMonday(t *testing.T) {
	cmt.AssertSource(t, "monday", []configs.TestConfig{configs.EmptyTestConfig})
}

func TestSourceResourceOlark(t *testing.T) {
	cmt.AssertSource(t, "olark", []configs.TestConfig{configs.EmptyTestConfig})
}

func TestSourceResourceOrtto(t *testing.T) {
	cmt.AssertSource(t, "ortto", []configs.TestConfig{configs.EmptyTestConfig})
}

func TestSourceResourcePagerDuty(t *testing.T) {
	cmt.AssertSource(t, "pagerduty", []configs.TestConfig{configs.EmptyTestConfig})
}

func TestSourceResourcePipedream(t *testing.T) {
	cmt.AssertSource(t, "pipedream", []configs.TestConfig{configs.EmptyTestConfig})
}

func TestSourceResourceRefiner(t *testing.T) {
	cmt.AssertSource(t, "refiner", []configs.TestConfig{configs.EmptyTestConfig})
}

func TestSourceResourceRevenuecat(t *testing.T) {
	cmt.AssertSource(t, "revenuecat", []configs.TestConfig{configs.EmptyTestConfig})
}

func TestSourceResourceSatisMeter(t *testing.T) {
	cmt.AssertSource(t, "satismeter", []configs.TestConfig{configs.EmptyTestConfig})
}

func TestSourceResourceSegment(t *testing.T) {
	cmt.AssertSource(t, "segment", []configs.TestConfig{configs.EmptyTestConfig})
}

func TestSourceResourceSIGNL4(t *testing.T) {
	cmt.AssertSource(t, "signl4", []configs.TestConfig{configs.EmptyTestConfig})
}

func TestSourceResourceSlack(t *testing.T) {
	cmt.AssertSource(t, "slack", []configs.TestConfig{configs.EmptyTestConfig})
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

func TestAccSourceAMP(t *testing.T) {
	acc.AccAssertSource(t, "amp", emptyTestConfigs)
}

func TestAccSourceAndroidKotlin(t *testing.T) {
	acc.AccAssertSource(t, "android_kotlin", emptyTestConfigs)
}

func TestAccSourceIOSSwift(t *testing.T) {
	acc.AccAssertSource(t, "ios_swift", emptyTestConfigs)
}

func TestAccSourceUnity(t *testing.T) {
	acc.AccAssertSource(t, "unity", emptyTestConfigs)
}

func TestAccSourceAppcenter(t *testing.T) {
	acc.AccAssertSource(t, "appcenter", emptyTestConfigs)
}

func TestAccSourceAppsflyer(t *testing.T) {
	acc.AccAssertSource(t, "appsflyer", emptyTestConfigs)
}

func TestAccSourceAuth0(t *testing.T) {
	acc.AccAssertSource(t, "auth0", emptyTestConfigs)
}

func TestAccSourceCanny(t *testing.T) {
	acc.AccAssertSource(t, "canny", emptyTestConfigs)
}

func TestAccSourceCloseCRM(t *testing.T) {
	acc.AccAssertSource(t, "close_crm", emptyTestConfigs)
}

func TestAccSourceCordial(t *testing.T) {
	acc.AccAssertSource(t, "cordial", emptyTestConfigs)
}

func TestAccSourceExtole(t *testing.T) {
	acc.AccAssertSource(t, "extole", emptyTestConfigs)
}

func TestAccSourceFormsort(t *testing.T) {
	acc.AccAssertSource(t, "formsort", emptyTestConfigs)
}

func TestAccSourceGainsightPX(t *testing.T) {
	acc.AccAssertSource(t, "gainsightpx", emptyTestConfigs)
}

func TestAccSourceIterable(t *testing.T) {
	acc.AccAssertSource(t, "iterable", emptyTestConfigs)
}

func TestAccSourceLooker(t *testing.T) {
	acc.AccAssertSource(t, "looker", emptyTestConfigs)
}

func TestAccSourceMailjet(t *testing.T) {
	acc.AccAssertSource(t, "mailjet", emptyTestConfigs)
}

func TestAccSourceMailmodo(t *testing.T) {
	acc.AccAssertSource(t, "mailmodo", emptyTestConfigs)
}

func TestAccSourceMoEngage(t *testing.T) {
	acc.AccAssertSource(t, "moengage", emptyTestConfigs)
}

func TestAccSourceMonday(t *testing.T) {
	acc.AccAssertSource(t, "monday", emptyTestConfigs)
}

func TestAccSourceOlark(t *testing.T) {
	acc.AccAssertSource(t, "olark", emptyTestConfigs)
}

func TestAccSourceOrtto(t *testing.T) {
	acc.AccAssertSource(t, "ortto", emptyTestConfigs)
}

func TestAccSourcePagerDuty(t *testing.T) {
	acc.AccAssertSource(t, "pagerduty", emptyTestConfigs)
}

func TestAccSourcePipedream(t *testing.T) {
	acc.AccAssertSource(t, "pipedream", emptyTestConfigs)
}

func TestAccSourceRefiner(t *testing.T) {
	acc.AccAssertSource(t, "refiner", emptyTestConfigs)
}

func TestAccSourceRevenuecat(t *testing.T) {
	acc.AccAssertSource(t, "revenuecat", emptyTestConfigs)
}

func TestAccSourceSatisMeter(t *testing.T) {
	acc.AccAssertSource(t, "satismeter", emptyTestConfigs)
}

func TestAccSourceSegment(t *testing.T) {
	acc.AccAssertSource(t, "segment", emptyTestConfigs)
}

func TestAccSourceSIGNL4(t *testing.T) {
	acc.AccAssertSource(t, "signl4", emptyTestConfigs)
}

func TestAccSourceSlack(t *testing.T) {
	acc.AccAssertSource(t, "slack", emptyTestConfigs)
}
