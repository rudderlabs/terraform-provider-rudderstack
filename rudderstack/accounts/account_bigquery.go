package accounts

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	c "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func init() {
	properties := []c.ConfigProperty{
		c.Simple("options.project", "project"),
		c.Simple("options.location", "location", c.SkipZeroValue),
		c.Simple("secret.credentials", "credentials"),
	}
	cfgSchema := map[string]*schema.Schema{
		"project": {
			Type:             schema.TypeString,
			Required:         true,
			Description:      "GCP project ID where the BigQuery dataset lives.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|^[a-z][a-z0-9.:-]{4,28}[a-z0-9]$"),
		},
		"location": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "BigQuery dataset location (e.g. US, EU).",
		},
		"credentials": {
			Type:             schema.TypeString,
			Required:         true,
			Sensitive:        true,
			Description:      "Full contents of the GCP service-account key JSON file (downloaded from the GCP console). Load it with Terraform's file() function instead of inlining the multi-line key.",
			ValidateDiagFunc: c.StringMatchesRegexp("(^env[.].+)|[\\s\\S]+"),
		},
	}
	c.Accounts.Register("bigquery", c.ConfigMeta{
		APIType:      "SOURCE_BIGQUERY",
		Properties:   properties,
		ConfigSchema: cfgSchema,
	})
}
