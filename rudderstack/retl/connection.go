package retl

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/rudderlabs/rudder-iac/api/client"
	"github.com/rudderlabs/rudder-iac/api/client/retl"
)

// ResourceConnection returns the schema for `rudderstack_retl_connection` —
// the generic resource covering JSON Mapper and Object Mapping flows.
//
// Flow detection (JSON Mapper vs Object Mapping) happens server-side from the
// destination definition. This resource only sends the flat schema and lets
// the API assemble the internal config. The `object` attribute distinguishes
// Object Mapping (set) from JSON Mapper (absent).
//
// Destination-specific flows (Customer.io Audience etc.) are intentionally
// out of scope here. Each destination-specific flow gets its own typed
// resource (e.g. `rudderstack_retl_connection_customerio_audience`) that
// exposes the destination's required fields as first-class schema attributes
// instead of stuffing them into a generic destination_config blob. The
// generic resource refuses to refresh a connection whose destinationConfig
// has been populated by a destination-specific flow so users get a clear
// signal at refresh time to switch to the typed resource.
func ResourceConnection() *schema.Resource {
	return &schema.Resource{
		Description: "A generic RETL connection between a RETL source and a destination, " +
			"covering JSON Mapper and Object Mapping flows. Destination-specific flows " +
			"(e.g. Customer.io Audience) have their own typed resources.",
		Schema:        genericConnectionSchema(),
		CreateContext: createConnection,
		ReadContext:   readConnection,
		UpdateContext: updateConnection,
		DeleteContext: deleteConnection,
		CustomizeDiff: customizeConnectionDiff,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// genericConnectionSchema composes the base RETL connection schema with the
// generic-flow-only fields: `event` (JSON Mapper), `constants`, `cursor_column`,
// `object` (Object Mapping). Identifiers ForceNew lives in the base schema
// because every flow treats identifier changes as breaking.
// `constants` ForceNew is conditional on Object Mapping and is applied in
// CustomizeDiff (Object Mapping = ForceNew; JSON Mapper = mutable), so the
// schema declares it without ForceNew.
func genericConnectionSchema() map[string]*schema.Schema {
	return mergeSchemas(baseConnectionSchema(), map[string]*schema.Schema{
		"event": {
			Type:     schema.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"type": {
						Type:         schema.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringInSlice([]string{"identify", "track"}, false),
					},
					"name":        {Type: schema.TypeString, Optional: true, ForceNew: true},
					"name_column": {Type: schema.TypeString, Optional: true, ForceNew: true},
				},
			},
			Description: "CDP event configuration. Optional in the Terraform schema; flow-specific " +
				"requirements (required for JSON Mapper, absent for Object Mapping) are enforced by the API.",
		},
		"constants": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"key":   {Type: schema.TypeString, Required: true},
					"value": {Type: schema.TypeString, Required: true},
				},
			},
			Description: "User-defined constants. Mutable for JSON Mapper; ForceNew for Object Mapping.",
		},
		"cursor_column": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "Column name for incremental upsert syncs (only valid when sync_behaviour is `upsert`).",
		},
		"object": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "Destination entity for Object Mapping flows (e.g. `Contact`, `Lead`).",
		},
	})
}

// customizeConnectionDiff applies the Object-Mapping-only `constants` ForceNew
// and rejects locally-detectable invalid combinations (cursor_column only with
// sync_behaviour=upsert) so users see the error at plan time instead of on
// apply.
func customizeConnectionDiff(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
	if cursor := d.Get("cursor_column").(string); cursor != "" {
		if sb := d.Get("sync_behaviour").(string); sb != "" && sb != "upsert" {
			return fmt.Errorf("cursor_column is only valid when sync_behaviour is %q, got %q", "upsert", sb)
		}
	}

	if d.Id() == "" {
		// Brand-new resource — ForceNew is irrelevant on create.
		return nil
	}

	// Object Mapping (object set) treats `constants` as immutable; JSON Mapper
	// allows in-place updates.
	if d.Get("object").(string) != "" && d.HasChange("constants") {
		if err := d.ForceNew("constants"); err != nil {
			return err
		}
	}
	return nil
}

func createConnection(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	svc, diags := service(m)
	if diags.HasError() {
		return diags
	}

	req, err := buildCreateRequest(d)
	if err != nil {
		return diag.FromErr(err)
	}

	created, err := svc.CreateConnection(ctx, req)
	if err != nil {
		return diag.FromErr(fmt.Errorf("could not create RETL connection: %w", err))
	}

	d.SetId(created.ID)
	return readConnection(ctx, d, m)
}

func readConnection(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	svc, diags := service(m)
	if diags.HasError() {
		return diags
	}

	conn, err := svc.GetConnection(ctx, d.Id())
	if err != nil {
		var apiErr *client.APIError
		if errors.As(err, &apiErr) && apiErr.HTTPStatusCode == 404 {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("could not read RETL connection: %w", err))
	}

	return diag.FromErr(storeGenericConnectionToState(d, conn))
}

func updateConnection(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	svc, diags := service(m)
	if diags.HasError() {
		return diags
	}

	req, err := buildUpdateRequest(d)
	if err != nil {
		return diag.FromErr(err)
	}
	if _, err := svc.UpdateConnection(ctx, d.Id(), req); err != nil {
		return diag.FromErr(fmt.Errorf("could not update RETL connection: %w", err))
	}

	return readConnection(ctx, d, m)
}

