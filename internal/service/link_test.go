package service_test

import (
	"testing"
	"zdzira/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
