package retl_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/rudderlabs/rudder-iac/api/client"
	iacretl "github.com/rudderlabs/rudder-iac/api/client/retl"
	"github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/retl"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) Do(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestServiceCreateConnectionOmittingSyncBehaviourDoesNotSendSyncBehaviour(t *testing.T) {
	httpClient := roundTripFunc(func(req *http.Request) (*http.Response, error) {
		require.Equal(t, "POST", req.Method)
		require.Equal(t, "https://api.rudderstack.com/v2/retl-connections", req.URL.String())

		body, err := io.ReadAll(req.Body)
		require.NoError(t, err)

		var payload map[string]interface{}
		require.NoError(t, json.Unmarshal(body, &payload))
		require.NotContains(t, payload, "syncBehaviour")
		require.Equal(t, map[string]interface{}{"object": "event"}, payload["destinationConfig"])

		return &http.Response{
			StatusCode: http.StatusOK,
			Body: io.NopCloser(strings.NewReader(`{
				"id": "conn-cio",
				"sourceId": "retl-src-1",
				"destinationId": "dest-cio",
				"enabled": true,
				"schedule": {"type": "manual"},
				"syncBehaviour": "upsert",
				"identifiers": [{"from": "email", "to": "email"}],
				"destinationConfig": {"object": "event"}
			}`)),
		}, nil
	})

	api, err := client.New("test-token", client.WithHTTPClient(httpClient))
	require.NoError(t, err)

	svc := retl.NewService(api)
	creator, ok := svc.(interface {
		CreateConnectionOmittingSyncBehaviour(context.Context, *iacretl.CreateRETLConnectionRequest) (*iacretl.RETLConnection, error)
	})
	require.True(t, ok)

	created, err := creator.CreateConnectionOmittingSyncBehaviour(context.Background(), &iacretl.CreateRETLConnectionRequest{
		SourceID:          "retl-src-1",
		DestinationID:     "dest-cio",
		Schedule:          iacretl.Schedule{Type: iacretl.ScheduleTypeManual},
		SyncBehaviour:     iacretl.SyncBehaviourUpsert,
		Identifiers:       []iacretl.Mapping{{From: "email", To: "email"}},
		DestinationConfig: json.RawMessage(`{"object":"event"}`),
	})
	require.NoError(t, err)
	require.Equal(t, "conn-cio", created.ID)
}
