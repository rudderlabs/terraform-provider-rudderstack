package retl

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/rudderlabs/rudder-iac/api/client/retl"
)

// baseConnectionSchema returns the fields shared by every
// rudderstack_retl_connection_* resource: the universal source / destination /
// schedule / sync / cursor fields, but NOT destination-flow-specific fields
// like `event`, `object`, or `audience_id`. Each per-flow resource (the
// generic JSON Mapper / Object Mapping resource and any typed
// destination-scoped resource such as
// rudderstack_retl_connection_customerio_audience) merges its own
// destination-specific fields onto this base via mergeSchemas.
//
// cursor_column lives here because it's a generic source-side field with a
// single cross-flow rule (sync_behaviour="upsert"); each resource enforces
// that rule via validateCursorColumnUpsertOnly in its CustomizeDiff.
//
// Identifiers are ForceNew at both the top level (catches list-size changes)
// and the nested from/to attributes (catches in-place value mutations). This
// applies uniformly across all flows — every RETL flow treats identifier
// changes as breaking and requires a new connection.
func baseConnectionSchema() map[string]*schema.Schema {
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
					"cron_expression": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Cron expression. Required when `type` is `cron`.",
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
			ForceNew: true,
			MinItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					// Nested ForceNew is required: the top-level ForceNew on a
					// TypeList triggers replacement only on list-size changes
					// (add/remove), not when an existing element's value
					// mutates in place.
					"from": {Type: schema.TypeString, Required: true, ForceNew: true},
					"to":   {Type: schema.TypeString, Required: true, ForceNew: true},
				},
			},
			Description: "Source-to-destination identifier mappings. ForceNew: any change recreates the connection.",
		},
		// cursor_column is a generic source-side field (the incremental
		// watermark column), shared by every flow. Its only constraint is
		// sync_behaviour="upsert", enforced uniformly via
		// validateCursorColumnUpsertOnly in each resource's CustomizeDiff.
		"cursor_column": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "Column name for incremental upsert syncs (only valid when sync_behaviour is `upsert`).",
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

// validateCursorColumnUpsertOnly enforces the one cross-flow rule on
// cursor_column: it's only valid when sync_behaviour is "upsert". Shared by
// every resource's CustomizeDiff so the error surfaces at plan time instead of
// an API rejection on apply. Takes resourceGetter so it works from
// *schema.ResourceDiff and is trivially testable.
func validateCursorColumnUpsertOnly(d resourceGetter) error {
	cursor, _ := d.Get("cursor_column").(string)
	if cursor == "" {
		return nil
	}
	if sb, _ := d.Get("sync_behaviour").(string); sb != "" && sb != "upsert" {
		return fmt.Errorf("cursor_column is only valid when sync_behaviour is %q, got %q", "upsert", sb)
	}
	return nil
}

// mergeSchemas combines a base map with overrides. Overrides win on key
// conflict — typed resources can replace a base field (e.g. tighten
// `identifiers` to ForceNew, or relax an unsupported field).
func mergeSchemas(base, override map[string]*schema.Schema) map[string]*schema.Schema {
	out := make(map[string]*schema.Schema, len(base)+len(override))
	for k, v := range base {
		out[k] = v
	}
	for k, v := range override {
		out[k] = v
	}
	return out
}

// applyBaseToCreateRequest populates the universal fields of a
// CreateRETLConnectionRequest from terraform state. Flow-specific fields
// (Event, Constants, Mappings, Object, DestinationConfig) are the caller's
// responsibility — e.g. `mappings` (field mappings) lives only on the generic
// resource, not on the destination-specific flows.
//
// cursor_column is handled here rather than per-resource: it's a generic
// source-side field, not destination-specific. GetOk is panic-safe on a
// resource that doesn't declare the field (e.g. the audience resource), so
// this is a no-op there.
func applyBaseToCreateRequest(d *schema.ResourceData, req *retl.CreateRETLConnectionRequest) error {
	schedule, err := scheduleFromState(d)
	if err != nil {
		return err
	}
	enabled := d.Get("enabled").(bool)
	req.SourceID = d.Get("source_id").(string)
	req.DestinationID = d.Get("destination_id").(string)
	req.Enabled = &enabled
	req.Schedule = schedule
	req.SyncBehaviour = retl.SyncBehaviour(d.Get("sync_behaviour").(string))
	req.Identifiers = mappingsFromState(d, "identifiers")

	if ss, ok := syncSettingsFromState(d); ok {
		req.SyncSettings = ss
	}
	if v, ok := d.GetOk("cursor_column"); ok {
		req.CursorColumn = v.(string)
	}
	return nil
}

// applyBaseToUpdateRequest populates the universal fields of an
// UpdateRETLConnectionRequest from terraform state, respecting d.HasChange so
// unchanged fields stay nil and don't get echoed back to the API.
func applyBaseToUpdateRequest(d *schema.ResourceData, req *retl.UpdateRETLConnectionRequest) error {
	schedule, err := scheduleFromState(d)
	if err != nil {
		return err
	}
	enabled := d.Get("enabled").(bool)
	req.Enabled = &enabled
	req.Schedule = schedule

	// Only forward sync_settings when the user actually changed the block.
	// Otherwise the Optional+Computed field can carry server-computed values
	// in state and we'd echo them back as if the user had set them.
	if d.HasChange("sync_settings") {
		if ss, ok := syncSettingsFromState(d); ok {
			req.SyncSettings = ss
		}
	}
	// Identifiers are ForceNew across all flows (see baseConnectionSchema),
	// so HasChange("identifiers") never fires on Update — terraform routes
	// identifier changes through destroy + create instead.
	return nil
}

// storeBaseConnectionToState writes the universal fields back to terraform
// state. Destination-flow-specific fields (event, constants, object,
// audience_id, ...) are written by the caller after this returns.
//
// cursor_column is a generic source-side field declared in
// baseConnectionSchema, so every connection resource has it and it is written
// back here unconditionally.
func storeBaseConnectionToState(d *schema.ResourceData, c *retl.RETLConnection) error {
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
	}
	for _, s := range setters {
		if err := d.Set(s.k, s.v); err != nil {
			return fmt.Errorf("set %s: %w", s.k, err)
		}
	}
	if err := d.Set("cursor_column", c.CursorColumn); err != nil {
		return fmt.Errorf("set cursor_column: %w", err)
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

// --- per-block helpers (shared by all retl connection resources) ---

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
	if v, _ := m["cron_expression"].(string); v != "" {
		s.CronExpression = &v
	}
	if s.Type == retl.ScheduleTypeBasic && s.EveryMinutes == nil {
		return retl.Schedule{}, fmt.Errorf("schedule.every_minutes is required when schedule.type is %q", s.Type)
	}
	if s.Type == retl.ScheduleTypeCron && s.CronExpression == nil {
		return retl.Schedule{}, fmt.Errorf("schedule.cron_expression is required when schedule.type is %q", s.Type)
	}
	return s, nil
}

func scheduleToState(s retl.Schedule) []map[string]interface{} {
	m := map[string]interface{}{"type": string(s.Type)}
	if s.EveryMinutes != nil {
		m["every_minutes"] = *s.EveryMinutes
	}
	if s.CronExpression != nil {
		m["cron_expression"] = *s.CronExpression
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
	// Empty sync_settings {} block (or one with empty sub-blocks) carries
	// nothing useful — treat it as "not set" so callers can omit the field.
	if ss.SyncLogsConfig == nil && ss.FailedKeysConfig == nil {
		return nil, false
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
