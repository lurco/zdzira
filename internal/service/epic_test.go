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
