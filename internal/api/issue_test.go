package api_test

import (
	"net/http"
	"testing"
	"zdzira/internal/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIssueCreate_Returns201(t *testing.T) {
	srv := newTestServer(t)
	defer srv.Close()

	resp := do(t, srv, http.MethodPost, "/projects", map[string]string{"name": "Issue Proj", "shortcut": "IS"})
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	resp.Body.Close()

	var issue model.Issue
	resp = do(t, srv, http.MethodPost, "/projects/issue-proj/issues", map[string]any{
		"name": "first issue", "type": "TASK", "priority": "LOW",
	})
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	decode(t, resp, &issue)
	assert.Equal(t, uint(1), issue.Number)
	assert.Equal(t, "first issue", issue.Name)
}

func TestIssueGet_Returns200(t *testing.T) {
	srv := newTestServer(t)
	defer srv.Close()

	resp := do(t, srv, http.MethodPost, "/projects", map[string]string{"name": "Issue Proj", "shortcut": "IS"})
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	resp.Body.Close()

	resp = do(t, srv, http.MethodPost, "/projects/issue-proj/issues", map[string]any{
		"name": "find me", "type": "TASK", "priority": "LOW",
	})
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	resp.Body.Close()

	var issue model.Issue
	resp = do(t, srv, http.MethodGet, "/projects/issue-proj/issues/IS-1", nil)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	decode(t, resp, &issue)
	assert.Equal(t, "find me", issue.Name)
}

func TestIssueList_Returns200(t *testing.T) {
	srv := newTestServer(t)
	defer srv.Close()

	resp := do(t, srv, http.MethodPost, "/projects", map[string]string{"name": "Issue Proj", "shortcut": "IS"})
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	resp.Body.Close()

	for _, name := range []string{"one", "two", "three"} {
		resp = do(t, srv, http.MethodPost, "/projects/issue-proj/issues", map[string]any{
			"name": name, "type": "TASK", "priority": "LOW",
		})
		require.Equal(t, http.StatusCreated, resp.StatusCode)
		resp.Body.Close()
	}

	var issues []model.Issue
	resp = do(t, srv, http.MethodGet, "/projects/issue-proj/issues", nil)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	decode(t, resp, &issues)
	assert.Len(t, issues, 3)
}

func TestIssueMove_Returns200(t *testing.T) {
	srv := newTestServer(t)
	defer srv.Close()

	resp := do(t, srv, http.MethodPost, "/projects", map[string]string{"name": "Issue Proj", "shortcut": "IS"})
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	resp.Body.Close()

	resp = do(t, srv, http.MethodPost, "/projects/issue-proj/issues", map[string]any{
		"name": "move me", "type": "TASK", "priority": "LOW",
	})
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	resp.Body.Close()

	var moved model.Issue
	resp = do(t, srv, http.MethodPost, "/projects/issue-proj/issues/IS-1/move", map[string]string{
		"swimlane": "Done",
	})
	require.Equal(t, http.StatusOK, resp.StatusCode)
	decode(t, resp, &moved)
	assert.NotZero(t, moved.SwimlaneID)
}

func TestIssueDelete_Returns204(t *testing.T) {
	srv := newTestServer(t)
	defer srv.Close()

	resp := do(t, srv, http.MethodPost, "/projects", map[string]string{"name": "Issue Proj", "shortcut": "IS"})
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	resp.Body.Close()

	resp = do(t, srv, http.MethodPost, "/projects/issue-proj/issues", map[string]any{
		"name": "delete me", "type": "TASK", "priority": "LOW",
	})
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	resp.Body.Close()

	resp = do(t, srv, http.MethodDelete, "/projects/issue-proj/issues/IS-1", nil)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	resp.Body.Close()
}

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
