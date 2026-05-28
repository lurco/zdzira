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

func getResultJSON(t *testing.T, result map[string]any) map[string]any {
	t.Helper()
	res := result["result"].(map[string]any)
	content := res["content"].([]any)
	raw := content[0].(map[string]any)["text"].(string)
	var out map[string]any
	require.NoError(t, json.Unmarshal([]byte(raw), &out))
	return out
}

func TestUpdateIssueTool_ChangesFields(t *testing.T) {
	s, svcs := newTestMCPServer(t)
	ctx := context.Background()

	_, err := svcs.Projects.Create(ctx, service.CreateProjectInput{Name: "Upd", Shortcut: "UPD"})
	require.NoError(t, err)
	_, err = svcs.Issues.Create(ctx, service.CreateIssueInput{
		ProjectSlug: "upd", Name: "original", Type: "TASK", Priority: "LOW",
	})
	require.NoError(t, err)

	result := callTool(t, s, "update_issue", map[string]any{
		"project":   "upd",
		"issue_ref": "UPD-1",
		"name":      "updated",
		"type":      "BUG",
		"priority":  "HIGH",
	})
	assert.False(t, isToolError(result), fmt.Sprintf("expected success, got: %v", result))

	out := getResultJSON(t, result)
	assert.Equal(t, "updated", out["name"])
	assert.Equal(t, "BUG", out["type"])
	assert.Equal(t, "HIGH", out["priority"])
}

func TestUpdateIssueTool_NonexistentIssue(t *testing.T) {
	s, svcs := newTestMCPServer(t)
	_, err := svcs.Projects.Create(context.Background(), service.CreateProjectInput{Name: "Err", Shortcut: "ERR"})
	require.NoError(t, err)

	result := callTool(t, s, "update_issue", map[string]any{
		"project":   "err",
		"issue_ref": "ERR-99",
		"name":      "ghost",
		"type":      "TASK",
		"priority":  "LOW",
	})
	assert.True(t, isToolError(result))
}

func TestListSwimlanesTool_ReturnsThreeDefault(t *testing.T) {
	s, svcs := newTestMCPServer(t)
	_, err := svcs.Projects.Create(context.Background(), service.CreateProjectInput{Name: "Swim", Shortcut: "SWM"})
	require.NoError(t, err)

	result := callTool(t, s, "list_swimlanes", map[string]any{"project": "swim"})
	assert.False(t, isToolError(result), fmt.Sprintf("expected success, got: %v", result))

	res := result["result"].(map[string]any)
	content := res["content"].([]any)
	raw := content[0].(map[string]any)["text"].(string)
	var swimlanes []map[string]any
	require.NoError(t, json.Unmarshal([]byte(raw), &swimlanes))
	assert.Len(t, swimlanes, 3)
	assert.Equal(t, "Backlog", swimlanes[0]["name"])
}

func TestUpdateEpicTool_ChangesName(t *testing.T) {
	s, svcs := newTestMCPServer(t)
	ctx := context.Background()

	_, err := svcs.Projects.Create(ctx, service.CreateProjectInput{Name: "Epic Proj", Shortcut: "EP"})
	require.NoError(t, err)
	_, err = svcs.Epics.Create(ctx, service.CreateEpicInput{ProjectSlug: "epic-proj", Name: "old name"})
	require.NoError(t, err)

	result := callTool(t, s, "update_epic", map[string]any{
		"project":  "epic-proj",
		"epic_ref": "EP-E1",
		"name":     "new name",
	})
	assert.False(t, isToolError(result), fmt.Sprintf("expected success, got: %v", result))

	out := getResultJSON(t, result)
	assert.Equal(t, "new name", out["name"])
}

func TestUpdateEpicTool_NonexistentEpic(t *testing.T) {
	s, svcs := newTestMCPServer(t)
	_, err := svcs.Projects.Create(context.Background(), service.CreateProjectInput{Name: "Err", Shortcut: "ERR"})
	require.NoError(t, err)

	result := callTool(t, s, "update_epic", map[string]any{
		"project":  "err",
		"epic_ref": "ERR-E99",
		"name":     "ghost",
	})
	assert.True(t, isToolError(result))
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