// deleteConnection deletes any RETL connection (by ID). Shared between the
// generic and per-destination resources — delete has no flow-specific logic.
func deleteConnection(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	svc, diags := service(m)
	if diags.HasError() {
		return diags
	}

	if err := svc.DeleteConnection(ctx, d.Id()); err != nil {
		var apiErr *client.APIError
		if errors.As(err, &apiErr) && apiErr.HTTPStatusCode == 404 {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("could not delete RETL connection: %w", err))
	}
	d.SetId("")
	return nil
}

func buildCreateRequest(d *schema.ResourceData) (*retl.CreateRETLConnectionRequest, error) {
	req := &retl.CreateRETLConnectionRequest{}
	if err := applyBaseToCreateRequest(d, req); err != nil {
		return nil, err
	}
	if event, ok := eventFromState(d); ok {
		req.Event = event
	}
	if constants := constantsFromState(d); len(constants) > 0 {
		req.Constants = constants
	}
	if v := d.Get("cursor_column").(string); v != "" {
		req.CursorColumn = v
	}
	if v := d.Get("object").(string); v != "" {
		req.Object = v
	}
	return req, nil
}

func buildUpdateRequest(d *schema.ResourceData) (*retl.UpdateRETLConnectionRequest, error) {
	req := &retl.UpdateRETLConnectionRequest{}
	if err := applyBaseToUpdateRequest(d, req); err != nil {
		return nil, err
	}
	if d.HasChange("constants") {
		constants := constantsFromState(d)
		req.Constants = &constants
	}
	return req, nil
}

func storeGenericConnectionToState(d *schema.ResourceData, c *retl.RETLConnection) error {
	if err := storeBaseConnectionToState(d, c); err != nil {
		return err
	}
	if err := d.Set("event", eventToState(c.Event)); err != nil {
		return fmt.Errorf("set event: %w", err)
	}
	if err := d.Set("constants", constantsToState(c.Constants)); err != nil {
		return fmt.Errorf("set constants: %w", err)
	}
	if err := d.Set("cursor_column", c.CursorColumn); err != nil {
		return fmt.Errorf("set cursor_column: %w", err)
	}
	if err := d.Set("object", c.Object); err != nil {
		return fmt.Errorf("set object: %w", err)
	}

	// The API returns destinationConfig as a non-empty JSON payload only for
	// destination-specific flows. The generic resource does not represent
	// those — fail loudly at refresh so the user knows to switch resources
	// rather than silently dropping config from state. JSON `null` is the
	// server's way of saying "no destination-specific config" — treat it as a
	// no-op.
	//
	// The payload itself is not echoed into the error: destinationConfig can
	// carry credentials for destination-specific flows we haven't yet typed
	// (e.g. API keys), and Terraform diagnostics surface in CI logs.
	if len(c.DestinationConfig) > 0 {
		var parsed any
		if err := json.Unmarshal(c.DestinationConfig, &parsed); err != nil {
			return fmt.Errorf("decode destinationConfig: %w", err)
		}
		if parsed != nil {
			return fmt.Errorf(
				"connection %q has destination-specific configuration (destinationConfig is %d bytes); "+
					"the generic rudderstack_retl_connection resource does not represent destination-specific flows. "+
					"Use a typed resource such as rudderstack_retl_connection_customerio_audience instead.",
				c.ID, len(c.DestinationConfig),
			)
		}
	}
	return nil
}

// --- per-block helpers specific to the generic resource ---

func eventFromState(d *schema.ResourceData) (*retl.Event, bool) {
	raw, ok := d.Get("event").([]interface{})
	if !ok || len(raw) == 0 || raw[0] == nil {
		return nil, false
	}
	m := raw[0].(map[string]interface{})
	e := &retl.Event{Type: retl.EventType(m["type"].(string))}
	if v, _ := m["name"].(string); v != "" {
		e.Name = v
	}
	if v, _ := m["name_column"].(string); v != "" {
		e.NameColumn = v
	}
	return e, true
}

func eventToState(e *retl.Event) []map[string]interface{} {
	if e == nil {
		return nil
	}
	m := map[string]interface{}{"type": string(e.Type)}
	if e.Name != "" {
		m["name"] = e.Name
	}
	if e.NameColumn != "" {
		m["name_column"] = e.NameColumn
	}
	return []map[string]interface{}{m}
}

func constantsFromState(d *schema.ResourceData) []retl.Constant {
	raw, _ := d.Get("constants").([]interface{})
	out := make([]retl.Constant, 0, len(raw))
	for _, item := range raw {
		m := item.(map[string]interface{})
		out = append(out, retl.Constant{
			Key:   m["key"].(string),
			Value: m["value"].(string),
		})
	}
	return out
}

func constantsToState(cs []retl.Constant) []map[string]interface{} {
	if len(cs) == 0 {
		return nil
	}
	out := make([]map[string]interface{}, 0, len(cs))
	for _, c := range cs {
		out = append(out, map[string]interface{}{"key": c.Key, "value": c.Value})
	}
	return out
}
