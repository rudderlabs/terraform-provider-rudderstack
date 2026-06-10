package rudderstack

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"
)

// TEMPORARY — remove when DEX-375 lands; types mirror CONTRACT-ACCT-V1 §9.3.
// This seam keeps account-resource CRUD decoupled from the yet-to-land
// rudder-iac Accounts client wiring.
type accountsAPI interface {
	Create(ctx context.Context, req *CreateAccountRequest) (*Account, error)
	Get(ctx context.Context, id string) (*Account, error)
	Update(ctx context.Context, id string, req *UpdateAccountRequest) (*Account, error)
	Delete(ctx context.Context, id string) error
}

type CreateAccountRequest struct {
	Name                  string          `json:"name"`
	AccountDefinitionName string          `json:"accountDefinitionName"`
	Options               json.RawMessage `json:"options"`
	Secret                json.RawMessage `json:"secret"`
}

type UpdateAccountRequest struct {
	Name    string          `json:"name"`
	Options json.RawMessage `json:"options"`
	Secret  json.RawMessage `json:"secret"`
}

type Account struct {
	ID         string          `json:"id"`
	Name       string          `json:"name"`
	Definition AccountTypeInfo `json:"definition"`
	Options    json.RawMessage `json:"options"`
	CreatedAt  *time.Time      `json:"createdAt"`
	UpdatedAt  *time.Time      `json:"updatedAt"`
}

type AccountTypeInfo struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Category string `json:"category"`
}

var ErrAccountNotFound = errors.New("account not found")

var accountServiceByClient sync.Map

func (c *Client) accountsClient() accountsAPI {
	if c == nil {
		return nil
	}
	if svc, ok := accountServiceByClient.Load(c); ok {
		return svc.(accountsAPI)
	}
	return nil
}
