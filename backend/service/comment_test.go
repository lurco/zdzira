package service_test

import (
	"testing"
	"zdzira/backend/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommentAddToEpic_SetsEpicIDOnly(t *testing.T) {
	svcs := newTestServices(t)
	_, err := svcs.Projects.Create(ctx, service.CreateProjectInput{Name: "Cmt Proj", Shortcut: "CMT"})
	require.NoError(t, err)
	_, err = svcs.Epics.Create(ctx, service.CreateEpicInput{ProjectSlug: "cmt-proj", Name: "epic"})
	require.NoError(t, err)

	c, err := svcs.Comments.AddToEpic(ctx, "cmt-proj", "CMT-E1", "epic comment")
	require.NoError(t, err)

	assert.NotNil(t, c.EpicID)
	assert.Nil(t, c.IssueID)
	assert.Nil(t, c.ProjectID)
}

func TestCommentAddToProject_SetsProjectIDOnly(t *testing.T) {
	svcs := newTestServices(t)
	_, err := svcs.Projects.Create(ctx, service.CreateProjectInput{Name: "Cmt Proj", Shortcut: "CMT"})
	require.NoError(t, err)

	c, err := svcs.Comments.AddToProject(ctx, "cmt-proj", "project comment")
	require.NoError(t, err)

	assert.NotNil(t, c.ProjectID)
	assert.Nil(t, c.IssueID)
	assert.Nil(t, c.EpicID)
}

func TestCommentListForIssue_ReturnsComments(t *testing.T) {
	svcs := newTestServices(t)
	_, err := svcs.Projects.Create(ctx, service.CreateProjectInput{Name: "Cmt Proj", Shortcut: "CMT"})
	require.NoError(t, err)
	_, err = svcs.Issues.Create(ctx, service.CreateIssueInput{
		ProjectSlug: "cmt-proj",
		Name:        "issue",
		Type:        "TASK",
		Priority:    "LOW",
	})
	require.NoError(t, err)

	_, err = svcs.Comments.AddToIssue(ctx, "cmt-proj", "CMT-1", "first")
	require.NoError(t, err)
	_, err = svcs.Comments.AddToIssue(ctx, "cmt-proj", "CMT-1", "second")
	require.NoError(t, err)

	comments, err := svcs.Comments.ListForIssue(ctx, "cmt-proj", "CMT-1")
	require.NoError(t, err)
	assert.Len(t, comments, 2)
}

func TestCommentDelete_RemovesComment(t *testing.T) {
	svcs := newTestServices(t)
	_, err := svcs.Projects.Create(ctx, service.CreateProjectInput{Name: "Cmt Proj", Shortcut: "CMT"})
	require.NoError(t, err)
	_, err = svcs.Issues.Create(ctx, service.CreateIssueInput{
		ProjectSlug: "cmt-proj",
		Name:        "issue",
		Type:        "TASK",
		Priority:    "LOW",
	})
	require.NoError(t, err)

	c, err := svcs.Comments.AddToIssue(ctx, "cmt-proj", "CMT-1", "delete me")
	require.NoError(t, err)

	require.NoError(t, svcs.Comments.Delete(ctx, c.ID))

	remaining, err := svcs.Comments.ListForIssue(ctx, "cmt-proj", "CMT-1")
	require.NoError(t, err)
	assert.Empty(t, remaining)
}

func TestCommentUpdateContents_ChangesTextKeepsParent(t *testing.T) {
	svcs := newTestServices(t)
	_, err := svcs.Projects.Create(ctx, service.CreateProjectInput{Name: "Cmt Proj", Shortcut: "CMT"})
	require.NoError(t, err)
	_, err = svcs.Issues.Create(ctx, service.CreateIssueInput{
		ProjectSlug: "cmt-proj",
		Name:        "issue",
		Type:        "TASK",
		Priority:    "LOW",
	})
	require.NoError(t, err)

	c, err := svcs.Comments.AddToIssue(ctx, "cmt-proj", "CMT-1", "typo here")
	require.NoError(t, err)

	updated, err := svcs.Comments.UpdateContents(ctx, c.ID, "fixed text")
	require.NoError(t, err)
	assert.Equal(t, "fixed text", updated.Contents)
	assert.Equal(t, c.IssueID, updated.IssueID)

	comments, err := svcs.Comments.ListForIssue(ctx, "cmt-proj", "CMT-1")
	require.NoError(t, err)
	require.Len(t, comments, 1)
	assert.Equal(t, "fixed text", comments[0].Contents)
}

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
