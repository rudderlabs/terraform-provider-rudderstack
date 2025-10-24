package destinations

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func init() {
	supportedSourceTypes := []string{"web", "android", "ios", "unity", "reactnative", "flutter", "cordova", "amp", "cloud", "warehouse", "shopify"}
	commonProperties, commonSchema := GetCommonConfigMeta(supportedSourceTypes)

	properties := []c.ConfigProperty{
		c.Simple("token", "token"),
		c.Simple("apiSecret", "api_secret", c.SkipZeroValue),
		c.Simple("dataResidency", "data_residency"),
		c.Simple("people", "people", c.SkipZeroValue),
		c.Simple("setAllTraitsByDefault", "set_all_traits_by_default", c.SkipZeroValue),
		c.Simple("consolidatedPageCalls", "consolidated_page_calls"),
		c.Simple("trackCategorizedPages", "track_categorized_pages", c.SkipZeroValue),
		c.Simple("trackNamedPages", "track_named_pages", c.SkipZeroValue),
		c.Simple("sourceName", "source_name", c.SkipZeroValue),
		c.Simple("crossSubdomainCookie", "cross_subdomain_cookie", c.SkipZeroValue),
		c.Simple("secureCookie", "secure_cookie", c.SkipZeroValue),
		c.ArrayWithStrings("superProperties", "property", "super_properties"),
		c.ArrayWithStrings("peopleProperties", "property", "people_properties"),
		c.ArrayWithStrings("eventIncrements", "property", "event_increments"),
		c.ArrayWithStrings("propIncrements", "property", "prop_increments"),
		c.ArrayWithStrings("groupKeySettings", "groupKey", "group_key_settings"),
		c.Simple("useNativeSDK.web", "use_native_sdk.0.web"),
		c.Discriminator("eventFilteringOption", c.DiscriminatorValues{
			"event_filtering.0.whitelist": "whitelistedEvents",
			"event_filtering.0.blacklist": "blacklistedEvents",
		}),
		c.ArrayWithStrings("whitelistedEvents", "eventName", "event_filtering.0.whitelist"),
		c.ArrayWithStrings("blacklistedEvents", "eventName", "event_filtering.0.blacklist"),
		c.Simple("useNewMapping", "use_new_mapping", c.SkipZeroValue),
		c.Simple("identityMergeApi", "identity_merge_api"),
		c.Simple("useUserDefinedPageEventName", "use_user_defined_page_event_name"),
		c.Simple("userDefinedPageEventTemplate", "user_defined_page_event_template"),
		c.Simple("useUserDefinedScreenEventName", "use_user_defined_screen_event_name"),
		c.Simple("userDefinedScreenEventTemplate", "user_defined_screen_event_template"),
		c.Simple("dropTraitsInTrackEvent", "drop_traits_in_track_event"),
		c.Simple("strictMode", "strict_mode"),
		c.ArrayWithStrings("setOnceProperties", "property", "set_once_properties"),
		c.ArrayWithStrings("unionProperties", "property", "union_properties"),
		c.ArrayWithStrings("appendProperties", "property", "append_properties"),
		c.Simple("userDeletionApi", "user_deletion_api"),
		c.Simple("gdprApiToken", "gdpr_api_token", c.SkipZeroValue),
		c.Simple("sessionReplayPercentage.web", "session_replay_percentage.0.web"),
		c.Simple("consentManagement.web", "consent_management.0.web"),
		c.Simple("consentManagement.android", "consent_management.0.android"),
		c.Simple("consentManagement.ios", "consent_management.0.ios"),
		c.Simple("consentManagement.unity", "consent_management.0.unity"),
		c.Simple("consentManagement.reactnative", "consent_management.0.react_native"),
		c.Simple("consentManagement.flutter", "consent_management.0.flutter"),
		c.Simple("consentManagement.cordova", "consent_management.0.cordova"),
		c.Simple("consentManagement.shopify", "consent_management.0.shopify"),
		c.Simple("consentManagement.cloud", "consent_management.0.cloud"),
		c.Simple("consentManagement.warehouse", "consent_management.0.warehouse"),
		c.Simple("persistenceName", "persistence_name"),
		c.Simple("persistenceType", "persistence_type"),
	}

	properties = append(properties, commonProperties...)

	schema := map[string]*schema.Schema{
		"token": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "Mixpanel API Token",
			Sensitive:        true,
			ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
		},
		"api_secret": {
			Type:             schema.TypeString,
			Optional:         true,
			Sensitive:        true,
			Description:      "Mixpanel API secret",
			ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
		},
		"data_residency": {
			Type:             schema.TypeString,
			Required:         true,
			Description:      "Mixpanel Server region either us/eu",
			ValidateDiagFunc: c.StringMatchesRegexp("^(us|eu|in)$"),
		},
		"people": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Boolean flag to send all of your identify calls to Mixpanel's People feature ",
		},
		"set_all_traits_by_default": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "While this is checked, our integration automatically sets all traits on identify calls as super properties and people properties if Mixpanel People is checked as well.",
		},
		"consolidated_page_calls": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "This will track Loaded a Page events to Mixpanel for all page method calls. We enable this by default as it's how Mixpanel suggests sending these calls.",
		},
		"track_categorized_pages": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "This will track events to Mixpanel for page method calls that have a category associated with them. For example page('Docs', 'Index') would translate to Viewed Docs Index Page.",
		},
		"track_named_pages": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "This will track events to Mixpanel for page method calls that have a name associated with them. For example page('Signup') would translate to Viewed Signup Page.",
		},
		"source_name": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "This value, if it's not blank, will be sent as rudderstack_source_name to Mixpanel for every event/page/screen call.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
		},
		"cross_subdomain_cookie": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "This will allow the Mixpanel cookie to persist between different pages of your application.",
		},
		"secure_cookie": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "This will mark the Mixpanel cookie as secure, meaning it will only be transmitted over https.",
		},
		"super_properties": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Property to send as super Properties.",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"people_properties": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Traits to set as People Properties.",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"event_increments": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Events to increment in People.",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"prop_increments": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Properties to increment in People",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"group_key_settings": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Group Key",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"use_native_sdk": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "Enable this setting to send the events via the device mode.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"web": {
						Type:     schema.TypeBool,
						Optional: true,
					},
				},
			},
		},
		"event_filtering": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "With this option, you can determine which events are blocked or allowed to flow through to Mixpanel.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"whitelist": {
						Type:         schema.TypeList,
						Optional:     true,
						Description:  "Enter the event names to be allowlisted.",
						ExactlyOneOf: []string{"config.0.event_filtering.0.whitelist", "config.0.event_filtering.0.blacklist"},
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"blacklist": {
						Type:         schema.TypeList,
						Optional:     true,
						Description:  "Enter the event names to be denylisted.",
						ExactlyOneOf: []string{"config.0.event_filtering.0.whitelist", "config.0.event_filtering.0.blacklist"},
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
				},
			},
		},
		"use_new_mapping": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "This value is true by default and when this flag is enabled, camel case fields are mapped to snake case fields while sending to Mixpanel. Please refer to https://www.rudderstack.com/docs/destinations/streaming-destinations/mixpanel/#connection-settings for more details.",
		},
		"connection_mode": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "Enable this setting to send the events via the cloud mode.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"web": {
						Type:             schema.TypeString,
						Optional:         true,
						ValidateDiagFunc: c.StringMatchesRegexp("^(cloud|device)$"),
					},
					"android": {
						Type:             schema.TypeString,
						Optional:         true,
						ValidateDiagFunc: c.StringMatchesRegexp("^(cloud)$"),
					},
					"ios": {
						Type:             schema.TypeString,
						Optional:         true,
						ValidateDiagFunc: c.StringMatchesRegexp("^(cloud)$"),
					},
					"unity": {
						Type:             schema.TypeString,
						Optional:         true,
						ValidateDiagFunc: c.StringMatchesRegexp("^(cloud)$"),
					},
					"reactnative": {
						Type:             schema.TypeString,
						Optional:         true,
						ValidateDiagFunc: c.StringMatchesRegexp("^(cloud)$"),
					},
					"flutter": {
						Type:             schema.TypeString,
						Optional:         true,
						ValidateDiagFunc: c.StringMatchesRegexp("^(cloud)$"),
					},
					"cordova": {
						Type:             schema.TypeString,
						Optional:         true,
						ValidateDiagFunc: c.StringMatchesRegexp("^(cloud)$"),
					},
					"shopify": {
						Type:             schema.TypeString,
						Optional:         true,
						ValidateDiagFunc: c.StringMatchesRegexp("^(cloud)$"),
					},
					"cloud": {
						Type:             schema.TypeString,
						Optional:         true,
						ValidateDiagFunc: c.StringMatchesRegexp("^(cloud)$"),
					},
					"warehouse": {
						Type:             schema.TypeString,
						Optional:         true,
						ValidateDiagFunc: c.StringMatchesRegexp("^(cloud)$"),
					},
				},
			},
		},
		"identity_merge_api": {
			Type:             schema.TypeString,
			Optional:         true,
			Default:          "original",
			Description:      "Mixpanel Identity Merge types",
			ValidateDiagFunc: c.StringMatchesRegexp("^(simplified|original)$"),
		},
		"use_user_defined_page_event_name": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Boolean flag to use user-defined page event names",
		},
		"user_defined_page_event_template": {
			Type:             schema.TypeString,
			Optional:         true,
			Default:          "Viewed {{ category }} {{ name }} page",
			Description:      "Template for user-defined page event names",
			ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,200})$"),
		},
		"use_user_defined_screen_event_name": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Boolean flag to use user-defined screen event names",
		},
		"user_defined_screen_event_template": {
			Type:             schema.TypeString,
			Optional:         true,
			Default:          "Viewed {{ category }} {{ name }} screen",
			Description:      "Template for user-defined screen event names",
			ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,200})$"),
		},
		"drop_traits_in_track_event": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Boolean flag to drop traits in track event calls",
		},
		"strict_mode": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Boolean flag to enable strict mode",
		},
		"set_once_properties": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Properties to set only once",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"union_properties": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Properties to union",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"append_properties": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Properties to append",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"user_deletion_api": {
			Type:             schema.TypeString,
			Optional:         true,
			Default:          "engage",
			ValidateDiagFunc: c.StringMatchesRegexp("^(engage|task)$"),
		},
		"gdpr_api_token": {
			Type:             schema.TypeString,
			Optional:         true,
			Sensitive:        true,
			ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
		},
		"session_replay_percentage": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "Percentage of SDK initializations that will qualify for replay data capture",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"web": {
						Type:             schema.TypeString,
						Optional:         true,
						ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(100|[1-9]?[0-9])$"),
					},
				},
			},
		},
		"persistence_name": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "Mixpanel Persistence Name",
			ValidateDiagFunc: c.StringMatchesRegexp("^(none|cookie|localStorage)$"),
		},
		"persistence_type": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "Mixpanel Persistence Type",
			ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{0,100})$"),
		},
	}

	for key, value := range commonSchema {
		schema[key] = value
	}

	c.Destinations.Register("mixpanel", c.ConfigMeta{
		APIType:      "MP",
		Properties:   properties,
		ConfigSchema: schema,
	})
}
