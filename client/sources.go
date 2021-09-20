package rudderclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// GetSources - Returns list of sources.
func (c *Client) GetSources() ([]Source, error) {
	host := c.WorkspaceHost
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/sources", host.Url), nil)
	if err != nil {
		return nil, err
	}

	body, err := host.doRequest(req)
	if err != nil {
		return nil, err
	}

	sources := []Source{}
	err = json.Unmarshal(body, &sources)
	if err != nil {
		return nil, err
	}

	return sources, nil
}

// GetSource - Returns source
func (c *Client) GetSource(sourceID string) (*Source, error) {
	host := c.WorkspaceHost

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/sources/%s", host.Url, sourceID), nil)
	if err != nil {
		return nil, err
	}

	body, err := host.doRequest(req)
	if err != nil {
		return nil, err
	}

	source := Source{}
	err = json.Unmarshal(body, &source)
	if err != nil {
		return nil, err
	}

	return &source, nil
}

// CreateSource - Create new source
func (c *Client) CreateSource(source Source) (*Source, error) {
	host := c.WorkspaceHost

	rb, err := json.Marshal(source)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/sources", host.Url), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := host.doRequest(req)
	if err != nil {
		return nil, err
	}

	newSource := Source{}
	err = json.Unmarshal(body, &newSource)
	if err != nil {
		return nil, err
	}

	return &newSource, nil
}
