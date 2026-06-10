package rudderstack

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"

	"github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
)

// resourceAccount builds a generic, ConfigMeta-driven account resource.
// The ConfigMeta property list must map account config fields using
// `options.<path>` and `secret.<path>` keys. Secret fields should be marked
// sensitive in schema, and accountDefinitionName is sourced from cm.APIType.
// The API does not return secret values on read/import. Read preserves secrets
// from prior state to avoid perpetual drift, while import initializes
// `secret: {}` because no prior state exists.
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
			Description: "Account specific configuration. Check nested block documentation for details.",
			Elem: &schema.Resource{
				Schema: cm.ConfigSchema,
			},
		},
	}
}

func resourceAccountCreate(cm configs.ConfigMeta) schema.CreateContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		api, diags := accountsService(m)
		if diags.HasError() {
			return diags
		}

		combined, err := accountConfigStateToAPI(cm, d)
		if err != nil {
			return diag.FromErr(err)
		}

		req := &CreateAccountRequest{
			Name:                  d.Get("name").(string),
			AccountDefinitionName: cm.APIType,
			Options:               nestedJSONOrEmpty(combined, "options"),
			Secret:                nestedJSONOrEmpty(combined, "secret"),
		}

		account, err := api.Create(ctx, req)
		if err != nil {
			return diag.FromErr(fmt.Errorf("could not create account: %w", err))
		}

		d.SetId(account.ID)
		return resourceAccountRead(cm)(ctx, d, m)
	}
}

func resourceAccountRead(cm configs.ConfigMeta) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		api, diags := accountsService(m)
		if diags.HasError() {
			return diags
		}

		account, err := api.Get(ctx, d.Id())
		if err != nil {
			if errors.Is(err, ErrAccountNotFound) {
				d.SetId("")
				return nil
			}
			return diag.FromErr(err)
		}

		if err := storeAccountToState(cm, account, d); err != nil {
			return diag.FromErr(err)
		}

		return nil
	}
}

func resourceAccountUpdate(cm configs.ConfigMeta) schema.UpdateContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		api, diags := accountsService(m)
		if diags.HasError() {
			return diags
		}

		combined, err := accountConfigStateToAPI(cm, d)
		if err != nil {
			return diag.FromErr(err)
		}

		// CONTRACT-ACCT-V1 §6.1 uses replacement semantics; send all fields.
		req := &UpdateAccountRequest{
			Name:    d.Get("name").(string),
			Options: nestedJSONOrEmpty(combined, "options"),
			Secret:  nestedJSONOrEmpty(combined, "secret"),
		}

		account, err := api.Update(ctx, d.Id(), req)
		if err != nil {
			return diag.FromErr(fmt.Errorf("could not update account: %w", err))
		}

		d.SetId(account.ID)
		return resourceAccountRead(cm)(ctx, d, m)
	}
}

func resourceAccountDelete(cm configs.ConfigMeta) schema.DeleteContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		api, diags := accountsService(m)
		if diags.HasError() {
			return diags
		}

		if err := api.Delete(ctx, d.Id()); err != nil {
			if errors.Is(err, ErrAccountNotFound) {
				d.SetId("")
				return nil
			}
			return diag.FromErr(err)
		}

		d.SetId("")
		return nil
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

func accountsService(m interface{}) (accountsAPI, diag.Diagnostics) {
	c, ok := m.(*Client)
	if !ok {
		return nil, diag.FromErr(fmt.Errorf("API client is not configured"))
	}

	api := c.accountsClient()
	if api == nil {
		return nil, diag.FromErr(fmt.Errorf("accounts API client is not configured"))
	}

	return api, nil
}

func accountConfigStateToAPI(cm configs.ConfigMeta, d *schema.ResourceData) (string, error) {
	state, err := json.Marshal(d.Get("config.0"))
	if err != nil {
		return "", err
	}

	return cm.StateToAPI(string(state))
}

func nestedJSONOrEmpty(raw, path string) json.RawMessage {
	value := gjson.Get(raw, path)
	if value.Exists() && value.IsObject() {
		return json.RawMessage(value.Raw)
	}
	return json.RawMessage("{}")
}

func storeAccountToState(cm configs.ConfigMeta, account *Account, d *schema.ResourceData) error {
	d.SetId(account.ID)

	if err := d.Set("name", account.Name); err != nil {
		return err
	}

	if account.CreatedAt != nil {
		if err := d.Set("created_at", account.CreatedAt.Format(time.RFC3339)); err != nil {
			return err
		}
	}

	if account.UpdatedAt != nil {
		if err := d.Set("updated_at", account.UpdatedAt.Format(time.RFC3339)); err != nil {
			return err
		}
	}

	combined, err := sjson.SetRaw("{}", "options", string(objectJSONOrEmpty(account.Options)))
	if err != nil {
		return err
	}
	combined, err = sjson.SetRaw(combined, "secret", string(accountSecretsFromState(cm, d)))
	if err != nil {
		return err
	}

	state, err := cm.APIToState(combined)
	if err != nil {
		return err
	}

	properties := make(map[string]interface{})
	if err := json.Unmarshal([]byte(state), &properties); err != nil {
		return err
	}

	if len(properties) > 0 {
		if err := d.Set("config", []interface{}{properties}); err != nil {
			return err
		}
	} else {
		if err := d.Set("config", []interface{}{}); err != nil {
			return err
		}
	}

	return nil
}

func accountSecretsFromState(cm configs.ConfigMeta, d *schema.ResourceData) json.RawMessage {
	currentConfig, ok := d.GetOk("config.0")
	if !ok || currentConfig == nil {
		return json.RawMessage("{}")
	}

	currentState, err := json.Marshal(currentConfig)
	if err != nil {
		return json.RawMessage("{}")
	}

	combined, err := cm.StateToAPI(string(currentState))
	if err != nil {
		return json.RawMessage("{}")
	}

	return nestedJSONOrEmpty(combined, "secret")
}

func objectJSONOrEmpty(raw json.RawMessage) json.RawMessage {
	value := gjson.ParseBytes(raw)
	if value.IsObject() {
		return json.RawMessage(value.Raw)
	}
	return json.RawMessage("{}")
}
