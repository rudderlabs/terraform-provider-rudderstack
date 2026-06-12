// TEMPORARY — remove when DEX-375 lands; types mirror CONTRACT-ACCT-V1 §9.3.

package rudderstack

import (
	"encoding/json"
	"time"
)

type CreateAccountRequest struct {
	Name                  string          `json:"name"`
	AccountDefinitionName string          `json:"accountDefinitionName"`
	Options               json.RawMessage `json:"options"`
	Secret                json.RawMessage `json:"secret"`
}

// UpdateAccountRequest intentionally omits accountDefinitionName — it is immutable and must not be sent on update.
type UpdateAccountRequest struct {
	Name    string          `json:"name"`
	Options json.RawMessage `json:"options"`
	Secret  json.RawMessage `json:"secret"`
}

type AccountDefinition struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Category string `json:"category"`
}

type Account struct {
	ID         string            `json:"id"`
	Name       string            `json:"name"`
	Definition AccountDefinition `json:"definition"`
	Options    json.RawMessage   `json:"options"`
	CreatedAt  *time.Time        `json:"createdAt"`
	UpdatedAt  *time.Time        `json:"updatedAt"`
}
