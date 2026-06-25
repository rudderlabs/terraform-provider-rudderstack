package rudderstack

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAccount() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAccountRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique identifier of the account to look up.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Human readable name of the account.",
			},
			"definition_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The account definition name (e.g. SOURCE_BIGQUERY).",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The account definition type (e.g. bigquery).",
			},
			"options": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Free-form non-secret options for the account, as key-value string pairs.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceAccountRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, ok := m.(*Client)
	if !ok {
		return diag.FromErr(fmt.Errorf("API client is not configured"))
	}

	id := d.Get("id").(string)

	acc, err := c.Accounts.Get(ctx, id)
	if err != nil {
		return diag.FromErr(fmt.Errorf("could not get account: %w", err))
	}

	d.SetId(acc.ID)

	if err := d.Set("name", acc.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("definition_name", acc.Definition.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("type", acc.Definition.Type); err != nil {
		return diag.FromErr(err)
	}

	// Unmarshal options from RawMessage into map[string]interface{}, then
	// convert all values to strings so they fit the TypeMap(TypeString) schema.
	var rawOpts map[string]interface{}
	if len(acc.Options) > 0 && string(acc.Options) != "null" {
		if err := json.Unmarshal(acc.Options, &rawOpts); err != nil {
			return diag.FromErr(fmt.Errorf("could not unmarshal account options: %w", err))
		}
	}

	optsMap := make(map[string]string, len(rawOpts))
	for k, v := range rawOpts {
		switch sv := v.(type) {
		case string:
			optsMap[k] = sv
		default:
			b, err := json.Marshal(v)
			if err != nil {
				return diag.FromErr(fmt.Errorf("could not stringify option key %q: %w", k, err))
			}
			optsMap[k] = string(b)
		}
	}

	if err := d.Set("options", optsMap); err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}
