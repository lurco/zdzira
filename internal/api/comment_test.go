package api_test

import (
	"net/http"
	"testing"
	"zdzira/internal/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommentAddToIssue_Returns201(t *testing.T) {
	srv := newTestServer(t)
	defer srv.Close()

	resp := do(t, srv, http.MethodPost, "/projects", map[string]string{"name": "Cmt Proj", "shortcut": "CMT"})
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	resp.Body.Close()

	resp = do(t, srv, http.MethodPost, "/projects/cmt-proj/issues", map[string]any{
		"name": "issue", "type": "TASK", "priority": "LOW",
	})
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	resp.Body.Close()

	var c model.Comment
	resp = do(t, srv, http.MethodPost, "/projects/cmt-proj/issues/CMT-1/comments", map[string]string{
		"contents": "looks good",
	})
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	decode(t, resp, &c)
	assert.Equal(t, "looks good", c.Contents)
	assert.NotNil(t, c.IssueID)
}

func TestCommentListForIssue_Returns200(t *testing.T) {
	srv := newTestServer(t)
	defer srv.Close()

	resp := do(t, srv, http.MethodPost, "/projects", map[string]string{"name": "Cmt Proj", "shortcut": "CMT"})
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	resp.Body.Close()

	resp = do(t, srv, http.MethodPost, "/projects/cmt-proj/issues", map[string]any{
		"name": "issue", "type": "TASK", "priority": "LOW",
	})
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	resp.Body.Close()

	for _, text := range []string{"first", "second"} {
		resp = do(t, srv, http.MethodPost, "/projects/cmt-proj/issues/CMT-1/comments", map[string]string{
			"contents": text,
		})
		require.Equal(t, http.StatusCreated, resp.StatusCode)
		resp.Body.Close()
	}

	var comments []model.Comment
	resp = do(t, srv, http.MethodGet, "/projects/cmt-proj/issues/CMT-1/comments", nil)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	decode(t, resp, &comments)
	assert.Len(t, comments, 2)
}
