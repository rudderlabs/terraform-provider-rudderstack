package rudderstack

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func resourceConfigCustomizeDiff(cm configs.ConfigMeta) schema.CustomizeDiffFunc {
	return func(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
		return validateResourceConfigDiff(cm, d)
	}
}

func validateResourceConfigDiff(cm configs.ConfigMeta, d *schema.ResourceDiff) error {
	if cm.SkipConfig || cm.CustomizeConfigDiff == nil {
		return nil
	}

	return cm.CustomizeConfigDiff(d)
}
