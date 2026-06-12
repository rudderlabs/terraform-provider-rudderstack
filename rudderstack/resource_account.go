// Package rudderstack provides the Terraform resource handler for warehouse accounts.
//
// ConfigMeta dot-notation maps HCL config.x ↔ API options.X/secret.X; secrets are
// Sensitive:true (set by the per-warehouse integration file, DEX-379); accountDefinitionName
// is immutable; no external_id this milestone; secret is never returned by the API so it
// is write-only from Terraform's perspective.

package rudderstack

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tidwall/gjson"

	"github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

func resourceAccount(cm configs.ConfigMeta) *schema.Resource {
	return &schema.Resource{
		Schema:        resourceAccountSchema(cm),
		CreateContext: resourceAccountCreate(cm),
		ReadContext:   resourceAccountRead(cm),
		UpdateContext: resourceAccountUpdate(cm),
		DeleteContext: resourceAccountDelete(cm),
		Importer: &schema.ResourceImporter{
			StateContext: resourceAccountImportState(cm),
		},
	}
}

func resourceAccountSchema(cm configs.ConfigMeta) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Human readable name of the account. The value has to be unique across all accounts.",
		},
		"created_at": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Time when the resource was created, in ISO 8601 format.",
		},
		"updated_at": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Time when the resource was last updated, in ISO 8601 format.",
		},
		"config": {
			Type:        schema.TypeList,
			Required:    true,
			MaxItems:    1,
			Description: "Account specific configuration. Check the nested block documentation for more information.",
			Elem: &schema.Resource{
				Schema: cm.ConfigSchema,
			},
		},
	}
}

// optionsOrEmpty returns the raw JSON of options, falling back to "{}".
func optionsOrEmpty(raw json.RawMessage) string {
	if len(raw) == 0 {
		return "{}"
	}
	return string(raw)
}

// splitOptionsSecret splits a combined StateToAPI result ({"options":{...},"secret":{...}})
// into separate options and secret raw JSON messages.
func splitOptionsSecret(combined string) (options, secret json.RawMessage) {
	optRaw := gjson.Get(combined, "options").Raw
	if optRaw == "" {
		optRaw = "{}"
	}
	secRaw := gjson.Get(combined, "secret").Raw
	if secRaw == "" {
		secRaw = "{}"
	}
	return json.RawMessage(optRaw), json.RawMessage(secRaw)
}

func resourceAccountCreate(cm configs.ConfigMeta) schema.CreateContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		c, ok := m.(*Client)
		if !ok {
			return diag.FromErr(fmt.Errorf("API client is not configured"))
		}

		state, err := json.Marshal(d.Get("config.0"))
		if err != nil {
			return diag.FromErr(fmt.Errorf("could not marshal config to JSON: %w", err))
		}

		combined, err := cm.StateToAPI(string(state))
		if err != nil {
			return diag.FromErr(fmt.Errorf("could not convert state to API config: %w", err))
		}

		options, secret := splitOptionsSecret(combined)

		req := &CreateAccountRequest{
			Name:                  d.Get("name").(string),
			AccountDefinitionName: cm.APIType,
			Options:               options,
			Secret:                secret,
		}

		acc, err := c.Accounts.Create(ctx, req)
		if err != nil {
			return diag.FromErr(fmt.Errorf("could not create account: %w", err))
		}

		d.SetId(acc.ID)

		return resourceAccountRead(cm)(ctx, d, m)
	}
}

