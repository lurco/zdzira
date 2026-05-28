package api_test

import (
	"net/http"
	"testing"
	"zdzira/backend/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProjectCreate_Returns201(t *testing.T) {
	srv := newTestServer(t)
	defer srv.Close()

	var p model.Project
	resp := do(t, srv, http.MethodPost, "/projects", map[string]string{
		"name":     "My Project",
		"shortcut": "MP",
	})
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	decode(t, resp, &p)

	assert.Equal(t, "my-project", p.Slug)
	assert.Equal(t, "MP", p.Shortcut)
}

func TestProjectGet_Returns200(t *testing.T) {
	srv := newTestServer(t)
	defer srv.Close()

	resp := do(t, srv, http.MethodPost, "/projects", map[string]string{"name": "Get Me", "shortcut": "GM"})
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	resp.Body.Close()

	var p model.Project
	resp = do(t, srv, http.MethodGet, "/projects/get-me", nil)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	decode(t, resp, &p)
	assert.Equal(t, "get-me", p.Slug)
}

func TestProjectList_Returns200(t *testing.T) {
	srv := newTestServer(t)
	defer srv.Close()

	for _, name := range []string{"Alpha", "Beta"} {
		resp := do(t, srv, http.MethodPost, "/projects", map[string]string{"name": name, "shortcut": name[:3]})
		require.Equal(t, http.StatusCreated, resp.StatusCode)
		resp.Body.Close()
	}

	var projects []model.Project
	resp := do(t, srv, http.MethodGet, "/projects", nil)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	decode(t, resp, &projects)
	assert.Len(t, projects, 2)
}

func TestProjectDelete_Returns204(t *testing.T) {
	srv := newTestServer(t)
	defer srv.Close()

	resp := do(t, srv, http.MethodPost, "/projects", map[string]string{"name": "Delete Me", "shortcut": "DM"})
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	resp.Body.Close()

	resp = do(t, srv, http.MethodDelete, "/projects/delete-me", nil)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	resp.Body.Close()
}
