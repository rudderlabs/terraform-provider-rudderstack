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
			"host": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("RUDDERSTACK_API_HOST", "https://api.rudderstack.com/v2"),
				Description: "The hostname of Rudderstack API interface",
			},
			"access_token": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("RUDDERSTACK_ACCESS_TOKEN", ""),
				Description: "The API access token used to authenticate you Rudderstack account",
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
	resources := map[string]*schema.Resource{}
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
	host := d.Get("host").(string)
	accessToken := d.Get("access_token").(string)
	client, err := NewAPIClient(accessToken, client.WithBaseURL(host))
	if err != nil {
		return nil, diag.FromErr(err)
	}
	return client, diag.Diagnostics{}
}
