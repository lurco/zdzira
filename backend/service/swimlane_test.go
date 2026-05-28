package service_test

import (
	"testing"
	"zdzira/backend/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSwimlaneCreate_AppendsToProject(t *testing.T) {
	svcs := newTestServices(t)
	p, err := svcs.Projects.Create(ctx, service.CreateProjectInput{Name: "Board", Shortcut: "BRD"})
	require.NoError(t, err)

	sl, err := svcs.Swimlanes.Create(ctx, service.CreateSwimlaneInput{
		ProjectSlug: "board",
		Name:        "In Review",
	})
	require.NoError(t, err)
	assert.Equal(t, "In Review", sl.Name)
	assert.Equal(t, uint(4), sl.Position, "appends after seeded Backlog/In Progress/Done")

	all, err := svcs.Swimlanes.ListByProject(ctx, p.ID)
	require.NoError(t, err)
	assert.Len(t, all, 4)
}

func TestSwimlaneRename_ChangesName(t *testing.T) {
	svcs := newTestServices(t)
	p, err := svcs.Projects.Create(ctx, service.CreateProjectInput{Name: "Board", Shortcut: "BRD"})
	require.NoError(t, err)

	swimlanes, err := svcs.Swimlanes.ListByProject(ctx, p.ID)
	require.NoError(t, err)
	backlog := swimlanes[0]

	updated, err := svcs.Swimlanes.Rename(ctx, service.RenameSwimlaneInput{
		ProjectSlug: "board",
		ID:          backlog.ID,
		Name:        "Todo",
	})
	require.NoError(t, err)
	assert.Equal(t, "Todo", updated.Name)
}
