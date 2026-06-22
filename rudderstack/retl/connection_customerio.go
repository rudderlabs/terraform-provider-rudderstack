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

// ResourceConnectionCustomerIO returns the schema for
// `rudderstack_retl_connection_customerio` — a RETL connection scoped to a
// Customer.io destination using the VDM v2 flow.
//
// The destination `object` is a first-class typed field. Internally it
// round-trips through the API's untyped `destinationConfig` JSON as
// `{"object": "..."}`. identifiers and mappings flow through the base schema
// (VDM v2 identifierMappings / fieldMappings); config-be assembles the VDM v2
// connectionConfig server-side from the Customer.io destination definition.
//
// This follows the same typed-destination pattern as
// rudderstack_retl_connection_customerio_audience: baseConnectionSchema()
// composed with the destination's required fields plus a small CRUD shim to
// pack/unpack destinationConfig.
func ResourceConnectionCustomerIO() *schema.Resource {
	return &schema.Resource{
		Description: "A RETL connection to a Customer.io destination (VDM v2). " +
			"Carries the destination object as a typed top-level field; ForceNew because the " +
			"destinationConfig shape is not mutable in place on this flow.",
		Schema: mergeSchemas(baseConnectionSchema(), map[string]*schema.Schema{
			"object": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				Description:  "Customer.io destination object (e.g. `customers`). Packed into destinationConfig.",
			},
			// VDM v2 supports only upsert and mirror — drop `full` from the base
			// schema's allowed set so users see a plan-time error instead of an
			// API rejection on apply.
			"sync_behaviour": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"upsert", "mirror"}, false),
				Description:  "How records are synced to the destination: `upsert` or `mirror` (VDM v2).",
			},
			// cursor_column is a generic source-side field (the incremental
			// watermark column), sent as a top-level request field — NOT inside
			// destinationConfig. Only valid when sync_behaviour is `upsert`
			// (enforced in customizeCustomerIOConnectionDiff).
			"cursor_column": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Column name for incremental upsert syncs (only valid when sync_behaviour is `upsert`).",
			},
		}),
		CreateContext: createCustomerIOConnection,
		ReadContext:   readCustomerIOConnection,
		UpdateContext: updateCustomerIOConnection,
		DeleteContext: deleteConnection,
		CustomizeDiff: customizeCustomerIOConnectionDiff,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// customizeCustomerIOConnectionDiff rejects cursor_column when sync_behaviour
// is not `upsert`, surfacing the error at plan time instead of an API
// rejection on apply (mirroring the generic resource's customizeConnectionDiff).
func customizeCustomerIOConnectionDiff(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
	if cursor := d.Get("cursor_column").(string); cursor != "" {
		if sb := d.Get("sync_behaviour").(string); sb != "" && sb != "upsert" {
			return fmt.Errorf("cursor_column is only valid when sync_behaviour is %q, got %q", "upsert", sb)
		}
	}
	return nil
}

func createCustomerIOConnection(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	svc, diags := service(m)
	if diags.HasError() {
		return diags
	}

	req := &retl.CreateRETLConnectionRequest{}
	if err := applyBaseToCreateRequest(d, req); err != nil {
		return diag.FromErr(err)
	}
	cfg, err := encodeCustomerIOObjectConfig(d)
	if err != nil {
		return diag.FromErr(err)
	}
	req.DestinationConfig = cfg
	// cursor_column is a generic top-level source field (not destinationConfig).
	if v := d.Get("cursor_column").(string); v != "" {
		req.CursorColumn = v
	}

	created, err := svc.CreateConnection(ctx, req)
	if err != nil {
		return diag.FromErr(fmt.Errorf("could not create RETL connection: %w", err))
	}
	d.SetId(created.ID)
	return readCustomerIOConnection(ctx, d, m)
}

