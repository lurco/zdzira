package api_test

import (
	"net/http"
	"testing"
	"zdzira/internal/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLinkCreate_Returns201(t *testing.T) {
	srv := newTestServer(t)
	defer srv.Close()

	resp := do(t, srv, http.MethodPost, "/projects", map[string]string{"name": "Lnk Proj", "shortcut": "LNK"})
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	resp.Body.Close()

	for _, name := range []string{"source", "target"} {
		resp = do(t, srv, http.MethodPost, "/projects/lnk-proj/issues", map[string]any{
			"name": name, "type": "TASK", "priority": "LOW",
		})
		require.Equal(t, http.StatusCreated, resp.StatusCode)
		resp.Body.Close()
	}

	var l model.Link
	resp = do(t, srv, http.MethodPost, "/projects/lnk-proj/issues/LNK-1/links", map[string]string{
		"target_ref": "LNK-2",
		"type":       "BLOCKS",
	})
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	decode(t, resp, &l)
	assert.Equal(t, model.LinkType("BLOCKS"), l.Type)
}

func TestLinkListForIssue_Returns200(t *testing.T) {
	srv := newTestServer(t)
	defer srv.Close()

	resp := do(t, srv, http.MethodPost, "/projects", map[string]string{"name": "Lnk Proj", "shortcut": "LNK"})
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	resp.Body.Close()

	for _, name := range []string{"one", "two", "three"} {
		resp = do(t, srv, http.MethodPost, "/projects/lnk-proj/issues", map[string]any{
			"name": name, "type": "TASK", "priority": "LOW",
		})
		require.Equal(t, http.StatusCreated, resp.StatusCode)
		resp.Body.Close()
	}

	// LNK-1 blocks LNK-2
	resp = do(t, srv, http.MethodPost, "/projects/lnk-proj/issues/LNK-1/links", map[string]string{
		"target_ref": "LNK-2", "type": "BLOCKS",
	})
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	resp.Body.Close()

	// LNK-3 blocks LNK-1 (LNK-1 is target)
	resp = do(t, srv, http.MethodPost, "/projects/lnk-proj/issues/LNK-3/links", map[string]string{
		"target_ref": "LNK-1", "type": "BLOCKS",
	})
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	resp.Body.Close()

	var links []model.Link
	resp = do(t, srv, http.MethodGet, "/projects/lnk-proj/issues/LNK-1/links", nil)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	decode(t, resp, &links)
	assert.Len(t, links, 2)
}
