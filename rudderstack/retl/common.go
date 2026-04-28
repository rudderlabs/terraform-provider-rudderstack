// Package retl holds the standalone RETL source resources
// (rudderstack_retl_source_model / _table / _s3_table). It is deliberately
// separate from rudderstack/integrations, which is reserved for the
// ConfigMeta-registry integrations. RETL sources are registered directly in
// rudderstack/provider.go via the exported Resource* constructors.
package retl

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/rudderlabs/rudder-iac/api/client"
	"github.com/rudderlabs/rudder-iac/api/client/retl"
)

// Service is the subset of the upstream retl.RETLStore the RETL source
// resources need. Kept narrow so tests can mock it cheaply.
type Service interface {
	CreateRetlSource(ctx context.Context, source *retl.RETLSourceCreateRequest) (*retl.RETLSource, error)
	GetRetlSource(ctx context.Context, id string) (*retl.RETLSource, error)
	UpdateRetlSource(ctx context.Context, id string, source *retl.RETLSourceUpdateRequest) (*retl.RETLSource, error)
	DeleteRetlSource(ctx context.Context, id string) error
}

// ClientProvider exposes the RETL service from the provider's configured
// client. Defined here (not in the rudderstack package) so this package does
// not need to import rudderstack — which would create a circular dependency
// since rudderstack imports this package to register its resources.
type ClientProvider interface {
	RETLSourcesClient() Service
}

// typeAdapter parameterises the per-resource bits that vary between the three
// RETL source resources: the SourceType sent to the API, how the schema's
// `config` block maps to the request config JSON, and how the response config
// JSON is read back into state.
type typeAdapter struct {
	sourceType retl.SourceType
	// fixedSourceDefinitionName, if non-empty, overrides the user-supplied
	// `source_definition_name` (used by the s3_table resource).
	fixedSourceDefinitionName string
	configSchema              func() map[string]*schema.Schema
	marshalConfig             func(map[string]any) (retl.RETLConfig, error)
	unmarshalConfig           func(retl.RETLConfig) ([]map[string]interface{}, error)
}

// buildResource constructs a *schema.Resource for one of the RETL source
// resources, given the per-type adapter. The CRUD implementations are shared.
func buildResource(adapter typeAdapter, description string) *schema.Resource {
	resourceSchema := map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Human-readable name of the RETL source.",
		},
		"account_id": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "ID of the RudderStack account used to connect to the data source.",
		},
		"enabled": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "Whether the source is enabled.",
		},
		"config": {
			Type:     schema.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: adapter.configSchema(),
			},
		},
		"created_at": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"updated_at": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}

	if adapter.fixedSourceDefinitionName == "" {
		resourceSchema["source_definition_name"] = &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Type of the data source (e.g. snowflake, bigquery, postgres, redshift).",
		}
	}

	return &schema.Resource{
		Description:   description,
		Schema:        resourceSchema,
		CreateContext: makeCreate(adapter),
		ReadContext:   makeRead(adapter),
		UpdateContext: makeUpdate(adapter),
		DeleteContext: deleteSource,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// service extracts the RETL client from the provider meta, returning a
// diagnostic error if the meta does not satisfy ClientProvider.
func service(m interface{}) (Service, diag.Diagnostics) {
	cp, ok := m.(ClientProvider)
	if !ok {
		return nil, diag.FromErr(fmt.Errorf("API client is not configured"))
	}
	return cp.RETLSourcesClient(), nil
}

// configBlock extracts the (single) `config` block from the resource state
// as a map, or returns an error if it is missing or empty.
func configBlock(d *schema.ResourceData) (map[string]interface{}, error) {
	raw, ok := d.Get("config").([]interface{})
	if !ok || len(raw) == 0 || raw[0] == nil {
		return nil, fmt.Errorf("config block is required")
	}
	cfg, ok := raw[0].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("config block has unexpected shape")
	}
	return cfg, nil
}

// stringField extracts an optional string value from a Terraform config block,
// returning "" if the key is absent or holds a nil value.
func stringField(cfg map[string]interface{}, key string) string {
	v, ok := cfg[key]
	if !ok || v == nil {
		return ""
	}
	s, _ := v.(string)
	return s
}

