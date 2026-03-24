package rudderstack

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/rudderlabs/rudder-api-go/client"
)

func resourceWarehouseAccount() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Account display name.",
			},
			"account_definition_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Warehouse type, e.g. \"bigquery\".",
			},
			"options": {
				Type:        schema.TypeMap,
				Required:    true,
				Sensitive:   true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Warehouse-specific credentials and options.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Time created, ISO 8601.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Time last updated, ISO 8601.",
			},
		},
		CreateContext: resourceWarehouseAccountCreate,
		ReadContext:   resourceWarehouseAccountRead,
		UpdateContext: resourceWarehouseAccountUpdate,
		DeleteContext: resourceWarehouseAccountDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceWarehouseAccountImportState,
		},
	}
}

func resourceWarehouseAccountCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, ok := m.(*Client)
	if !ok {
		return diag.FromErr(fmt.Errorf("API client is not configured"))
	}

	input := &client.AccountCreateInput{}
	populateAccountCreateInputFromState(input, d)

	account, err := c.Accounts.Create(ctx, input)
	if err != nil {
		return diag.FromErr(fmt.Errorf("could not create warehouse account: %w", err))
	}

	d.SetId(account.ID)

	return resourceWarehouseAccountRead(ctx, d, m)
}

func resourceWarehouseAccountRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, ok := m.(*Client)
	if !ok {
		return diag.FromErr(fmt.Errorf("API client is not configured"))
	}

	id := d.Id()

	account, err := c.Accounts.Get(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := storeWarehouseAccountToState(account, d); err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}

func resourceWarehouseAccountUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, ok := m.(*Client)
	if !ok {
		return diag.FromErr(fmt.Errorf("API client is not configured"))
	}

	input := &client.AccountUpdateInput{}
	populateAccountUpdateInputFromState(input, d)

	account, err := c.Accounts.Update(ctx, d.Id(), input)
	if err != nil {
		return diag.FromErr(fmt.Errorf("could not update warehouse account: %w", err))
	}

	d.SetId(account.ID)

	return resourceWarehouseAccountRead(ctx, d, m)
}

func resourceWarehouseAccountDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

func resourceWarehouseAccountImportState(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	diagnostics := resourceWarehouseAccountRead(ctx, d, m)
	if diagnostics.HasError() {
		for _, diagnostic := range diagnostics {
			if diagnostic.Severity == diag.Error {
				return nil, fmt.Errorf("could not import warehouse account: %s", diagnostic.Summary)
			}
		}
	}
	return []*schema.ResourceData{d}, nil
}

func populateAccountCreateInputFromState(input *client.AccountCreateInput, d *schema.ResourceData) {
	input.Name = d.Get("name").(string)
	input.AccountDefinitionName = d.Get("account_definition_name").(string)

	if opts, ok := d.GetOk("options"); ok {
		raw := opts.(map[string]interface{})
		options := make(map[string]interface{}, len(raw))
		for k, v := range raw {
			options[k] = v
		}
		input.Options = options
	}
}

func populateAccountUpdateInputFromState(input *client.AccountUpdateInput, d *schema.ResourceData) {
	input.Name = d.Get("name").(string)

	if opts, ok := d.GetOk("options"); ok {
		raw := opts.(map[string]interface{})
		options := make(map[string]interface{}, len(raw))
		for k, v := range raw {
			options[k] = v
		}
		input.Options = options
	}
}

func storeWarehouseAccountToState(account *client.Account, d *schema.ResourceData) error {
	d.SetId(account.ID)

	if err := d.Set("name", account.Name); err != nil {
		return err
	}
	if account.Definition != nil {
		if err := d.Set("account_definition_name", account.Definition.Type); err != nil {
			return err
		}
	}
	// Convert map[string]interface{} to map[string]string for Terraform state
	if account.Options != nil {
		stringOpts := make(map[string]string, len(account.Options))
		for k, v := range account.Options {
			stringOpts[k] = fmt.Sprintf("%v", v)
		}
		if err := d.Set("options", stringOpts); err != nil {
			return err
		}
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
