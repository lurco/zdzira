package api_test

import (
	"net/http"
	"testing"
	"zdzira/internal/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIssueUpdate_ChangesFieldsViaAPI(t *testing.T) {
	srv := newTestServer(t)
	defer srv.Close()

	resp := do(t, srv, http.MethodPost, "/projects", map[string]string{
		"name":     "API Test",
		"shortcut": "API",
	})
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	resp.Body.Close()

	resp = do(t, srv, http.MethodPost, "/projects/api-test/issues", map[string]any{
		"name":     "original",
		"type":     "TASK",
		"priority": "LOW",
	})
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	resp.Body.Close()

	var updated model.Issue
	resp = do(t, srv, http.MethodPut, "/projects/api-test/issues/API-1", map[string]any{
		"name":     "updated via api",
		"type":     "BUG",
		"priority": "HIGH",
	})
	require.Equal(t, http.StatusOK, resp.StatusCode)
	decode(t, resp, &updated)

	assert.Equal(t, "updated via api", updated.Name)
	assert.Equal(t, model.IssueType("BUG"), updated.Type)
	assert.Equal(t, model.Priority("HIGH"), updated.Priority)
}
