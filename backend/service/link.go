package service

import (
	"context"
	"zdzira/backend/model"
	"zdzira/backend/store"
)

type LinkService struct {
	stores *store.Stores
}

type CreateLinkInput struct {
	ProjectSlug string
	SourceRef   string
	TargetRef   string
	Type        model.LinkType
}

func (s *LinkService) Create(ctx context.Context, in CreateLinkInput) (*model.Link, error) {
	p, err := s.stores.Projects.GetBySlug(ctx, in.ProjectSlug)
	if err != nil {
		return nil, err
	}

	sourceNum, err := parseIssueRef(p.Shortcut, in.SourceRef)
	if err != nil {
		return nil, err
	}
	targetNum, err := parseIssueRef(p.Shortcut, in.TargetRef)
	if err != nil {
		return nil, err
	}

	source, err := s.stores.Issues.GetByRef(ctx, p.ID, sourceNum)
	if err != nil {
		return nil, err
	}
	target, err := s.stores.Issues.GetByRef(ctx, p.ID, targetNum)
	if err != nil {
		return nil, err
	}

	l := &model.Link{
		Type:   in.Type,
		IssueA: source.ID,
		IssueB: target.ID,
	}
	if err := s.stores.Links.Create(ctx, l); err != nil {
		return nil, err
	}
	return l, nil
}

func (s *LinkService) ListForIssue(ctx context.Context, projectSlug, issueRef string) ([]model.Link, error) {
	p, err := s.stores.Projects.GetBySlug(ctx, projectSlug)
	if err != nil {
		return nil, err
	}
	num, err := parseIssueRef(p.Shortcut, issueRef)
	if err != nil {
		return nil, err
	}
	issue, err := s.stores.Issues.GetByRef(ctx, p.ID, num)
	if err != nil {
		return nil, err
	}
	return s.stores.Links.ListByIssue(ctx, issue.ID)
}

func (s *LinkService) Delete(ctx context.Context, id uint) error {
	return s.stores.Links.Delete(ctx, id)
}
