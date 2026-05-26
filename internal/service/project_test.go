package service_test

import (
	"testing"
	"zdzira/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProjectCreate_SlugDerivedFromName(t *testing.T) {
	svcs := newTestServices(t)

	p, err := svcs.Projects.Create(ctx, service.CreateProjectInput{
		Name:     "My Project 123!",
		Shortcut: "PROJ",
	})

	require.NoError(t, err)
	assert.Equal(t, "my-project-123", p.Slug)
}

func TestProjectCreate_ShortcutUppercased(t *testing.T) {
	svcs := newTestServices(t)

	p, err := svcs.Projects.Create(ctx, service.CreateProjectInput{
		Name:     "Alpha",
		Shortcut: "alp",
	})

	require.NoError(t, err)
	assert.Equal(t, "ALP", p.Shortcut)
}

func TestProjectDelete_CascadesToChildren(t *testing.T) {
	svcs := newTestServices(t)

	p, err := svcs.Projects.Create(ctx, service.CreateProjectInput{Name: "Gamma", Shortcut: "GAM"})
	require.NoError(t, err)

	_, err = svcs.Issues.Create(ctx, service.CreateIssueInput{
		ProjectSlug: "gamma",
		Name:        "issue one",
		Type:        "TASK",
		Priority:    "LOW",
	})
	require.NoError(t, err)

	_, err = svcs.Epics.Create(ctx, service.CreateEpicInput{ProjectSlug: "gamma", Name: "epic one"})
	require.NoError(t, err)

	require.NoError(t, svcs.Projects.Delete(ctx, "gamma"))

	issues, err := svcs.Issues.List(ctx, "gamma")
	assert.Error(t, err, "listing issues on deleted project should fail")
	assert.Empty(t, issues)

	swimlanes, err := svcs.Swimlanes.ListByProject(ctx, p.ID)
	require.NoError(t, err)
	assert.Empty(t, swimlanes, "swimlanes should be soft-deleted")
}

func TestProjectCreate_SeedsSwimlanes(t *testing.T) {
	svcs := newTestServices(t)

	p, err := svcs.Projects.Create(ctx, service.CreateProjectInput{
		Name:     "Beta",
		Shortcut: "BET",
	})
	require.NoError(t, err)

	swimlanes, err := svcs.Swimlanes.ListByProject(ctx, p.ID)
	require.NoError(t, err)

	require.Len(t, swimlanes, 3)
	assert.Equal(t, "Backlog", swimlanes[0].Name)
	assert.Equal(t, uint(1), swimlanes[0].Position)
	assert.Equal(t, "In Progress", swimlanes[1].Name)
	assert.Equal(t, uint(2), swimlanes[1].Position)
	assert.Equal(t, "Done", swimlanes[2].Name)
	assert.Equal(t, uint(3), swimlanes[2].Position)
}
