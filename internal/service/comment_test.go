package service_test

import (
	"testing"
	"zdzira/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommentAddToIssue_SetsIssueIDOnly(t *testing.T) {
	svcs := newTestServices(t)
	_, err := svcs.Projects.Create(ctx, service.CreateProjectInput{Name: "Comment Project", Shortcut: "COM"})
	require.NoError(t, err)

	_, err = svcs.Issues.Create(ctx, service.CreateIssueInput{
		ProjectSlug: "comment-project",
		Name:        "issue",
		Type:        "TASK",
		Priority:    "LOW",
	})
	require.NoError(t, err)

	c, err := svcs.Comments.AddToIssue(ctx, "comment-project", "COM-1", "looks good")
	require.NoError(t, err)

	assert.Equal(t, "looks good", c.Contents)
	assert.NotNil(t, c.IssueID, "IssueID should be set")
	assert.Nil(t, c.EpicID, "EpicID should be nil")
	assert.Nil(t, c.ProjectID, "ProjectID should be nil")
}
