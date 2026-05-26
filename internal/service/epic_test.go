package service_test

import (
	"testing"
	"zdzira/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEpicCreate_SequentialNumbers(t *testing.T) {
	svcs := newTestServices(t)
	_, err := svcs.Projects.Create(ctx, service.CreateProjectInput{Name: "Epic Project", Shortcut: "EP"})
	require.NoError(t, err)

	for i := uint(1); i <= 3; i++ {
		e, err := svcs.Epics.Create(ctx, service.CreateEpicInput{
			ProjectSlug: "epic-project",
			Name:        "epic",
		})
		require.NoError(t, err)
		assert.Equal(t, i, e.Number)
	}
}

func TestEpicList_ReturnsEpicsForProject(t *testing.T) {
	svcs := newTestServices(t)
	_, err := svcs.Projects.Create(ctx, service.CreateProjectInput{Name: "Epic Project", Shortcut: "EP"})
	require.NoError(t, err)

	for i := 0; i < 3; i++ {
		_, err = svcs.Epics.Create(ctx, service.CreateEpicInput{ProjectSlug: "epic-project", Name: "epic"})
		require.NoError(t, err)
	}

	epics, err := svcs.Epics.List(ctx, "epic-project")
	require.NoError(t, err)
	assert.Len(t, epics, 3)
}

func TestEpicDelete_RemovesEpic(t *testing.T) {
	svcs := newTestServices(t)
	_, err := svcs.Projects.Create(ctx, service.CreateProjectInput{Name: "Epic Project", Shortcut: "EP"})
	require.NoError(t, err)

	_, err = svcs.Epics.Create(ctx, service.CreateEpicInput{ProjectSlug: "epic-project", Name: "to delete"})
	require.NoError(t, err)

	require.NoError(t, svcs.Epics.Delete(ctx, "epic-project", "EP-E1"))

	_, err = svcs.Epics.Get(ctx, "epic-project", "EP-E1")
	assert.Error(t, err, "deleted epic should not be found")
}

func TestEpicGet_ByRef(t *testing.T) {
	svcs := newTestServices(t)
	_, err := svcs.Projects.Create(ctx, service.CreateProjectInput{Name: "Epic Project", Shortcut: "EP"})
	require.NoError(t, err)

	created, err := svcs.Epics.Create(ctx, service.CreateEpicInput{ProjectSlug: "epic-project", Name: "the epic"})
	require.NoError(t, err)

	got, err := svcs.Epics.Get(ctx, "epic-project", "EP-E1")
	require.NoError(t, err)
	assert.Equal(t, created.ID, got.ID)
}
