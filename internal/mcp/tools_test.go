package mcp_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	internalmcp "zdzira/internal/mcp"
	"zdzira/internal/service"
	"zdzira/internal/store"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestMCPServer(t *testing.T) (*internalmcp.Server, *service.Services) {
	t.Helper()
	db, err := store.Open(":memory:")
	require.NoError(t, err)
	svcs := service.New(store.New(db))
	return internalmcp.NewServer(svcs), svcs
}

func callTool(t *testing.T, s *internalmcp.Server, name string, params map[string]any) map[string]any {
	t.Helper()
	msg, err := json.Marshal(map[string]any{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "tools/call",
		"params":  map[string]any{"name": name, "arguments": params},
	})
	require.NoError(t, err)

	raw := s.HandleMessage(context.Background(), msg)
	b, err := json.Marshal(raw)
	require.NoError(t, err)

	var result map[string]any
	require.NoError(t, json.Unmarshal(b, &result))
	return result
}

func isToolError(result map[string]any) bool {
	res, ok := result["result"].(map[string]any)
	if !ok {
		return true
	}
	isError, _ := res["isError"].(bool)
	return isError
}

func TestDeleteIssueTool_RemovesIssue(t *testing.T) {
	s, svcs := newTestMCPServer(t)

	_, err := svcs.Projects.Create(context.Background(), service.CreateProjectInput{
		Name: "MCP Project", Shortcut: "MCP",
	})
	require.NoError(t, err)

	_, err = svcs.Issues.Create(context.Background(), service.CreateIssueInput{
		ProjectSlug: "mcp-project", Name: "to delete", Type: "TASK", Priority: "LOW",
	})
	require.NoError(t, err)

	result := callTool(t, s, "delete_issue", map[string]any{
		"project":   "mcp-project",
		"issue_ref": "MCP-1",
	})
	assert.False(t, isToolError(result), fmt.Sprintf("expected success, got: %v", result))

	_, err = svcs.Issues.Get(context.Background(), "mcp-project", "MCP-1")
	assert.Error(t, err, "issue should be soft-deleted")
}
