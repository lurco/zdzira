package service

import (
	"context"
	"fmt"
	"zdzira/backend/model"
	"zdzira/backend/store"
)

type CommentService struct {
	stores *store.Stores
}

func (s *CommentService) AddToIssue(ctx context.Context, projectSlug, issueRef, contents string) (*model.Comment, error) {
	issue, err := s.resolveIssue(ctx, projectSlug, issueRef)
	if err != nil {
		return nil, err
	}
	c := &model.Comment{Contents: contents, IssueID: &issue.ID}
	return c, s.stores.Comments.Create(ctx, c)
}

func (s *CommentService) AddToEpic(ctx context.Context, projectSlug, epicRef, contents string) (*model.Comment, error) {
	p, err := s.stores.Projects.GetBySlug(ctx, projectSlug)
	if err != nil {
		return nil, fmt.Errorf("project %q not found", projectSlug)
	}
	number, err := parseEpicRef(p.Shortcut, epicRef)
	if err != nil {
		return nil, err
	}
	epic, err := s.stores.Epics.GetByRef(ctx, p.ID, number)
	if err != nil {
		return nil, err
	}
	c := &model.Comment{Contents: contents, EpicID: &epic.ID}
	return c, s.stores.Comments.Create(ctx, c)
}

func (s *CommentService) ListForEpic(ctx context.Context, projectSlug, epicRef string) ([]model.Comment, error) {
	p, err := s.stores.Projects.GetBySlug(ctx, projectSlug)
	if err != nil {
		return nil, fmt.Errorf("project %q not found", projectSlug)
	}
	number, err := parseEpicRef(p.Shortcut, epicRef)
	if err != nil {
		return nil, err
	}
	epic, err := s.stores.Epics.GetByRef(ctx, p.ID, number)
	if err != nil {
		return nil, err
	}
	return s.stores.Comments.ListByEpic(ctx, epic.ID)
}

func (s *CommentService) AddToProject(ctx context.Context, projectSlug, contents string) (*model.Comment, error) {
	p, err := s.stores.Projects.GetBySlug(ctx, projectSlug)
	if err != nil {
		return nil, fmt.Errorf("project %q not found", projectSlug)
	}
	c := &model.Comment{Contents: contents, ProjectID: &p.ID}
	return c, s.stores.Comments.Create(ctx, c)
}

func (s *CommentService) ListForIssue(ctx context.Context, projectSlug, issueRef string) ([]model.Comment, error) {
	issue, err := s.resolveIssue(ctx, projectSlug, issueRef)
	if err != nil {
		return nil, err
	}
	return s.stores.Comments.ListByIssue(ctx, issue.ID)
}

func (s *CommentService) UpdateContents(ctx context.Context, id uint, contents string) (*model.Comment, error) {
	c, err := s.stores.Comments.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	c.Contents = contents
	if err := s.stores.Comments.Update(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

// CountsByIssueIDs returns issue ID → comment count for a batch of issues, so
// list and board views can surface a comment cue without an N+1 of list calls.
func (s *CommentService) CountsByIssueIDs(ctx context.Context, issueIDs []uint) (map[uint]uint, error) {
	return s.stores.Comments.CountByIssueIDs(ctx, issueIDs)
}

// CountsByEpicIDs is the epic-keyed counterpart of CountsByIssueIDs.
func (s *CommentService) CountsByEpicIDs(ctx context.Context, epicIDs []uint) (map[uint]uint, error) {
	return s.stores.Comments.CountByEpicIDs(ctx, epicIDs)
}

func (s *CommentService) Delete(ctx context.Context, id uint) error {
	return s.stores.Comments.Delete(ctx, id)
}

func (s *CommentService) resolveIssue(ctx context.Context, projectSlug, ref string) (*model.Issue, error) {
	p, err := s.stores.Projects.GetBySlug(ctx, projectSlug)
	if err != nil {
		return nil, fmt.Errorf("project %q not found", projectSlug)
	}
	number, err := parseIssueRef(p.Shortcut, ref)
	if err != nil {
		return nil, err
	}
	return s.stores.Issues.GetByRef(ctx, p.ID, number)
}
