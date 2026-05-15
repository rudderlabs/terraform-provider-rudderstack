package retl

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/rudderlabs/rudder-iac/api/client"
	"github.com/rudderlabs/rudder-iac/api/client/retl"
)

// ResourceConnectionCustomerIOAudience returns the schema for
// `rudderstack_retl_connection_customerio_audience` — a RETL connection
// scoped to a Customer.io Audience destination.
//
// The audience ID is a first-class typed field. Internally it round-trips
// through the API's untyped `destinationConfig` JSON as `{"audienceId": N}`
// so the server-side flow detection still works, but Terraform users see and
// validate an integer field rather than an opaque blob.
//
// Customer.io Audience is the first destination-specific RETL flow exposed
// this way. Adding another typed destination follows the same pattern: a
// new ResourceConnection<Destination> function composes
// baseConnectionSchema() with the destination's required fields and adds a
// small CRUD shim to pack/unpack destinationConfig.
func ResourceConnectionCustomerIOAudience() *schema.Resource {
	return &schema.Resource{
		Description: "A RETL connection to a Customer.io Audience destination. " +
			"Carries the audience ID as a typed top-level field; ForceNew because the " +
			"Customer.io Audience API does not accept destinationConfig changes on update.",
		Schema: mergeSchemas(baseConnectionSchema(), map[string]*schema.Schema{
			"audience_id": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntAtLeast(1),
				Description:  "Customer.io audience ID (positive integer).",
			},
		}),
		CreateContext: createCustomerIOAudienceConnection,
		ReadContext:   readCustomerIOAudienceConnection,
		UpdateContext: updateCustomerIOAudienceConnection,
		DeleteContext: deleteConnection,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func createCustomerIOAudienceConnection(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	svc, diags := service(m)
	if diags.HasError() {
		return diags
	}

	req := &retl.CreateRETLConnectionRequest{}
	if err := applyBaseToCreateRequest(d, req); err != nil {
		return diag.FromErr(err)
	}
	cfg, err := encodeCustomerIOAudienceConfig(d)
	if err != nil {
		return diag.FromErr(err)
	}
	req.DestinationConfig = cfg

	created, err := svc.CreateConnection(ctx, req)
	if err != nil {
		return diag.FromErr(fmt.Errorf("could not create RETL connection: %w", err))
	}
	d.SetId(created.ID)
	return readCustomerIOAudienceConnection(ctx, d, m)
}

func readCustomerIOAudienceConnection(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
	if len(conn.DestinationConfig) == 0 {
		// No destinationConfig on a typed Customer.io Audience resource is a
		// server-side inconsistency — clear the field so the next plan shows
		// the gap rather than silently keeping a stale audience_id.
		if err := d.Set("audience_id", 0); err != nil {
			return diag.FromErr(fmt.Errorf("set audience_id: %w", err))
		}
		return nil
	}
	id, err := decodeCustomerIOAudienceID(conn.DestinationConfig)
	if err != nil {
		if errors.Is(err, errCustomerIOAudienceNullConfig) {
			// JSON `null` carries the same "no typed config" meaning as an empty
			// payload — clear the field for the same reason.
			if err := d.Set("audience_id", 0); err != nil {
				return diag.FromErr(fmt.Errorf("set audience_id: %w", err))
			}
			return nil
		}
		return diag.FromErr(err)
	}
	if err := d.Set("audience_id", id); err != nil {
		return diag.FromErr(fmt.Errorf("set audience_id: %w", err))
	}
	return nil
}

func updateCustomerIOAudienceConnection(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	svc, diags := service(m)
	if diags.HasError() {
		return diags
	}

	// audience_id is ForceNew, so an Update never sees a change to it —
	// applyBaseToUpdateRequest covers the rest.
	req := &retl.UpdateRETLConnectionRequest{}
	if err := applyBaseToUpdateRequest(d, req); err != nil {
		return diag.FromErr(err)
	}
	if _, err := svc.UpdateConnection(ctx, d.Id(), req); err != nil {
		return diag.FromErr(fmt.Errorf("could not update RETL connection: %w", err))
	}
	return readCustomerIOAudienceConnection(ctx, d, m)
}

// encodeCustomerIOAudienceConfig packs the top-level `audience_id` field
// into the destinationConfig JSON shape the API expects.
func encodeCustomerIOAudienceConfig(d *schema.ResourceData) (json.RawMessage, error) {
	id, ok := d.Get("audience_id").(int)
	if !ok {
		// Defensive: schema declares audience_id as TypeInt so the SDK should
		// always hand us an int. Surface the divergence loudly rather than
		// silently POST `{"audienceId": 0}`.
		return nil, fmt.Errorf("audience_id has unexpected type %T", d.Get("audience_id"))
	}
	out, err := json.Marshal(map[string]any{"audienceId": id})
	if err != nil {
		return nil, fmt.Errorf("encode customerio_audience destinationConfig: %w", err)
	}
	return out, nil
}

// errCustomerIOAudienceNullConfig signals that destinationConfig was the JSON
// literal `null` — semantically "no destination-specific config", not a
// malformed payload. Callers distinguish this from a hard decode error.
var errCustomerIOAudienceNullConfig = errors.New("customerio_audience destinationConfig is null")

// decodeCustomerIOAudienceID extracts the integer audienceId from a
// destinationConfig JSON blob. Returns errCustomerIOAudienceNullConfig (a
// soft signal) for JSON `null`. Returns a hard error when the payload is a
// non-null shape without a numeric `audienceId` — that signals an
// unsupported destination-specific connection (e.g. imported from a
// destination that isn't Customer.io Audience), which this resource does
// not represent.
func decodeCustomerIOAudienceID(raw json.RawMessage) (int, error) {
	var parsed map[string]interface{}
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return 0, fmt.Errorf("decode customerio_audience destinationConfig: %w", err)
	}
	if parsed == nil {
		// JSON `null` unmarshalls into a nil map — treat as "no typed config".
		return 0, errCustomerIOAudienceNullConfig
	}
	v, ok := parsed["audienceId"]
	if !ok {
		return 0, fmt.Errorf("destinationConfig has no audienceId — only Customer.io Audience destination-specific connections are supported")
	}
	n, ok := v.(float64) // json.Unmarshal decodes JSON numbers as float64
	if !ok {
		return 0, fmt.Errorf("destinationConfig audienceId is %T, expected number", v)
	}
	if math.IsNaN(n) || math.IsInf(n, 0) || math.Trunc(n) != n {
		return 0, fmt.Errorf("destinationConfig audienceId %v is not an integer", n)
	}
	// int(n) is exact for |n| < 2^53 (float64 mantissa). Customer.io audience
	// IDs are well below that bound in practice; the integrality check above
	// rejects fractional values like 42.5 that would otherwise truncate silently.
	return int(n), nil
}
