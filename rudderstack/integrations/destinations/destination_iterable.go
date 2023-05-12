package destinations

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func init() {
	c.Destinations.Register("iterable", c.ConfigMeta{
		APIType: "ITERABLE",
		Properties: []c.ConfigProperty{
			c.Simple("apiKey", "api_key", c.SkipZeroValue),
			c.Simple("mapToSingleEvent", "map_to_single_event"),
			c.Simple("trackAllPages", "track_all_pages", c.SkipZeroValue),
			c.Simple("trackCategorisedPages", "track_categorized_pages"),
			c.Simple("trackNamedPages", "track_named_pages"),
			c.Simple("useNativeSDK.web", "use_native_sdk.0.web"),
			c.Simple("initialisationIdentifier.web", "initialisation_identifier.0.web", c.SkipZeroValue),
			c.ArrayWithStrings("getInAppEventMapping.web", "eventName", "get_in_app_event_mapping.0.web"),
			c.ArrayWithStrings("purchaseEventMapping.web", "eventName", "purchase_event_mapping.0.web"),
			c.Simple("sendTrackForInapp.web", "send_track_for_inapp.0.web",),
			c.Simple("animationDuration.web", "animation_duration.0.web"),
			c.Simple("displayInterval.web", "display_interval.0.web"),
			c.Simple("onOpenScreenReaderMessage.web", "on_open_screen_reader_message.0.web"),
			c.Simple("onOpenNodeToTakeFocus.web", "on_open_node_to_take_focus.0.web"),
			c.Simple("packageName.web", "package_name.0.web"),
			c.Simple("rightOffset.web", "right_offset.0.web"),
			c.Simple("topOffset.web", "top_offset.0.web"),
			c.Simple("bottomOffset.web", "bottom_offset.0.web"),
			c.Simple("handleLinks.web", "handle_links.0.web"),
			c.Simple("closeButtonColor.web", "close_button_color.0.web"),
			c.Simple("closeButtonSize.web", "close_button_size.0.web"),
			c.Simple("closeButtonColorTopOffset.web", "close_button_color_top_offset.0.web"),
			c.Simple("closeButtonColorSideOffset.web", "close_button_color_side_offset.0.web"),
			c.Simple("iconPath.web", "icon_path.0.web"),
			c.Simple("isRequiredToDismissMessage.web", "is_required_to_dismiss_message.0.web"),
			c.Simple("closeButtonPosition.web", "close_button_position.0.web"),
			c.ArrayWithStrings("oneTrustCookieCategories", "oneTrustCookieCategory", "onetrust_cookie_categories"),
		},
		ConfigSchema: map[string]*schema.Schema{
			"api_key": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "Enter your Iterable Api Key",
				ValidateDiagFunc: c.StringMatchesRegexp("(^\\{\\{.*\\|\\|(.*)\\}\\}$)|(^env[.].+)|^(.{1,100})$"),
			},
			"map_to_single_event": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Map All Pages to Single Event Name",
			},
			"track_all_pages": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Track All Pages",
			},
			"track_categorized_pages": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Track Categorised Pages",
			},
			"track_named_pages": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Track Named Pages",
			},
			"use_native_sdk": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Enable this setting to send events to Iterable via the device mode.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"web": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
					},
				},
			},
			"initialisation_identifier": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Identifier to identify a user over a session",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"web": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"get_in_app_event_mapping": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Mapping to trigger the getInApp messages",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"web": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"purchase_event_mapping": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Mapping to trigger the getInApp messages",
				ConfigMode:  schema.SchemaConfigModeAttr,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"web": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"send_track_for_inapp": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Trigger a track event for web in-app push",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"web": {
							Type:     schema.TypeBool,
							Optional: true,
							Default: false,
						},
					},
				},
			},
			"animation_duration": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Time (in ms) for messages to animate in and out.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"web": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"display_interval": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Wait time for next message",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"web": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"on_open_screen_reader_message": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Screen Reader Text",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"web": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"on_open_node_to_take_focus": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Focus Element",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"web": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"package_name": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Package Name",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"web": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"right_offset": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Space (px or %) between screen right & messages.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"web": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"top_offset": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Space (px or %) between screen top & messages.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"web": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"bottom_offset": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Space (px or %) between screen bottom & messages.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"web": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"handle_links": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Control how to open links.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"web": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"close_button_color": {
				Type:            schema.TypeList,
				MaxItems:        1,
				Optional:        true,
				Description:     "Color of Close button",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"web": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"close_button_size": {
				Type:            schema.TypeList,
				MaxItems:    	 1,
				Optional:        true,
				Description:     "Size of Close button",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"web": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"close_button_color_top_offset": {
				Type:            schema.TypeList,
				MaxItems:    	 1,
				Optional:        true,
				Description:     "Space between button & container top",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"web": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"close_button_color_side_offset": {
				Type:            schema.TypeList,
				MaxItems:    	 1,
				Optional:        true,
				Description:     "Space between button & container side",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"web": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"icon_path": {
				Type:            schema.TypeList,
				MaxItems:    	 1,
				Optional:        true,
				Description:     "Custom pathname",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"web": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"is_required_to_dismiss_message": {	
				Type:            schema.TypeList,
				MaxItems:    	 1,
				Optional:        true,
				Description:     "Prevent user dismissing in-app message by clicking outside message",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"web": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			"close_button_position": {
				Type:            schema.TypeList,
				MaxItems:    	 1,
				Optional:        true,
				Description:     "Position",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"web": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"onetrust_cookie_categories": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Specify the OneTrust category name for mapping the OneTrust consent settings to RudderStack's consent purposes.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	})
}