package retl

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/rudderlabs/rudder-iac/api/client"
	"github.com/rudderlabs/rudder-iac/api/client/retl"
)

// ResourceConnection returns the schema for `rudderstack_retl_connection`.
//
// Flow detection (JSON Mapper / Object mapping / Destination-specific) happens
// server-side based on the destination definition; this resource only sends
// the flat schema and lets the API assemble the internal config. ForceNew on
// `identifiers` and `constants` is conditional on the detected flow and is
// applied via CustomizeDiff using the same field-presence signals the API uses.
func ResourceConnection() *schema.Resource {
	return &schema.Resource{
		Description: "A RETL connection between a RETL source and a destination. " +
			"Flow type (JSON Mapper, Object mapping, Destination-specific) is " +
			"determined by the destination definition.",
		Schema:        connectionSchema(),
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

func connectionSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"source_id": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "ID of the RETL source.",
		},
		"destination_id": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "ID of the destination.",
		},
		"enabled": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "Whether the connection is enabled.",
		},
		"external_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Optional external identifier for CLI/IaC state tracking.",
		},
		"schedule": {
			Type:     schema.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"type": {
						Type:         schema.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice([]string{"basic", "manual", "cron"}, false),
						Description:  "Schedule type: `basic`, `manual`, or `cron`.",
					},
					"every_minutes": {
						Type:         schema.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntAtLeast(5),
						Description:  "Sync interval in minutes. Required when `type` is `basic`.",
					},
				},
			},
		},
		"sync_settings": {
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"sync_logs_config": {
						Type:     schema.TypeList,
						Optional: true,
						Computed: true,
						MaxItems: 1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"enabled": {
									Type:     schema.TypeBool,
									Optional: true,
									Default:  true,
								},
								"log_retention_in_days": {
									Type:     schema.TypeInt,
									Optional: true,
									Default:  30,
								},
								"snapshots_to_retain": {
									Type:     schema.TypeInt,
									Optional: true,
									Default:  5,
								},
							},
						},
					},
					"failed_keys_config": {
						Type:     schema.TypeList,
						Optional: true,
						Computed: true,
						MaxItems: 1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"enable_failed_keys_retry": {
									Type:     schema.TypeBool,
									Optional: true,
									Default:  true,
								},
							},
						},
					},
				},
			},
		},
		"sync_behaviour": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice([]string{"upsert", "mirror", "full"}, false),
			Description:  "How records are synced to the destination: `upsert`, `mirror`, or `full`.",
		},
		"identifiers": {
			Type:     schema.TypeList,
			Required: true,
			MinItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"from": {
						Type:     schema.TypeString,
						Required: true,
					},
					"to": {
						Type:     schema.TypeString,
						Required: true,
					},
				},
			},
			Description: "Source-to-destination identifier mappings. ForceNew for JSON Mapper " +
				"and Object Mapping flows; mutable for destination-specific flows.",
		},
		"mappings": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"from": {
						Type:     schema.TypeString,
						Required: true,
					},
					"to": {
						Type:     schema.TypeString,
						Required: true,
					},
				},
			},
			Description: "Source-to-destination field mappings (mutable for all flows).",
		},
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
						ValidateFunc: validation.StringInSlice([]string{"identify", "track"}, false),
					},
					"name": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"name_column": {
						Type:     schema.TypeString,
						Optional: true,
					},
				},
			},
			Description: "CDP event configuration. Optional in the Terraform schema; flow-specific " +
				"requirements (required for JSON Mapper, absent for other flows) are enforced by the API.",
		},
		"constants": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"key": {
						Type:     schema.TypeString,
						Required: true,
					},
					"value": {
						Type:     schema.TypeString,
						Required: true,
					},
				},
			},
			Description: "User-defined constants. Mutable for JSON Mapper; ForceNew for Object Mapping " +
				"and destination-specific flows.",
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
		"destination_config": {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsJSON,
			Description:  "Destination-specific configuration as a JSON-encoded string.",
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
}

// customizeConnectionDiff applies flow-dependent ForceNew on `identifiers` and
// `constants`, and rejects locally-detectable invalid combinations (cursor_column
// only with sync_behaviour=upsert) so users see the error at plan time instead
// of on apply.
//
// Flow is inferred from the same signals the server uses: `object` set =>
// Object Mapping, `destination_config` set => Destination-specific, otherwise
// => JSON Mapper.
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

	objectSet := d.Get("object").(string) != ""
	destConfigSet := d.Get("destination_config").(string) != ""

	identifiersForceNew, constantsForceNew := flowForceNewRules(objectSet, destConfigSet)

	if identifiersForceNew && d.HasChange("identifiers") {
		if err := d.ForceNew("identifiers"); err != nil {
			return err
		}
	}
	if constantsForceNew && d.HasChange("constants") {
		if err := d.ForceNew("constants"); err != nil {
			return err
		}
	}
	return nil
}

