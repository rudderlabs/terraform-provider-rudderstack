package rudderstack

import (
	"context"

	"github.com/rudderlabs/rudder-iac/api/client"
	iacretl "github.com/rudderlabs/rudder-iac/api/client/retl"

	"github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/retl"
)

type Client struct {
	Sources      SourcesService
	Destinations DestinationsService
	Connections  ConnectionsService
	RETLSources  retl.Service
}

// RETLSourcesClient satisfies retl.ClientProvider so the RETL resources in the
// retl subpackage can reach the configured service without importing this
// package (which would create a circular dependency).
func (c *Client) RETLSourcesClient() retl.Service {
	return c.RETLSources
}

type SourcesService interface {
	Create(ctx context.Context, source *client.Source) (*client.Source, error)
	Get(ctx context.Context, id string) (*client.Source, error)
	Update(ctx context.Context, source *client.Source) (*client.Source, error)
	Delete(ctx context.Context, id string) error
}

type DestinationsService interface {
	Create(ctx context.Context, destination *client.Destination) (*client.Destination, error)
	Get(ctx context.Context, id string) (*client.Destination, error)
	Update(ctx context.Context, destination *client.Destination) (*client.Destination, error)
	Delete(ctx context.Context, id string) error
}

type ConnectionsService interface {
	Create(ctx context.Context, connection *client.Connection) (*client.Connection, error)
	Get(ctx context.Context, id string) (*client.Connection, error)
	Update(ctx context.Context, connection *client.Connection) (*client.Connection, error)
	Delete(ctx context.Context, id string) error
}

func NewAPIClient(accessToken string, options ...client.Option) (*Client, error) {
	api, err := client.New(accessToken, options...)
	if err != nil {
		return nil, err
	}

	return &Client{
		Sources:      api.Sources,
		Destinations: api.Destinations,
		Connections:  api.Connections,
		RETLSources:  iacretl.NewRudderRETLStore(api),
	}, nil
}
