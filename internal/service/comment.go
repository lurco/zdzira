package service

import (
	"context"
	"fmt"
	"zdzira/internal/model"
	"zdzira/internal/store"
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
