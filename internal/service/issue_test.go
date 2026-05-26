package service_test

import (
	"testing"
	"zdzira/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newProjectWithIssues(t *testing.T, svcs *service.Services) {
	t.Helper()
	_, err := svcs.Projects.Create(ctx, service.CreateProjectInput{Name: "Tracker", Shortcut: "TRK"})
	require.NoError(t, err)
}

func TestIssueList_ReturnsIssuesForProject(t *testing.T) {
	svcs := newTestServices(t)
	newProjectWithIssues(t, svcs)

	for i := 0; i < 3; i++ {
		_, err := svcs.Issues.Create(ctx, service.CreateIssueInput{
			ProjectSlug: "tracker",
			Name:        "issue",
			Type:        "TASK",
			Priority:    "LOW",
		})
		require.NoError(t, err)
	}

	issues, err := svcs.Issues.List(ctx, "tracker")
	require.NoError(t, err)
	assert.Len(t, issues, 3)
}

func TestIssueDelete_RemovesIssue(t *testing.T) {
	svcs := newTestServices(t)
	newProjectWithIssues(t, svcs)

	_, err := svcs.Issues.Create(ctx, service.CreateIssueInput{
		ProjectSlug: "tracker",
		Name:        "to delete",
		Type:        "TASK",
		Priority:    "LOW",
	})
	require.NoError(t, err)

	require.NoError(t, svcs.Issues.Delete(ctx, "tracker", "TRK-1"))

	_, err = svcs.Issues.Get(ctx, "tracker", "TRK-1")
	assert.Error(t, err, "deleted issue should not be found")
}

func TestIssueCreate_SequentialNumbers(t *testing.T) {
	svcs := newTestServices(t)
	newProjectWithIssues(t, svcs)

	for i := uint(1); i <= 3; i++ {
		issue, err := svcs.Issues.Create(ctx, service.CreateIssueInput{
			ProjectSlug: "tracker",
			Name:        "issue",
			Type:        "TASK",
			Priority:    "LOW",
		})
		require.NoError(t, err)
		assert.Equal(t, i, issue.Number)
	}
}

func TestIssueCreate_DefaultsToFirstSwimlane(t *testing.T) {
	svcs := newTestServices(t)
	p, err := svcs.Projects.Create(ctx, service.CreateProjectInput{Name: "Swim", Shortcut: "SWM"})
	require.NoError(t, err)

	issue, err := svcs.Issues.Create(ctx, service.CreateIssueInput{
		ProjectSlug: "swim",
		Name:        "no swimlane given",
		Type:        "BUG",
		Priority:    "HIGH",
	})
	require.NoError(t, err)

	swimlanes, err := svcs.Swimlanes.ListByProject(ctx, p.ID)
	require.NoError(t, err)
	assert.Equal(t, swimlanes[0].ID, issue.SwimlaneID, "should land in Backlog (first swimlane)")
}

func TestIssueCreate_AttachesToEpic(t *testing.T) {
	svcs := newTestServices(t)
	newProjectWithIssues(t, svcs)

	epic, err := svcs.Epics.Create(ctx, service.CreateEpicInput{ProjectSlug: "tracker", Name: "big feature"})
	require.NoError(t, err)

	epicRef := "TRK-E1"
	issue, err := svcs.Issues.Create(ctx, service.CreateIssueInput{
		ProjectSlug: "tracker",
		Name:        "part of epic",
		Type:        "STORY",
		Priority:    "HIGH",
		EpicRef:     &epicRef,
	})
	require.NoError(t, err)
	require.NotNil(t, issue.EpicID)
	assert.Equal(t, epic.ID, *issue.EpicID)
}

func TestIssueGet_ByRef(t *testing.T) {
	svcs := newTestServices(t)
	newProjectWithIssues(t, svcs)

	created, err := svcs.Issues.Create(ctx, service.CreateIssueInput{
		ProjectSlug: "tracker",
		Name:        "find me",
		Type:        "TASK",
		Priority:    "LOW",
	})
	require.NoError(t, err)

	got, err := svcs.Issues.Get(ctx, "tracker", "TRK-1")
	require.NoError(t, err)
	assert.Equal(t, created.ID, got.ID)
}

func TestIssueGet_RefCaseInsensitive(t *testing.T) {
	svcs := newTestServices(t)
	newProjectWithIssues(t, svcs)

	_, err := svcs.Issues.Create(ctx, service.CreateIssueInput{
		ProjectSlug: "tracker",
		Name:        "case test",
		Type:        "TASK",
		Priority:    "LOW",
	})
	require.NoError(t, err)

	_, err = svcs.Issues.Get(ctx, "tracker", "trk-1")
	assert.NoError(t, err, "lowercase ref should resolve same as uppercase")
}

func TestIssueGet_InvalidRef(t *testing.T) {
	svcs := newTestServices(t)
	newProjectWithIssues(t, svcs)

	_, err := svcs.Issues.Get(ctx, "tracker", "TRK-abc")
	assert.Error(t, err, "non-numeric number part should error")

	_, err = svcs.Issues.Get(ctx, "tracker", "OTHER-1")
	assert.Error(t, err, "wrong shortcut should error")
}

func TestIssueMove_ChangesSwimlane(t *testing.T) {
	svcs := newTestServices(t)
	newProjectWithIssues(t, svcs)

	_, err := svcs.Issues.Create(ctx, service.CreateIssueInput{
		ProjectSlug: "tracker",
		Name:        "move me",
		Type:        "TASK",
		Priority:    "LOW",
	})
	require.NoError(t, err)

	moved, err := svcs.Issues.Move(ctx, service.MoveIssueInput{
		ProjectSlug:  "tracker",
		IssueRef:     "TRK-1",
		SwimlaneName: "Done",
	})
	require.NoError(t, err)

	p, err := svcs.Projects.Get(ctx, "tracker")
	require.NoError(t, err)
	swimlanes, err := svcs.Swimlanes.ListByProject(ctx, p.ID)
	require.NoError(t, err)

	doneSwimlane := swimlanes[2] // Backlog(0), In Progress(1), Done(2)
	assert.Equal(t, doneSwimlane.ID, moved.SwimlaneID)
}

func TestIssueMove_InvalidSwimlane(t *testing.T) {
	svcs := newTestServices(t)
	newProjectWithIssues(t, svcs)

	_, err := svcs.Issues.Create(ctx, service.CreateIssueInput{
		ProjectSlug: "tracker",
		Name:        "move me",
		Type:        "TASK",
		Priority:    "LOW",
	})
	require.NoError(t, err)

	_, err = svcs.Issues.Move(ctx, service.MoveIssueInput{
		ProjectSlug:  "tracker",
		IssueRef:     "TRK-1",
		SwimlaneName: "Nonexistent",
	})
	assert.Error(t, err)
}
