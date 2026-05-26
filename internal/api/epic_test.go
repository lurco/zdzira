package api_test

import (
	"net/http"
	"testing"
	"zdzira/internal/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEpicCreate_Returns201(t *testing.T) {
	srv := newTestServer(t)
	defer srv.Close()

	resp := do(t, srv, http.MethodPost, "/projects", map[string]string{"name": "Epic Proj", "shortcut": "EP"})
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	resp.Body.Close()

	var e model.Epic
	resp = do(t, srv, http.MethodPost, "/projects/epic-proj/epics", map[string]string{"name": "big feature"})
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	decode(t, resp, &e)
	assert.Equal(t, uint(1), e.Number)
	assert.Equal(t, "big feature", e.Name)
}

func TestEpicGet_Returns200(t *testing.T) {
	srv := newTestServer(t)
	defer srv.Close()

	resp := do(t, srv, http.MethodPost, "/projects", map[string]string{"name": "Epic Proj", "shortcut": "EP"})
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	resp.Body.Close()

	resp = do(t, srv, http.MethodPost, "/projects/epic-proj/epics", map[string]string{"name": "the one"})
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	resp.Body.Close()

	var e model.Epic
	resp = do(t, srv, http.MethodGet, "/projects/epic-proj/epics/EP-E1", nil)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	decode(t, resp, &e)
	assert.Equal(t, "the one", e.Name)
}

func TestEpicList_Returns200(t *testing.T) {
	srv := newTestServer(t)
	defer srv.Close()

	resp := do(t, srv, http.MethodPost, "/projects", map[string]string{"name": "Epic Proj", "shortcut": "EP"})
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	resp.Body.Close()

	for _, name := range []string{"epic one", "epic two"} {
		resp = do(t, srv, http.MethodPost, "/projects/epic-proj/epics", map[string]string{"name": name})
		require.Equal(t, http.StatusCreated, resp.StatusCode)
		resp.Body.Close()
	}

	var epics []model.Epic
	resp = do(t, srv, http.MethodGet, "/projects/epic-proj/epics", nil)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	decode(t, resp, &epics)
	assert.Len(t, epics, 2)
}

func TestEpicDelete_Returns204(t *testing.T) {
	srv := newTestServer(t)
	defer srv.Close()

	resp := do(t, srv, http.MethodPost, "/projects", map[string]string{"name": "Epic Proj", "shortcut": "EP"})
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	resp.Body.Close()

	resp = do(t, srv, http.MethodPost, "/projects/epic-proj/epics", map[string]string{"name": "to delete"})
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	resp.Body.Close()

	resp = do(t, srv, http.MethodDelete, "/projects/epic-proj/epics/EP-E1", nil)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	resp.Body.Close()
}
