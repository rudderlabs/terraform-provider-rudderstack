package rudderstack

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/rudderlabs/rudder-api-go/client"
	"github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func init() {
	schema.DescriptionKind = schema.StringMarkdown
}

type ConfigureClientFunc func(ctx context.Context, d *schema.ResourceData) (*Client, diag.Diagnostics)

func NewWithConfigureClientFunc(f ConfigureClientFunc) *schema.Provider {
	p := &schema.Provider{
		ConfigureContextFunc: func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
			return f(ctx, d)
		},
		Schema: map[string]*schema.Schema{
			"api_url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("RUDDERSTACK_API_URL", "https://api.rudderstack.com/v2"),
				Description: "The base URL of Rudderstack API. If not set, the provider will first look for a value in the " +
					"`RUDDERSTACK_API_URL` environmental value, and finally default to `https://api.rudderstack.com/v2` " +
					"if that is missing.",
			},
			"access_token": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("RUDDERSTACK_ACCESS_TOKEN", ""),
				Description: "The Rudderstack API access token used to authenticate you Rudderstack account. If not set, the provider " +
					"will look for that value in the `RUDDERSTACK_ACCESS_TOKEN` environmental value, " +
					"and fail with an error if that is missing.",
			},
			"workspace_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("RUDDERSTACK_WORKSPACE_ID", ""),
				Description: "The RudderStack workspace ID, required for RETL resources.",
			},
		},
		ResourcesMap: resourcesMap(),
	}

	return p
}

func New() *schema.Provider {
	return NewWithConfigureClientFunc(configureClient)
}

func resourcesMap() map[string]*schema.Resource {
	resources := map[string]*schema.Resource{
		"rudderstack_connection":        resourceConnection(),
		"rudderstack_warehouse_account": resourceWarehouseAccount(),
		"rudderstack_retl_source":       resourceRETLSource(),
		"rudderstack_retl_destination":  resourceRETLDestination(),
		"rudderstack_retl_connection":   resourceRETLConnection(),
	}

	// append sources and destinations from integration registries
	for k, v := range configs.Sources.Entries() {
		key := fmt.Sprintf("rudderstack_source_%s", k)
		resource := resourceSource(v)
		resources[key] = resource
	}
	for k, v := range configs.Destinations.Entries() {
		key := fmt.Sprintf("rudderstack_destination_%s", k)
		resource := resourceDestination(v)
		resources[key] = resource
	}
	return resources
}

func configureClient(ctx context.Context, d *schema.ResourceData) (*Client, diag.Diagnostics) {
	apiUrl := d.Get("api_url").(string)
	accessToken := d.Get("access_token").(string)

	opts := []client.Option{
		client.WithBaseURL(apiUrl),
		client.WithUserAgent("terraform-provider-rudderstack/4.2.0"),
	}

	if workspaceID, ok := d.GetOk("workspace_id"); ok && workspaceID.(string) != "" {
		opts = append(opts, client.WithWorkspaceID(workspaceID.(string)))
	}

	c, err := NewAPIClient(accessToken, opts...)
	if err != nil {
		return nil, diag.FromErr(err)
	}
	return c, diag.Diagnostics{}
}