func resourceAccountRead(cm configs.ConfigMeta) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		c, ok := m.(*Client)
		if !ok {
			return diag.FromErr(fmt.Errorf("API client is not configured"))
		}

		acc, err := c.Accounts.Get(ctx, d.Id())
		if err != nil {
			return diag.FromErr(fmt.Errorf("could not get account: %w", err))
		}

		// The API never returns secret; provide an empty secret so copyToState
		// leaves secret/credentials fields unset (they only write when Exists()).
		combinedForState := fmt.Sprintf(`{"options":%s,"secret":{}}`, optionsOrEmpty(acc.Options))

		stateJSON, err := cm.APIToState(combinedForState)
		if err != nil {
			return diag.FromErr(fmt.Errorf("could not convert API config to state: %w", err))
		}

		// apiProps contains only the fields the API returned (options, no secret).
		apiProps := make(map[string]interface{})
		if err := json.Unmarshal([]byte(stateJSON), &apiProps); err != nil {
			return diag.FromErr(fmt.Errorf("could not unmarshal state JSON: %w", err))
		}

		// Merge: start from API-returned fields (authoritative for all non-secret/options
		// fields), then preserve write-only secret fields from prior state — the API
		// never returns them, so without this they'd be dropped from state and cause a
		// perpetual diff.
		mergedProps := make(map[string]interface{})
		for k, v := range apiProps {
			mergedProps[k] = v
		}
		// Preserve write-only secret fields from prior state — the API never returns them,
		// so without this they'd be dropped from state and cause a perpetual diff.
		if existing, ok := d.GetOk("config"); ok {
			if list, ok := existing.([]interface{}); ok && len(list) > 0 {
				if priorMap, ok := list[0].(map[string]interface{}); ok {
					for key, sch := range cm.ConfigSchema {
						if sch.Sensitive {
							if val, ok := priorMap[key]; ok {
								mergedProps[key] = val
							}
						}
					}
				}
			}
		}

		if err := d.Set("name", acc.Name); err != nil {
			return diag.FromErr(err)
		}

		if acc.CreatedAt != nil {
			if err := d.Set("created_at", acc.CreatedAt.Format(time.RFC3339)); err != nil {
				return diag.FromErr(err)
			}
		}

		if acc.UpdatedAt != nil {
			if err := d.Set("updated_at", acc.UpdatedAt.Format(time.RFC3339)); err != nil {
				return diag.FromErr(err)
			}
		}

		if len(mergedProps) > 0 {
			if err := d.Set("config", []interface{}{mergedProps}); err != nil {
				return diag.FromErr(err)
			}
		} else {
			if err := d.Set("config", []interface{}{}); err != nil {
				return diag.FromErr(err)
			}
		}

		return diag.Diagnostics{}
	}
}

func resourceAccountUpdate(cm configs.ConfigMeta) schema.UpdateContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		c, ok := m.(*Client)
		if !ok {
			return diag.FromErr(fmt.Errorf("API client is not configured"))
		}

		state, err := json.Marshal(d.Get("config.0"))
		if err != nil {
			return diag.FromErr(fmt.Errorf("could not marshal config to JSON: %w", err))
		}

		combined, err := cm.StateToAPI(string(state))
		if err != nil {
			return diag.FromErr(fmt.Errorf("could not convert state to API config: %w", err))
		}

		options, secret := splitOptionsSecret(combined)

		req := &UpdateAccountRequest{
			Name:    d.Get("name").(string),
			Options: options,
			Secret:  secret,
		}

		acc, err := c.Accounts.Update(ctx, d.Id(), req)
		if err != nil {
			return diag.FromErr(fmt.Errorf("could not update account: %w", err))
		}

		d.SetId(acc.ID)

		return resourceAccountRead(cm)(ctx, d, m)
	}
}

func resourceAccountDelete(cm configs.ConfigMeta) schema.DeleteContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		c, ok := m.(*Client)
		if !ok {
			return diag.FromErr(fmt.Errorf("API client is not configured"))
		}

		if err := c.Accounts.Delete(ctx, d.Id()); err != nil {
			return diag.FromErr(fmt.Errorf("could not delete account: %w", err))
		}

		d.SetId("")
		return diag.Diagnostics{}
	}
}

func resourceAccountImportState(cm configs.ConfigMeta) schema.StateContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
		diagnostics := resourceAccountRead(cm)(ctx, d, m)
		if diagnostics.HasError() {
			for _, diagnostic := range diagnostics {
				if diagnostic.Severity == diag.Error {
					return nil, fmt.Errorf("could not import account: %s", diagnostic.Summary)
				}
			}
		}
		return []*schema.ResourceData{d}, nil
	}
}
