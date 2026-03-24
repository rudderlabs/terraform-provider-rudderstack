package rudderstack

import (
	"context"

	"github.com/rudderlabs/rudder-api-go/client"
)

type Client struct {
	Sources         SourcesService
	Destinations    DestinationsService
	Connections     ConnectionsService
	RETLSources     RETLSourcesService
	RETLConnections RETLConnectionsService
	Accounts        AccountsService
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

type RETLSourcesService interface {
	Create(ctx context.Context, source *client.RETLSource) (*client.RETLSource, error)
	Get(ctx context.Context, id string) (*client.RETLSource, error)
	Update(ctx context.Context, source *client.RETLSource) (*client.RETLSource, error)
	Delete(ctx context.Context, id string) error
}

type RETLConnectionsService interface {
	Create(ctx context.Context, connection *client.RETLConnection) (*client.RETLConnection, error)
	Get(ctx context.Context, id string) (*client.RETLConnection, error)
	Update(ctx context.Context, connection *client.RETLConnection) (*client.RETLConnection, error)
	Delete(ctx context.Context, id string) error
}

type AccountsService interface {
	Create(ctx context.Context, input *client.AccountCreateInput) (*client.Account, error)
	Get(ctx context.Context, id string) (*client.Account, error)
	Update(ctx context.Context, id string, input *client.AccountUpdateInput) (*client.Account, error)
	Delete(ctx context.Context, id string) error
}

func NewAPIClient(accessToken string, options ...client.Option) (*Client, error) {
	api, err := client.New(accessToken, options...)
	if err != nil {
		return nil, err
	}

	c := &Client{
		Sources:      api.Sources,
		Destinations: api.Destinations,
		Connections:  api.Connections,
	}

	if api.RETLSources != nil {
		c.RETLSources = api.RETLSources
	}
	if api.RETLConnections != nil {
		c.RETLConnections = api.RETLConnections
	}
	if api.Accounts != nil {
		c.Accounts = api.Accounts
	}

	return c, nil
}