func readCustomerIOConnection(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	if err := storeBaseConnectionToState(d, conn); err != nil {
		return diag.FromErr(err)
	}
	// cursor_column is a generic source-side field returned at top level
	// (independent of destinationConfig validity).
	if err := d.Set("cursor_column", conn.CursorColumn); err != nil {
		return diag.FromErr(fmt.Errorf("set cursor_column: %w", err))
	}
	// Empty or JSON-`null` destinationConfig on a Customer.io connection is a
	// server-side inconsistency — the connection MUST carry an object. Don't
	// zero the field in state: `object` is ForceNew with StringIsNotEmpty, so a
	// zero value would produce a plan that can never reconcile. Leave the prior
	// state value intact and surface a warning. (Same rationale as the audience
	// resource — see warnMissingCustomerIOAudienceConfig.)
	if len(conn.DestinationConfig) == 0 {
		return diag.Diagnostics{warnMissingCustomerIOObjectConfig(conn.ID, "destinationConfig is empty")}
	}
	object, err := decodeCustomerIOObject(conn.DestinationConfig)
	if err != nil {
		if errors.Is(err, errCustomerIONullConfig) {
			return diag.Diagnostics{warnMissingCustomerIOObjectConfig(conn.ID, "destinationConfig is JSON null")}
		}
		return diag.FromErr(err)
	}
	if err := d.Set("object", object); err != nil {
		return diag.FromErr(fmt.Errorf("set object: %w", err))
	}
	return nil
}

// warnMissingCustomerIOObjectConfig formats a warning diagnostic for the
// inconsistent-server-state case where a Customer.io connection comes back
// without a usable destinationConfig. Returns a Warning (not an Error) so
// refresh still succeeds — silently zeroing `object` would produce a
// never-reconcilable plan (ForceNew + StringIsNotEmpty). See the audience
// resource's warnMissingCustomerIOAudienceConfig for the full rationale.
func warnMissingCustomerIOObjectConfig(connID, reason string) diag.Diagnostic {
	return diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  fmt.Sprintf("Customer.io connection %q is missing object", connID),
		Detail: fmt.Sprintf(
			"The server returned a connection with no object (%s). "+
				"This is an inconsistent server state — object is mandatory for Customer.io VDM v2 destinations. "+
				"Terraform preserved the prior object value in state and will NOT automatically reconcile this. "+
				"To recover: fix the connection in the RudderStack UI, or force a replacement with "+
				"`terraform apply -replace=<resource address>`.",
			reason,
		),
	}
}

func updateCustomerIOConnection(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	svc, diags := service(m)
	if diags.HasError() {
		return diags
	}

	// `object` is ForceNew, so an Update never sees a change to it —
	// applyBaseToUpdateRequest covers the mutable fields.
	req := &retl.UpdateRETLConnectionRequest{}
	if err := applyBaseToUpdateRequest(d, req); err != nil {
		return diag.FromErr(err)
	}
	if _, err := svc.UpdateConnection(ctx, d.Id(), req); err != nil {
		return diag.FromErr(fmt.Errorf("could not update RETL connection: %w", err))
	}
	return readCustomerIOConnection(ctx, d, m)
}

// encodeCustomerIOObjectConfig packs the top-level `object` field into the
// destinationConfig JSON shape the API expects.
func encodeCustomerIOObjectConfig(d *schema.ResourceData) (json.RawMessage, error) {
	object, ok := d.Get("object").(string)
	if !ok {
		// Defensive: schema declares object as TypeString so the SDK should
		// always hand us a string.
		return nil, fmt.Errorf("object has unexpected type %T", d.Get("object"))
	}
	out, err := json.Marshal(map[string]any{"object": object})
	if err != nil {
		return nil, fmt.Errorf("encode customerio destinationConfig: %w", err)
	}
	return out, nil
}

// errCustomerIONullConfig signals that destinationConfig was the JSON literal
// `null` — semantically "no destination-specific config", not a malformed
// payload. Callers distinguish this from a hard decode error.
var errCustomerIONullConfig = errors.New("customerio destinationConfig is null")

// decodeCustomerIOObject extracts the destination object from a
// destinationConfig JSON blob. Returns errCustomerIONullConfig (a soft signal)
// for JSON `null`. Returns a hard error when the payload is a non-null shape
// without a non-empty string `object` — that signals an unsupported
// destination-specific connection (e.g. imported from a destination that isn't
// Customer.io VDM v2), which this resource does not represent.
func decodeCustomerIOObject(raw json.RawMessage) (string, error) {
	var parsed map[string]interface{}
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return "", fmt.Errorf("decode customerio destinationConfig: %w", err)
	}
	if parsed == nil {
		// JSON `null` unmarshalls into a nil map — treat as "no typed config".
		return "", errCustomerIONullConfig
	}
	v, ok := parsed["object"]
	if !ok {
		return "", fmt.Errorf("destinationConfig has no object — only Customer.io VDM v2 destination-specific connections are supported")
	}
	s, ok := v.(string)
	if !ok {
		return "", fmt.Errorf("destinationConfig object is %T, expected string", v)
	}
	if s == "" {
		return "", fmt.Errorf("destinationConfig object is empty")
	}
	return s, nil
}
