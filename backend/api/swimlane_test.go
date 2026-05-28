package api_test

import (
	"fmt"
	"net/http"
	"testing"
	"zdzira/backend/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSwimlaneList_Returns200(t *testing.T) {
	srv := newTestServer(t)
	defer srv.Close()

	resp := do(t, srv, http.MethodPost, "/projects", map[string]string{"name": "Board", "shortcut": "BRD"})
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	resp.Body.Close()

	var swimlanes []model.Swimlane
	resp = do(t, srv, http.MethodGet, "/projects/board/swimlanes", nil)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	decode(t, resp, &swimlanes)
	assert.Len(t, swimlanes, 3, "seeded Backlog / In Progress / Done")
	assert.Equal(t, "Backlog", swimlanes[0].Name)
}

func TestSwimlaneCreate_Returns201(t *testing.T) {
	srv := newTestServer(t)
	defer srv.Close()

	resp := do(t, srv, http.MethodPost, "/projects", map[string]string{"name": "Board", "shortcut": "BRD"})
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	resp.Body.Close()

	var sl model.Swimlane
	resp = do(t, srv, http.MethodPost, "/projects/board/swimlanes", map[string]string{"name": "In Review"})
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	decode(t, resp, &sl)
	assert.Equal(t, "In Review", sl.Name)
	assert.Equal(t, uint(4), sl.Position)
}

func TestSwimlaneRename_Returns200(t *testing.T) {
	srv := newTestServer(t)
	defer srv.Close()

	resp := do(t, srv, http.MethodPost, "/projects", map[string]string{"name": "Board", "shortcut": "BRD"})
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	resp.Body.Close()

	var swimlanes []model.Swimlane
	resp = do(t, srv, http.MethodGet, "/projects/board/swimlanes", nil)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	decode(t, resp, &swimlanes)
	backlogID := swimlanes[0].ID

	var sl model.Swimlane
	resp = do(t, srv, http.MethodPut,
		fmt.Sprintf("/projects/board/swimlanes/%d", backlogID),
		map[string]string{"name": "Todo"})
	require.Equal(t, http.StatusOK, resp.StatusCode)
	decode(t, resp, &sl)
	assert.Equal(t, "Todo", sl.Name)
}
