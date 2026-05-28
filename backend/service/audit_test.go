package service_test

import (
	"testing"
	"zdzira/backend/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuditIssueUpdate_RecordsUpdated(t *testing.T) {
	svcs := newTestServices(t)
	_, err := svcs.Projects.Create(ctx, service.CreateProjectInput{Name: "Audit Proj", Shortcut: "AUD"})
	require.NoError(t, err)
	_, err = svcs.Issues.Create(ctx, service.CreateIssueInput{
		ProjectSlug: "audit-proj", Name: "issue", Type: "TASK", Priority: "LOW",
	})
	require.NoError(t, err)

	_, err = svcs.Issues.Update(ctx, service.UpdateIssueInput{
		ProjectSlug: "audit-proj", IssueRef: "AUD-1", Name: "renamed", Type: "BUG", Priority: "HIGH",
	})
	require.NoError(t, err)

	entries, err := svcs.Audit.ListForProject(ctx, "audit-proj")
	require.NoError(t, err)
	assert.Equal(t, "updated", entries[0].Action)
	assert.Equal(t, "AUD-1", entries[0].Ref)
}

func TestAuditIssueMoved_RecordsMoved(t *testing.T) {
	svcs := newTestServices(t)
	_, err := svcs.Projects.Create(ctx, service.CreateProjectInput{Name: "Audit Proj", Shortcut: "AUD"})
	require.NoError(t, err)
	_, err = svcs.Issues.Create(ctx, service.CreateIssueInput{
		ProjectSlug: "audit-proj", Name: "issue", Type: "TASK", Priority: "LOW",
	})
	require.NoError(t, err)

	_, err = svcs.Issues.Move(ctx, service.MoveIssueInput{
		ProjectSlug: "audit-proj", IssueRef: "AUD-1", SwimlaneName: "Done",
	})
	require.NoError(t, err)

	entries, err := svcs.Audit.ListForProject(ctx, "audit-proj")
	require.NoError(t, err)
	assert.Equal(t, "moved", entries[0].Action)
}

func TestAuditIssueDelete_RecordsDeleted(t *testing.T) {
	svcs := newTestServices(t)
	_, err := svcs.Projects.Create(ctx, service.CreateProjectInput{Name: "Audit Proj", Shortcut: "AUD"})
	require.NoError(t, err)
	_, err = svcs.Issues.Create(ctx, service.CreateIssueInput{
		ProjectSlug: "audit-proj", Name: "issue", Type: "TASK", Priority: "LOW",
	})
	require.NoError(t, err)

	err = svcs.Issues.Delete(ctx, "audit-proj", "AUD-1")
	require.NoError(t, err)

	entries, err := svcs.Audit.ListForProject(ctx, "audit-proj")
	require.NoError(t, err)
	assert.Equal(t, "deleted", entries[0].Action)
}

func TestAuditEpicCreate_RecordsCreated(t *testing.T) {
	svcs := newTestServices(t)
	_, err := svcs.Projects.Create(ctx, service.CreateProjectInput{Name: "Audit Proj", Shortcut: "AUD"})
	require.NoError(t, err)

	_, err = svcs.Epics.Create(ctx, service.CreateEpicInput{ProjectSlug: "audit-proj", Name: "big feature"})
	require.NoError(t, err)

	entries, err := svcs.Audit.ListForProject(ctx, "audit-proj")
	require.NoError(t, err)
	require.Len(t, entries, 1)
	assert.Equal(t, "epic", entries[0].EntityType)
	assert.Equal(t, "AUD-E1", entries[0].Ref)
	assert.Equal(t, "created", entries[0].Action)
}

func TestAuditEpicDelete_RecordsDeleted(t *testing.T) {
	svcs := newTestServices(t)
	_, err := svcs.Projects.Create(ctx, service.CreateProjectInput{Name: "Audit Proj", Shortcut: "AUD"})
	require.NoError(t, err)
	_, err = svcs.Epics.Create(ctx, service.CreateEpicInput{ProjectSlug: "audit-proj", Name: "epic"})
	require.NoError(t, err)

	err = svcs.Epics.Delete(ctx, "audit-proj", "AUD-E1")
	require.NoError(t, err)

	entries, err := svcs.Audit.ListForProject(ctx, "audit-proj")
	require.NoError(t, err)
	assert.Equal(t, "deleted", entries[0].Action)
	assert.Equal(t, "AUD-E1", entries[0].Ref)
}

func TestAuditEpicUpdate_RecordsUpdated(t *testing.T) {
	svcs := newTestServices(t)
	_, err := svcs.Projects.Create(ctx, service.CreateProjectInput{Name: "Audit Proj", Shortcut: "AUD"})
	require.NoError(t, err)
	_, err = svcs.Epics.Create(ctx, service.CreateEpicInput{ProjectSlug: "audit-proj", Name: "epic"})
	require.NoError(t, err)

	_, err = svcs.Epics.Update(ctx, service.UpdateEpicInput{
		ProjectSlug: "audit-proj", EpicRef: "AUD-E1", Name: "renamed",
	})
	require.NoError(t, err)

	entries, err := svcs.Audit.ListForProject(ctx, "audit-proj")
	require.NoError(t, err)
	assert.Equal(t, "updated", entries[0].Action)
	assert.Equal(t, "AUD-E1", entries[0].Ref)
}

func TestAuditListForProject_ReturnsEntries(t *testing.T) {
	svcs := newTestServices(t)
	_, err := svcs.Projects.Create(ctx, service.CreateProjectInput{Name: "Audit Proj", Shortcut: "AUD"})
	require.NoError(t, err)

	_, err = svcs.Issues.Create(ctx, service.CreateIssueInput{
		ProjectSlug: "audit-proj", Name: "first", Type: "TASK", Priority: "LOW",
	})
	require.NoError(t, err)

	entries, err := svcs.Audit.ListForProject(ctx, "audit-proj")
	require.NoError(t, err)
	require.Len(t, entries, 1)
	assert.Equal(t, "issue", entries[0].EntityType)
	assert.Equal(t, "AUD-1", entries[0].Ref)
	assert.Equal(t, "created", entries[0].Action)
}
