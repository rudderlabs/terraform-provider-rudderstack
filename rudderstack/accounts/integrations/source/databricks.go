package source

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	a "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/accounts"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func init() {
	err := a.Accounts.Register("databricks", a.AccountConfigMeta{
		Category: a.CategorySource,
		ConfigMeta: c.ConfigMeta{
			APIType: "databricks",
			ConfigSchema: map[string]*schema.Schema{
				"host": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The server hostname.",
				},
				"port": {
					Type:        schema.TypeInt,
					Optional:    true,
					Default:     443,
					Description: "The port number.",
				},
				"path": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The HTTP path.",
				},
				"catalog": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "The name of your Unity catalog.",
				},
				"token": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The personal access token.",
					Sensitive:   true,
				},
			},
			Properties: []c.ConfigProperty{
				c.Simple("host", "host"),
				c.Simple("port", "port"),
				c.Simple("path", "path"),
				c.Simple("catalog", "catalog", c.SkipZeroValue),
			},
		},
		Secret: []c.ConfigProperty{
			c.Simple("token", "token"),
		},
	})

	if err != nil {
		panic(err)
	}
}
