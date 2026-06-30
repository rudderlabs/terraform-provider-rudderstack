package retl

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/rudderlabs/rudder-iac/api/client"
	iacretl "github.com/rudderlabs/rudder-iac/api/client/retl"
)

const retlConnectionsBasePath = "/v2/retl-connections"

type rawRETLClient interface {
	Do(ctx context.Context, method, path string, body io.Reader) ([]byte, error)
}

type syncBehaviourOmittingConnectionCreator interface {
	CreateConnectionOmittingSyncBehaviour(ctx context.Context, req *iacretl.CreateRETLConnectionRequest) (*iacretl.RETLConnection, error)
}

type serviceWithCustomCreate struct {
	Service
	client rawRETLClient
}

// NewService wraps the upstream RETL store with provider-local request shapes
// that are not yet available in rudder-iac.
func NewService(api *client.Client) Service {
	return &serviceWithCustomCreate{
		Service: iacretl.NewRudderRETLStore(api),
		client:  api,
	}
}

type createRETLConnectionWithoutSyncBehaviourRequest struct {
	SourceID          string                `json:"sourceId"`
	DestinationID     string                `json:"destinationId"`
	Enabled           *bool                 `json:"enabled,omitempty"`
	ExternalID        string                `json:"externalId,omitempty"`
	Schedule          iacretl.Schedule      `json:"schedule"`
	SyncSettings      *iacretl.SyncSettings `json:"syncSettings,omitempty"`
	Identifiers       []iacretl.Mapping     `json:"identifiers"`
	Mappings          []iacretl.Mapping     `json:"mappings,omitempty"`
	Event             *iacretl.Event        `json:"event,omitempty"`
	Constants         []iacretl.Constant    `json:"constants,omitempty"`
	CursorColumn      string                `json:"cursorColumn,omitempty"`
	Object            string                `json:"object,omitempty"`
	DestinationConfig json.RawMessage       `json:"destinationConfig,omitempty"`
}

func (s *serviceWithCustomCreate) CreateConnectionOmittingSyncBehaviour(ctx context.Context, req *iacretl.CreateRETLConnectionRequest) (*iacretl.RETLConnection, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}
	if req.Schedule.Type == "" {
		return nil, fmt.Errorf("schedule.type is required")
	}

	body, err := json.Marshal(createRETLConnectionWithoutSyncBehaviourRequest{
		SourceID:          req.SourceID,
		DestinationID:     req.DestinationID,
		Enabled:           req.Enabled,
		ExternalID:        req.ExternalID,
		Schedule:          req.Schedule,
		SyncSettings:      req.SyncSettings,
		Identifiers:       req.Identifiers,
		Mappings:          req.Mappings,
		Event:             req.Event,
		Constants:         req.Constants,
		CursorColumn:      req.CursorColumn,
		Object:            req.Object,
		DestinationConfig: req.DestinationConfig,
	})
	if err != nil {
		return nil, fmt.Errorf("marshalling connection: %w", err)
	}

	resp, err := s.client.Do(ctx, "POST", retlConnectionsBasePath, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("creating RETL connection: %w", err)
	}

	var result iacretl.RETLConnection
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("unmarshalling response: %w", err)
	}

	return &result, nil
}