// flowForceNewRules returns (identifiersForceNew, constantsForceNew) for the
// detected flow.
func flowForceNewRules(objectSet, destConfigSet bool) (identifiers, constants bool) {
	switch {
	case objectSet:
		return true, true // Object Mapping
	case destConfigSet:
		return false, true // Destination-specific
	default:
		return true, false // JSON Mapper
	}
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

	return diag.FromErr(storeConnectionToState(d, conn))
}

func updateConnection(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	svc, diags := service(m)
	if diags.HasError() {
		return diags
	}

	// Apply the main payload first so a failure there does not leave a partial
	// external_id change applied server-side. external_id is touched only after
	// the main update succeeds.
	if hasConnectionUpdatePayload(d) {
		req := buildUpdateRequest(d)
		if _, err := svc.UpdateConnection(ctx, d.Id(), req); err != nil {
			return diag.FromErr(fmt.Errorf("could not update RETL connection: %w", err))
		}
	}

	if d.HasChange("external_id") {
		extID := d.Get("external_id").(string)
		if err := svc.SetConnectionExternalId(ctx, &retl.SetRETLConnectionExternalIDRequest{
			ID:         d.Id(),
			ExternalID: extID,
		}); err != nil {
			return diag.FromErr(fmt.Errorf("could not set RETL connection external id: %w", err))
		}
	}

	return readConnection(ctx, d, m)
}

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

// hasConnectionUpdatePayload returns true when at least one field in the main
// PUT payload changed (i.e. an actual update endpoint call is needed). When
// only external_id changed, the dedicated set-external-id endpoint suffices.
func hasConnectionUpdatePayload(d *schema.ResourceData) bool {
	for _, k := range []string{"enabled", "schedule", "sync_settings", "mappings", "constants", "identifiers"} {
		if d.HasChange(k) {
			return true
		}
	}
	return false
}

