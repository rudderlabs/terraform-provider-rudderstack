package rudderstack

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rudderlabs/rudder-api-go/client"
	"github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/accounts"
)

func resourceAccount(cm accounts.AccountConfigMeta) *schema.Resource {
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

func resourceAccountSchema(cm accounts.AccountConfigMeta) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Human readable name of the account.",
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
			Type:     schema.TypeList,
			Required: true,
			Description: "Account specific configuration. Check the nested block documenation " +
				"for more information.",
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: cm.ConfigSchema,
			},
		},
	}
}

func resourceAccountRead(cm accounts.AccountConfigMeta) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		c, ok := m.(*Client)
		if !ok {
			return diag.FromErr(fmt.Errorf("API client is not configured"))
		}

		id := d.Id()

		account, err := c.Accounts.Get(ctx, id)
		if err != nil {
			return diag.FromErr(err)
		}

		storeAccountToState(cm, account, d)

		return diag.Diagnostics{}
	}
}

func resourceAccountCreate(cm accounts.AccountConfigMeta) schema.CreateContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		c, ok := m.(*Client)
		if !ok {
			return diag.FromErr(fmt.Errorf("API client is not configured"))
		}

		accountWithSecret := &client.AccountWithSecret{}
		err := populateAccountFromState(cm, accountWithSecret, d)
		if err != nil {
			return diag.FromErr(err)
		}

		account, err := c.Accounts.Create(ctx, accountWithSecret)
		if err != nil {
			return diag.FromErr(fmt.Errorf("could not create account: %w", err))
		}

		d.SetId(account.ID)

		storeAccountToState(cm, account, d)
		storeOptionsAndSecretToState(cm, account, accountWithSecret.Secret, d)
		return diag.Diagnostics{}
		// return resourceAccountRead(cm)(ctx, d, m)
	}
}

func resourceAccountUpdate(cm accounts.AccountConfigMeta) schema.UpdateContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		c, ok := m.(*Client)
		if !ok {
			return diag.FromErr(fmt.Errorf("API client is not configured"))
		}

		accountWithSecret := &client.AccountWithSecret{}
		err := populateAccountFromState(cm, accountWithSecret, d)
		if err != nil {
			return diag.FromErr(err)
		}

		account, err := c.Accounts.Update(ctx, accountWithSecret)
		if err != nil {
			return diag.FromErr(fmt.Errorf("could not update account: %w", err))
		}

		d.SetId(account.ID)

		storeAccountToState(cm, account, d)
		storeOptionsAndSecretToState(cm, account, accountWithSecret.Secret, d)
		return diag.Diagnostics{}
		// return resourceAccountRead(cm)(ctx, d, m)
	}
}

func resourceAccountDelete(cm accounts.AccountConfigMeta) schema.DeleteContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		c, ok := m.(*Client)
		if !ok {
			return diag.FromErr(fmt.Errorf("API client is not configured"))
		}

		if err := c.Accounts.Delete(ctx, d.Id()); err != nil {
			return diag.FromErr(err)
		}

		d.SetId("")
		return diag.Diagnostics{}
	}
}

func resourceAccountImportState(cm accounts.AccountConfigMeta) schema.StateContextFunc {
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

func populateAccountFromState(cm accounts.AccountConfigMeta, account *client.AccountWithSecret, d *schema.ResourceData) error {
	account.ID = d.Id()
	account.Type = cm.APIType
	account.Category = string(cm.Category)
	account.Name = d.Get("name").(string)

	if c := d.Get("config.0"); c != nil {
		state, err := json.Marshal(c)
		if err != nil {
			return err
		}

		options, err := cm.StateToAPI(string(state))
		if err != nil {
			return err
		}

		secret, err := cm.SecretStateToAPI(string(state))
		if err != nil {
			return err
		}

		account.Options = json.RawMessage(options)
		account.Secret = json.RawMessage(secret)
	}

	return nil
}

func storeAccountToState(cm accounts.AccountConfigMeta, account *client.Account, d *schema.ResourceData) error {
	d.SetId(account.ID)
	if err := d.Set("name", account.Name); err != nil {
		return err
	}
	if account.CreatedAt != nil {
		createdAt := account.CreatedAt.Format(time.RFC3339)
		if err := d.Set("created_at", createdAt); err != nil {
			return err
		}
	}
	if account.UpdatedAt != nil {
		updatedAt := account.UpdatedAt.Format(time.RFC3339)
		if err := d.Set("updated_at", updatedAt); err != nil {
			return err
		}
	}

	return nil
}

func storeOptionsAndSecretToState(cm accounts.AccountConfigMeta, account *client.Account, secret json.RawMessage, d *schema.ResourceData) error {
	optionsState, err := cm.APIToState(string(account.Options))
	if err != nil {
		return err
	}

	secretState, err := cm.SecretAPIToState(string(secret))
	if err != nil {
		return err
	}

	optionsProperties := make(map[string]interface{})
	if err := json.Unmarshal([]byte(optionsState), &optionsProperties); err != nil {
		return err
	}

	secretProperties := make(map[string]interface{})
	if err := json.Unmarshal([]byte(secretState), &secretProperties); err != nil {
		return err
	}

	mergedProperties := make(map[string]interface{})
	for k, v := range optionsProperties {
		mergedProperties[k] = v
	}
	for k, v := range secretProperties {
		mergedProperties[k] = v
	}

	if len(mergedProperties) > 0 {
		if err := d.Set("config", []interface{}{mergedProperties}); err != nil {
			return err
		}
	} else {
		if err := d.Set("config", []interface{}{map[string]interface{}{}}); err != nil {
			return err
		}
	}

	return nil
}
