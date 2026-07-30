package configs

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTestConfigDestinationAPIResponseDefaultsToRequest(t *testing.T) {
	cfg := TestConfig{
		APICreate: `{"kept":"create","requestOnly":true}`,
		APIUpdate: `{"kept":"update","requestOnly":true}`,
	}

	require.JSONEq(t, cfg.APICreate, cfg.DestinationAPICreateResponse())
	require.JSONEq(t, cfg.APIUpdate, cfg.DestinationAPIUpdateResponse())
}

func TestTestConfigDestinationAPIResponseOverridesRequest(t *testing.T) {
	cfg := TestConfig{
		APICreate:         `{"kept":"create","requestOnly":true}`,
		APICreateResponse: `{"kept":"create"}`,
		APIUpdate:         `{"kept":"update","requestOnly":true}`,
		APIUpdateResponse: `{"kept":"update"}`,
	}

	require.JSONEq(t, `{"kept":"create"}`, cfg.DestinationAPICreateResponse())
	require.JSONEq(t, `{"kept":"update"}`, cfg.DestinationAPIUpdateResponse())
}

func TestTestConfigDestinationAPIResponsePrunesRequestOnlyKeys(t *testing.T) {
	cfg := TestConfig{
		APICreate:           `{"kept":"create","requestOnly":true}`,
		APICreatePrunedKeys: []string{"requestOnly"},
		APIUpdate:           `{"kept":"update","requestOnly":true,"anotherRequestOnly":1}`,
		APIUpdatePrunedKeys: []string{"requestOnly", "anotherRequestOnly"},
	}

	require.JSONEq(t, `{"kept":"create"}`, cfg.DestinationAPICreateResponse())
	require.JSONEq(t, `{"kept":"update"}`, cfg.DestinationAPIUpdateResponse())

	var create map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(cfg.APICreate), &create))
	require.Contains(t, create, "requestOnly")
}
