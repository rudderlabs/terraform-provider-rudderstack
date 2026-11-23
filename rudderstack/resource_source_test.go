package rudderstack

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func TestResourceSourceCustomizeConfigDiff(t *testing.T) {
	configSchema := map[string]*schema.Schema{
		"value": {Type: schema.TypeString, Optional: true},
	}

	tests := []struct {
		name        string
		cm          configs.ConfigMeta
		expectError *regexp.Regexp
	}{
		{
			name: "hook error is returned during planning",
			cm: configs.ConfigMeta{
				APIType:      "TEST",
				ConfigSchema: configSchema,
				CustomizeConfigDiff: func(_ *schema.ResourceDiff) error {
					return errors.New("source config rejected")
				},
			},
			expectError: regexp.MustCompile(`source config rejected`),
		},
		{
			name: "source without hook is unaffected",
			cm: configs.ConfigMeta{
				APIType:      "TEST",
				ConfigSchema: configSchema,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: map[string]func() (*schema.Provider, error){
					"rudderstack": func() (*schema.Provider, error) {
						return &schema.Provider{
							ConfigureContextFunc: func(_ context.Context, _ *schema.ResourceData) (interface{}, diag.Diagnostics) {
								return nil, nil
							},
							ResourcesMap: map[string]*schema.Resource{
								"rudderstack_source_test": resourceSource(tt.cm),
							},
						}, nil
					},
				},
				Steps: []resource.TestStep{
					{
						PlanOnly:           true,
						ExpectNonEmptyPlan: tt.expectError == nil,
						ExpectError:        tt.expectError,
						Config: `
							resource "rudderstack_source_test" "example" {
								name = "test-source"
								config { value = "configured" }
							}
						`,
					},
				},
			})
		})
	}
}
