package service_test

import (
	"testing"
	"zdzira/backend/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBoardGet_GroupsIssuesByLane(t *testing.T) {
	svcs := newTestServices(t)
	_, err := svcs.Projects.Create(ctx, service.CreateProjectInput{Name: "Board", Shortcut: "BRD"})
	require.NoError(t, err)

	swimlanes, err := svcs.Swimlanes.ListForProject(ctx, "board")
	require.NoError(t, err)
	require.Len(t, swimlanes, 3)
	backlog, inProgress := swimlanes[0], swimlanes[1]

	_, err = svcs.Issues.Create(ctx, service.CreateIssueInput{
		ProjectSlug: "board", Name: "First", SwimlaneID: &backlog.ID,
	})
	require.NoError(t, err)
	_, err = svcs.Issues.Create(ctx, service.CreateIssueInput{
		ProjectSlug: "board", Name: "Second", SwimlaneID: &backlog.ID,
	})
	require.NoError(t, err)
	_, err = svcs.Issues.Create(ctx, service.CreateIssueInput{
		ProjectSlug: "board", Name: "Third", SwimlaneID: &inProgress.ID,
	})
	require.NoError(t, err)

	view, err := svcs.Board.Get(ctx, "board")
	require.NoError(t, err)
	require.Len(t, view.Swimlanes, 3)

	assert.Equal(t, backlog.ID, view.Swimlanes[0].ID)
	assert.Len(t, view.Swimlanes[0].Issues, 2)
	assert.Equal(t, "BRD-1", view.Swimlanes[0].Issues[0].Ref)
	assert.Equal(t, "BRD-2", view.Swimlanes[0].Issues[1].Ref)

	assert.Equal(t, inProgress.ID, view.Swimlanes[1].ID)
	assert.Len(t, view.Swimlanes[1].Issues, 1)
	assert.Equal(t, "BRD-3", view.Swimlanes[1].Issues[0].Ref)

	assert.Empty(t, view.Swimlanes[2].Issues)
}
