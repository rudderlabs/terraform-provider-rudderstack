package rudderclient

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// Auth kinds using in RudderStack APIs.
type AuthKind int

const (
	BasicAuth AuthKind = iota
	TokenAuth          = iota
)

// Client -
type Client struct {
	HTTPClient    *http.Client
	WorkspaceHost HostAccessStruct
	CatalogHost   HostAccessStruct
}

// HostAccessStruct -
type HostAccessStruct struct {
	HTTPClient *http.Client
	Url        string   `json:"hosturl"`
	Token      string   `json:"token"`
	AuthKind   AuthKind `json:"authKind"`
}

// NewClient -
func NewClient(workspaceHost, workspaceToken, catalogHost, catalogToken *string) (*Client, error) {
	httpClient := &http.Client{Timeout: 10 * time.Second}
	c := Client{
		HTTPClient: httpClient,
		WorkspaceHost: HostAccessStruct{
			HTTPClient: httpClient,
			Url:        *workspaceHost,
			Token:      *workspaceToken,
			AuthKind:   TokenAuth,
		},
		CatalogHost: HostAccessStruct{
			HTTPClient: httpClient,
			Url:        *catalogHost,
			Token:      *catalogToken,
			AuthKind:   BasicAuth,
		},
	}

	return &c, nil
}

func (ha *HostAccessStruct) doRequest(req *http.Request) ([]byte, error) {
	if ha.AuthKind == BasicAuth {
		req.SetBasicAuth(ha.Token, "")
	} else {
		req.Header.Set("Authorization", "Bearer " + ha.Token)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	res, err := ha.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return body, err
}
