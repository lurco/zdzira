package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"zdzira/backend/api"
	"zdzira/backend/service"
	"zdzira/backend/store"

	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	db, err := store.Open(":memory:")
	require.NoError(t, err)
	return httptest.NewServer(api.NewRouter(service.New(store.New(db)), slog.Default()))
}

func do(t *testing.T, srv *httptest.Server, method, path string, body any) *http.Response {
	t.Helper()
	var b bytes.Buffer
	if body != nil {
		require.NoError(t, json.NewEncoder(&b).Encode(body))
	}
	req, err := http.NewRequestWithContext(context.Background(), method, srv.URL+path, &b)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	return resp
}

func decode(t *testing.T, resp *http.Response, v any) {
	t.Helper()
	defer resp.Body.Close()
	require.NoError(t, json.NewDecoder(resp.Body).Decode(v))
}
