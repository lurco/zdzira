package api_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuditList_ReturnsEntriesForProject(t *testing.T) {
	srv := newTestServer(t)
	defer srv.Close()

	resp := do(t, srv, http.MethodPost, "/projects", map[string]string{"name": "Audit Proj", "shortcut": "AUD"})
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	resp.Body.Close()

	resp = do(t, srv, http.MethodPost, "/projects/audit-proj/issues", map[string]any{
		"name": "first", "type": "TASK", "priority": "LOW",
	})
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	resp.Body.Close()

	resp = do(t, srv, http.MethodPost, "/projects/audit-proj/epics", map[string]string{"name": "big feature"})
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	resp.Body.Close()

	var entries []map[string]any
	resp = do(t, srv, http.MethodGet, "/projects/audit-proj/audit", nil)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	decode(t, resp, &entries)

	require.Len(t, entries, 2)
	assert.Equal(t, "created", entries[0]["action"])
	assert.Equal(t, "epic", entries[0]["entity_type"])
	assert.Equal(t, "created", entries[1]["action"])
	assert.Equal(t, "issue", entries[1]["entity_type"])
}
