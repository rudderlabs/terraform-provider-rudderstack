package rudderstack

import (
	"context"

	"github.com/rudderlabs/rudder-api-go/client"
)

type Client struct {
	Sources      SourcesService
	Destinations DestinationsService
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

func NewAPIClient(accessToken string, options ...client.Option) (*Client, error) {
	api, err := client.New(accessToken, options...)
	if err != nil {
		return nil, err
	}

	return &Client{
		Sources:      api.Sources,
		Destinations: api.Destinations,
	}, nil
}