// sourceDefinitionName returns either the adapter-fixed value (s3_table) or
// the user-supplied schema field.
func (a typeAdapter) sourceDefinitionName(d *schema.ResourceData) string {
	if a.fixedSourceDefinitionName != "" {
		return a.fixedSourceDefinitionName
	}
	return d.Get("source_definition_name").(string)
}

func makeCreate(adapter typeAdapter) schema.CreateContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		svc, diags := service(m)
		if diags.HasError() {
			return diags
		}

		cfgBlock, err := configBlock(d)
		if err != nil {
			return diag.FromErr(err)
		}
		cfg, err := adapter.marshalConfig(cfgBlock)
		if err != nil {
			return diag.FromErr(fmt.Errorf("could not marshal config: %w", err))
		}

		req := &retl.RETLSourceCreateRequest{
			Name:                 d.Get("name").(string),
			Config:               cfg,
			SourceType:           adapter.sourceType,
			SourceDefinitionName: adapter.sourceDefinitionName(d),
			AccountID:            d.Get("account_id").(string),
			Enabled:              d.Get("enabled").(bool),
		}

		created, err := svc.CreateRetlSource(ctx, req)
		if err != nil {
			return diag.FromErr(fmt.Errorf("could not create RETL source: %w", err))
		}

		d.SetId(created.ID)
		return makeRead(adapter)(ctx, d, m)
	}
}

func makeRead(adapter typeAdapter) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		svc, diags := service(m)
		if diags.HasError() {
			return diags
		}

		source, err := svc.GetRetlSource(ctx, d.Id())
		if err != nil {
			var apiErr *client.APIError
			if errors.As(err, &apiErr) && apiErr.HTTPStatusCode == 404 {
				d.SetId("")
				return nil
			}
			return diag.FromErr(fmt.Errorf("could not read RETL source: %w", err))
		}

		return diag.FromErr(storeToState(adapter, d, source))
	}
}

func makeUpdate(adapter typeAdapter) schema.UpdateContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		svc, diags := service(m)
		if diags.HasError() {
			return diags
		}

		cfgBlock, err := configBlock(d)
		if err != nil {
			return diag.FromErr(err)
		}
		cfg, err := adapter.marshalConfig(cfgBlock)
		if err != nil {
			return diag.FromErr(fmt.Errorf("could not marshal config: %w", err))
		}

		req := &retl.RETLSourceUpdateRequest{
			Name:      d.Get("name").(string),
			Config:    cfg,
			IsEnabled: d.Get("enabled").(bool),
			AccountID: d.Get("account_id").(string),
		}

		if _, err := svc.UpdateRetlSource(ctx, d.Id(), req); err != nil {
			return diag.FromErr(fmt.Errorf("could not update RETL source: %w", err))
		}

		return makeRead(adapter)(ctx, d, m)
	}
}

func deleteSource(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	svc, diags := service(m)
	if diags.HasError() {
		return diags
	}
	if err := svc.DeleteRetlSource(ctx, d.Id()); err != nil {
		var apiErr *client.APIError
		if errors.As(err, &apiErr) && apiErr.HTTPStatusCode == 404 {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("could not delete RETL source: %w", err))
	}
	d.SetId("")
	return nil
}

func storeToState(adapter typeAdapter, d *schema.ResourceData, source *retl.RETLSource) error {
	d.SetId(source.ID)
	if err := d.Set("name", source.Name); err != nil {
		return err
	}
	if adapter.fixedSourceDefinitionName == "" {
		if err := d.Set("source_definition_name", source.SourceDefinitionName); err != nil {
			return err
		}
	}
	if err := d.Set("account_id", source.AccountID); err != nil {
		return err
	}
	if err := d.Set("enabled", source.IsEnabled); err != nil {
		return err
	}

	cfg, err := adapter.unmarshalConfig(source.Config)
	if err != nil {
		return fmt.Errorf("could not unmarshal config: %w", err)
	}
	if err := d.Set("config", cfg); err != nil {
		return err
	}

	if source.CreatedAt != nil {
		if err := d.Set("created_at", source.CreatedAt.Format(time.RFC3339)); err != nil {
			return err
		}
	}
	if source.UpdatedAt != nil {
		if err := d.Set("updated_at", source.UpdatedAt.Format(time.RFC3339)); err != nil {
			return err
		}
	}
	return nil
}
