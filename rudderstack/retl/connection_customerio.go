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
// `{"object": "..."}`. identifiers flow through the base schema; VDM v2 does
// not support field mappings, so this resource has no `mappings`. config-be
// assembles the VDM v2 connectionConfig server-side from the Customer.io
// destination definition.
//
// This follows the same typed-destination pattern as
// rudderstack_retl_connection_customerio_audience: baseConnectionSchema()
// composed with the destination's required fields plus a small CRUD shim to
// pack/unpack destinationConfig.
func ResourceConnectionCustomerIO() *schema.Resource {
	return &schema.Resource{
		Description: "A RETL connection to a Customer.io destination. " +
			"Carries the destination object as a typed top-level field; ForceNew because the " +
			"object cannot be changed in place — changing it recreates the connection.",
		Schema: mergeSchemas(baseConnectionSchema(), map[string]*schema.Schema{
			// Customer.io supports exactly one object, whose on-the-wire value is
			// `person` (the `value` from the listObjects API). Restrict to it so
			// typos fail at plan time instead of on apply. If Customer.io ever
			// adds objects, extend this slice.
			"object": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"person"}, false),
				Description:  "Customer.io destination object. Only `person` is supported.",
			},
			// Only upsert and mirror are supported — drop `full` from the base
			// schema's allowed set so users see a plan-time error instead of an
			// API rejection on apply.
			"sync_behaviour": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"upsert", "mirror"}, false),
				Description:  "How records are synced to the destination: `upsert` or `mirror`.",
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
// rejection on apply.
func customizeCustomerIOConnectionDiff(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
	return validateCursorColumnUpsertOnly(d)
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
	// A missing/empty/malformed object is a hard error (see decodeCustomerIOObject)
	// rather than a warning — surfacing it at refresh beats masking a broken
	// connection behind a plan that stays a silent no-op.
	object, err := decodeCustomerIOObject(conn.DestinationConfig)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("object", object); err != nil {
		return diag.FromErr(fmt.Errorf("set object: %w", err))
	}
	return nil
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

// decodeCustomerIOObject extracts the destination object from the connection's
// destination config. A successful GetConnection (200) returns a complete body,
// so any shape without a non-empty string `object` — empty input, JSON `null`,
// missing key, non-string, or empty string — is a persistent server-side
// inconsistency (this resource always writes a non-empty object on create) or
// signals a connection that isn't a Customer.io connection at all. All of these
// are hard errors so the problem surfaces at refresh instead of being masked by
// a warning that turns the plan into a silent no-op.
func decodeCustomerIOObject(raw json.RawMessage) (string, error) {
	var parsed map[string]interface{}
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return "", fmt.Errorf("decode customerio destination config: %w", err)
	}
	v, ok := parsed["object"]
	if !ok {
		// nil map (JSON `null`) and an object-less payload land here together.
		return "", fmt.Errorf("connection has no object — only Customer.io connections are supported by this resource")
	}
	s, ok := v.(string)
	if !ok {
		return "", fmt.Errorf("connection object is %T, expected string", v)
	}
	if s == "" {
		return "", fmt.Errorf("connection object is empty")
	}
	return s, nil
}
