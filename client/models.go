package rudderclient

import (
    "time"
)

// Source Definition -
type SourceDefinition struct
{
    ID            int                         `json:"id,omitempty"`
    Name          string                      `json:"name"`
    Category      string                      `json:"string"`
    CreatedAt     time.Time                   `json:"createdAt"`
    UpdatedAt     time.Time                   `json:"updatedAt"`

    Config        SourceDefinitionConfig      `json:"config"`
}

// Destination Definition -
type DestinationDefinition struct
{
    ID            int                         `json:"id,omitempty"`
    Name          string                      `json:"name"`
    Category      string                      `json:"string"`
    CreatedAt     time.Time                   `json:"createdAt"`
    UpdatedAt     time.Time                   `json:"updatedAt"`

    Config        DestinationDefinitionConfig `json:"config"`
}

// Sources -
type Source struct {
    ID            int                         `json:"id,omitempty"`
    Name          string                      `json:"name"`
    Category      string                      `json:"category"`
    CreatedAt     time.Time                   `json:"createdAt"`
    UpdatedAt     time.Time                   `json:"updatedAt"`

    Config        SourceConfig                `json:"config"`
}

// Destinations -
type Destination struct {
    ID            int                         `json:"id,omitempty"`
    Name          string                      `json:"name"`
    Category      string                      `json:"category"`
    CreatedAt     time.Time                   `json:"createdAt"`
    UpdatedAt     time.Time                   `json:"updatedAt"`

    Config        DestinationConfig           `json:"config"`
}

type SourceConfig struct {
}

type DestinationConfig struct {
}

type SourceDefinitionConfig struct {
}

type DestinationDefinitionConfig struct {
}
