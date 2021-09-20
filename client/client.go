package rudderclient

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// HostURL - Default Hashicups URL
const HostURL string = "http://localhost:19090"

// Client -
type Client struct {
	HTTPClient     *http.Client
	WorkspaceHost  HostAccessStruct
	CatalogHost    HostAccessStruct
}

// AuthStruct -
type HostAccessStruct struct {
	HTTPClient *http.Client
	HostUrl     string `json:"hosturl"`
	Token string `json:"token"`
	AuthKind    bool `json:"authKind"`
}

// NewClient -
func NewClient(workspaceHost, workspaceToken, catalogHost, catalogToken *string) (*Client, error) {
	httpClient := &http.Client{Timeout: 10 * time.Second}
	c := Client{
		HTTPClient: httpClient,
		WorkspaceHost: HostAccessStruct{
		        HTTPClient: httpClient,
			HostUrl: *workspaceHost,
			Token: *workspaceToken,
			AuthKind: true,
		},
		CatalogHost: HostAccessStruct{
		        HTTPClient: httpClient,
			HostUrl: *catalogHost,
			Token: *catalogToken,
			AuthKind: true,
		},
	}

	return &c, nil
}

func (ha *HostAccessStruct) doRequest(req *http.Request) ([]byte, error) {
	req.Header.Set("Authorization", ha.Token)

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
