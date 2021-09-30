package rudderclient

import (
	"encoding/json"
	"fmt"
	// "log"
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

type srcResultBodyType struct{
	Source Source `json:"source"`
}

// CreateSource - Create new source
func (c *Client) CreateSource(source Source) (*Source, error) {
	host := c.WorkspaceHost
	rb, err := json.Marshal(source)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%ssources/", host.Url)
	req, err := http.NewRequest("POST", url, strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := host.doRequest(req)
	if err != nil {
		return nil, err
	}

	resultBody := srcResultBodyType{}
	err = json.Unmarshal(body, &resultBody)
	if err != nil {
		return nil, err
	}

	return &resultBody.Source, nil
}

// DeleteSource - Delete existing source
func (c *Client) DeleteSource(sourceId string) error {
	host := c.WorkspaceHost

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/sources/%d", host.Url, sourceId), nil)
	if err != nil {
		return err
	}

	body, err := host.doRequest(req)
	_ = body
	if err != nil {
		return err
	}

	return nil
}
