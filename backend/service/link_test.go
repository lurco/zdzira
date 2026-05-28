package service_test

import (
	"testing"
	"zdzira/backend/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLinkListForIssue_ReturnsBothDirections(t *testing.T) {
	svcs := newTestServices(t)
	_, err := svcs.Projects.Create(ctx, service.CreateProjectInput{Name: "Link Project", Shortcut: "LNK"})
	require.NoError(t, err)

	for i := 0; i < 3; i++ {
		_, err = svcs.Issues.Create(ctx, service.CreateIssueInput{
			ProjectSlug: "link-project", Name: "issue", Type: "TASK", Priority: "LOW",
		})
		require.NoError(t, err)
	}

	// LNK-1 blocks LNK-2 (LNK-1 is source)
	_, err = svcs.Links.Create(ctx, service.CreateLinkInput{
		ProjectSlug: "link-project", SourceRef: "LNK-1", TargetRef: "LNK-2", Type: "BLOCKS",
	})
	require.NoError(t, err)

	// LNK-3 blocks LNK-1 (LNK-1 is target)
	_, err = svcs.Links.Create(ctx, service.CreateLinkInput{
		ProjectSlug: "link-project", SourceRef: "LNK-3", TargetRef: "LNK-1", Type: "BLOCKS",
	})
	require.NoError(t, err)

	links, err := svcs.Links.ListForIssue(ctx, "link-project", "LNK-1")
	require.NoError(t, err)
	assert.Len(t, links, 2, "should return links where issue is source OR target")
}

func TestLinkDelete_RemovesLink(t *testing.T) {
	svcs := newTestServices(t)
	_, err := svcs.Projects.Create(ctx, service.CreateProjectInput{Name: "Link Project", Shortcut: "LNK"})
	require.NoError(t, err)

	_, err = svcs.Issues.Create(ctx, service.CreateIssueInput{
		ProjectSlug: "link-project", Name: "source", Type: "TASK", Priority: "LOW",
	})
	require.NoError(t, err)
	_, err = svcs.Issues.Create(ctx, service.CreateIssueInput{
		ProjectSlug: "link-project", Name: "target", Type: "TASK", Priority: "LOW",
	})
	require.NoError(t, err)

	l, err := svcs.Links.Create(ctx, service.CreateLinkInput{
		ProjectSlug: "link-project", SourceRef: "LNK-1", TargetRef: "LNK-2", Type: "BLOCKS",
	})
	require.NoError(t, err)

	require.NoError(t, svcs.Links.Delete(ctx, l.ID))

	remaining, err := svcs.Links.ListForIssue(ctx, "link-project", "LNK-1")
	require.NoError(t, err)
	assert.Empty(t, remaining)
}

func TestLinkCreate_DirectionPreserved(t *testing.T) {
	svcs := newTestServices(t)
	_, err := svcs.Projects.Create(ctx, service.CreateProjectInput{Name: "Link Project", Shortcut: "LNK"})
	require.NoError(t, err)

	source, err := svcs.Issues.Create(ctx, service.CreateIssueInput{
		ProjectSlug: "link-project", Name: "source", Type: "TASK", Priority: "LOW",
	})
	require.NoError(t, err)

	target, err := svcs.Issues.Create(ctx, service.CreateIssueInput{
		ProjectSlug: "link-project", Name: "target", Type: "TASK", Priority: "LOW",
	})
	require.NoError(t, err)

	l, err := svcs.Links.Create(ctx, service.CreateLinkInput{
		ProjectSlug: "link-project",
		SourceRef:   "LNK-1",
		TargetRef:   "LNK-2",
		Type:        "BLOCKS",
	})
	require.NoError(t, err)
	assert.Equal(t, source.ID, l.IssueA, "IssueA should be source")
	assert.Equal(t, target.ID, l.IssueB, "IssueB should be target")
}