func buildCreateRequest(d *schema.ResourceData) (*retl.CreateRETLConnectionRequest, error) {
	schedule, err := scheduleFromState(d)
	if err != nil {
		return nil, err
	}

	enabled := d.Get("enabled").(bool)
	req := &retl.CreateRETLConnectionRequest{
		SourceID:      d.Get("source_id").(string),
		DestinationID: d.Get("destination_id").(string),
		Enabled:       &enabled,
		Schedule:      schedule,
		SyncBehaviour: retl.SyncBehaviour(d.Get("sync_behaviour").(string)),
		Identifiers:   mappingsFromState(d, "identifiers"),
	}

	if v := d.Get("external_id").(string); v != "" {
		req.ExternalID = v
	}
	if ss, ok := syncSettingsFromState(d); ok {
		req.SyncSettings = ss
	}
	if mappings := mappingsFromState(d, "mappings"); len(mappings) > 0 {
		req.Mappings = mappings
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
	if v := d.Get("destination_config").(string); v != "" {
		req.DestinationConfig = json.RawMessage(v)
	}
	return req, nil
}

func buildUpdateRequest(d *schema.ResourceData) *retl.UpdateRETLConnectionRequest {
	enabled := d.Get("enabled").(bool)
	schedule, _ := scheduleFromState(d)

	req := &retl.UpdateRETLConnectionRequest{
		Enabled:  &enabled,
		Schedule: schedule,
	}
	if ss, ok := syncSettingsFromState(d); ok {
		req.SyncSettings = ss
	}
	if d.HasChange("mappings") {
		mappings := mappingsFromState(d, "mappings")
		req.Mappings = &mappings
	}
	if d.HasChange("constants") {
		constants := constantsFromState(d)
		req.Constants = &constants
	}
	if d.HasChange("identifiers") {
		req.Identifiers = mappingsFromState(d, "identifiers")
	}
	return req
}

func storeConnectionToState(d *schema.ResourceData, c *retl.RETLConnection) error {
	d.SetId(c.ID)
	setters := []struct {
		k string
		v interface{}
	}{
		{"source_id", c.SourceID},
		{"destination_id", c.DestinationID},
		{"enabled", c.Enabled},
		{"sync_behaviour", string(c.SyncBehaviour)},
		{"schedule", scheduleToState(c.Schedule)},
		{"sync_settings", syncSettingsToState(c.SyncSettings)},
		{"identifiers", mappingsToState(c.Identifiers)},
		{"mappings", mappingsToState(c.Mappings)},
		{"event", eventToState(c.Event)},
		{"constants", constantsToState(c.Constants)},
		{"cursor_column", c.CursorColumn},
		{"object", c.Object},
	}
	for _, s := range setters {
		if err := d.Set(s.k, s.v); err != nil {
			return fmt.Errorf("set %s: %w", s.k, err)
		}
	}

	// Always set, even when empty, so state can be cleared if the server
	// returns an empty value (otherwise Terraform would see perpetual diffs
	// against the stale local value).
	if err := d.Set("external_id", c.ExternalID); err != nil {
		return err
	}
	if err := d.Set("destination_config", string(c.DestinationConfig)); err != nil {
		return err
	}
	if c.CreatedAt != nil {
		if err := d.Set("created_at", c.CreatedAt.Format(time.RFC3339)); err != nil {
			return err
		}
	}
	if c.UpdatedAt != nil {
		if err := d.Set("updated_at", c.UpdatedAt.Format(time.RFC3339)); err != nil {
			return err
		}
	}
	return nil
}

// --- per-block helpers ---

func scheduleFromState(d *schema.ResourceData) (retl.Schedule, error) {
	raw, ok := d.Get("schedule").([]interface{})
	if !ok || len(raw) == 0 || raw[0] == nil {
		return retl.Schedule{}, fmt.Errorf("schedule block is required")
	}
	m := raw[0].(map[string]interface{})
	s := retl.Schedule{Type: retl.ScheduleType(m["type"].(string))}
	if v, ok := m["every_minutes"].(int); ok && v > 0 {
		s.EveryMinutes = &v
	}
	if s.Type == retl.ScheduleTypeBasic && s.EveryMinutes == nil {
		return retl.Schedule{}, fmt.Errorf("schedule.every_minutes is required when schedule.type is %q", s.Type)
	}
	return s, nil
}

func scheduleToState(s retl.Schedule) []map[string]interface{} {
	m := map[string]interface{}{"type": string(s.Type)}
	if s.EveryMinutes != nil {
		m["every_minutes"] = *s.EveryMinutes
	}
	return []map[string]interface{}{m}
}

func syncSettingsFromState(d *schema.ResourceData) (*retl.SyncSettings, bool) {
	raw, ok := d.Get("sync_settings").([]interface{})
	if !ok || len(raw) == 0 || raw[0] == nil {
		return nil, false
	}
	m := raw[0].(map[string]interface{})
	ss := &retl.SyncSettings{}

	if logs, ok := m["sync_logs_config"].([]interface{}); ok && len(logs) > 0 && logs[0] != nil {
		lm := logs[0].(map[string]interface{})
		cfg := &retl.SyncLogsConfig{}
		if v, ok := lm["enabled"].(bool); ok {
			cfg.Enabled = &v
		}
		if v, ok := lm["log_retention_in_days"].(int); ok && v > 0 {
			cfg.LogRetentionInDays = &v
		}
		if v, ok := lm["snapshots_to_retain"].(int); ok && v > 0 {
			cfg.SnapshotsToRetain = &v
		}
		ss.SyncLogsConfig = cfg
	}
	if fk, ok := m["failed_keys_config"].([]interface{}); ok && len(fk) > 0 && fk[0] != nil {
		fkm := fk[0].(map[string]interface{})
		cfg := &retl.FailedKeysConfig{}
		if v, ok := fkm["enable_failed_keys_retry"].(bool); ok {
			cfg.EnableFailedKeysRetry = &v
		}
		ss.FailedKeysConfig = cfg
	}
	return ss, true
}

func syncSettingsToState(ss *retl.SyncSettings) []map[string]interface{} {
	if ss == nil {
		return nil
	}
	out := map[string]interface{}{}
	if ss.SyncLogsConfig != nil {
		lm := map[string]interface{}{}
		if ss.SyncLogsConfig.Enabled != nil {
			lm["enabled"] = *ss.SyncLogsConfig.Enabled
		}
		if ss.SyncLogsConfig.LogRetentionInDays != nil {
			lm["log_retention_in_days"] = *ss.SyncLogsConfig.LogRetentionInDays
		}
		if ss.SyncLogsConfig.SnapshotsToRetain != nil {
			lm["snapshots_to_retain"] = *ss.SyncLogsConfig.SnapshotsToRetain
		}
		out["sync_logs_config"] = []map[string]interface{}{lm}
	}
	if ss.FailedKeysConfig != nil {
		fkm := map[string]interface{}{}
		if ss.FailedKeysConfig.EnableFailedKeysRetry != nil {
			fkm["enable_failed_keys_retry"] = *ss.FailedKeysConfig.EnableFailedKeysRetry
		}
		out["failed_keys_config"] = []map[string]interface{}{fkm}
	}
	return []map[string]interface{}{out}
}

func mappingsFromState(d *schema.ResourceData, key string) []retl.Mapping {
	raw, _ := d.Get(key).([]interface{})
	out := make([]retl.Mapping, 0, len(raw))
	for _, item := range raw {
		m := item.(map[string]interface{})
		out = append(out, retl.Mapping{
			From: m["from"].(string),
			To:   m["to"].(string),
		})
	}
	return out
}

func mappingsToState(ms []retl.Mapping) []map[string]interface{} {
	if len(ms) == 0 {
		return nil
	}
	out := make([]map[string]interface{}, 0, len(ms))
	for _, m := range ms {
		out = append(out, map[string]interface{}{"from": m.From, "to": m.To})
	}
	return out
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
