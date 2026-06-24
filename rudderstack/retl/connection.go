package retl

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/rudderlabs/rudder-iac/api/client"
	"github.com/rudderlabs/rudder-iac/api/client/retl"
)

// resourceGetter is the subset of *schema.ResourceData / *schema.ResourceDiff
// used by helpers that need to read state without caring which one they got.
type resourceGetter interface {
	Get(string) interface{}
}

// validJSONMapperIdentifierTargets enumerates the only `identifiers[*].to`
// values the server accepts for JSON Mapper. Kept in sync with the server-side
// IDENTIFIER_TARGETS constant in
// rudder-config-backend/.../api-gateway/connection-config/constants.ts.
var validJSONMapperIdentifierTargets = []string{"user_id", "anonymous_id"}

// validateJSONMapperIdentifierTargets enforces the JSON Mapper rule at plan
// time. Caller must already have determined this is the JSON Mapper flow
// (`object` unset). Accepts resourceGetter so the helper is usable from both
// schema.ResourceDiff and schema.ResourceData, and trivially testable with
// TestResourceData.
func validateJSONMapperIdentifierTargets(d resourceGetter) error {
	raw, ok := d.Get("identifiers").([]interface{})
	if !ok {
		return nil
	}
	for i, item := range raw {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		to, _ := m["to"].(string)
		if !slices.Contains(validJSONMapperIdentifierTargets, to) {
			return fmt.Errorf(
				`identifiers[%d].to must be one of [%s] for JSON Mapper flow (got %q); set "object" to use Object Mapping or change identifiers[%d].to to a valid value`,
				i, strings.Join(validJSONMapperIdentifierTargets, ", "), to, i,
			)
		}
	}
	return nil
}

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
// generic-flow-only fields: `event` (JSON Mapper), `constants`, `object`
// (Object Mapping), and `mappings` (field mappings). `cursor_column` and
// `identifiers` live in the base schema because every flow shares them
// (identifiers are mutable; changes are forwarded on update).
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
			// terraform-plugin-sdk v2: ForceNew on a TypeList parent only fires when the
			// list's element count changes — nested-field edits (e.g. event.0.type
			// flipping from "identify" to "track") plan as in-place updates and are
			// silently dropped because `event` is excluded from buildUpdateRequest.
			// Mirror ForceNew onto every inner field to force destroy+create on any
			// change.
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
		"object": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "Destination entity for Object Mapping flows (e.g. `Contact`, `Lead`).",
		},
		// mappings (field mappings) is generic-flow-only. Destination-specific
		// flows (e.g. customerio VDM v2) do not expose it.
		"mappings": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"from": {Type: schema.TypeString, Required: true},
					"to":   {Type: schema.TypeString, Required: true},
				},
			},
			Description: "Source-to-destination field mappings (mutable).",
		},
	})
}

// customizeConnectionDiff applies the Object-Mapping-only `constants` ForceNew
// and rejects locally-detectable invalid combinations (cursor_column only with
// sync_behaviour=upsert) so users see the error at plan time instead of on
// apply.
func customizeConnectionDiff(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
	if err := validateCursorColumnUpsertOnly(d); err != nil {
		return err
	}

	// JSON Mapper restricts identifiers[*].to to {user_id, anonymous_id}. The
	// server rejects other values with a confusing generic error at apply time;
	// surface the same rule at plan time so users can fix the .tf without an
	// API round-trip. Runs on both create and update.
	//
	// In the typed-resource design, the generic resource only covers JSON Mapper
	// and Object Mapping — destination-specific flows live in their own typed
	// resources (e.g. rudderstack_retl_connection_customerio_audience), so the
	// JSON Mapper check fires whenever `object` is unset.
	if d.Get("object").(string) == "" {
		if err := validateJSONMapperIdentifierTargets(d); err != nil {
			return err
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
	if v := d.Get("object").(string); v != "" {
		req.Object = v
	}
	if mappings := mappingsFromState(d, "mappings"); len(mappings) > 0 {
		req.Mappings = mappings
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
	if d.HasChange("mappings") {
		mappings := mappingsFromState(d, "mappings")
		req.Mappings = &mappings
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
	if err := d.Set("object", c.Object); err != nil {
		return fmt.Errorf("set object: %w", err)
	}
	if err := d.Set("mappings", mappingsToState(c.Mappings)); err != nil {
		return fmt.Errorf("set mappings: %w", err)
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
					"Use the typed rudderstack_retl_connection_<destination> resource that matches this connection's destination "+
					"(e.g. rudderstack_retl_connection_customerio_audience for Customer.io Audience).",
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
