package service

import (
	"context"
	"fmt"
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

// LinkedIssue is a lightweight issue reference used inside EnrichedLink.
type LinkedIssue struct {
	Ref  string `json:"ref"`
	Name string `json:"name"`
}

// EnrichedLink replaces the raw model.Link in API/MCP responses so callers
// see human-readable refs instead of internal database IDs.
type EnrichedLink struct {
	ID     uint           `json:"id"`
	Type   model.LinkType `json:"type"`
	Source LinkedIssue    `json:"source"`
	Target LinkedIssue    `json:"target"`
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

func (s *LinkService) ListForIssue(ctx context.Context, projectSlug, issueRef string) ([]EnrichedLink, error) {
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
	links, err := s.stores.Links.ListByIssue(ctx, issue.ID)
	if err != nil {
		return nil, err
	}

	allIssues, err := s.stores.Issues.ListByProject(ctx, p.ID)
	if err != nil {
		return nil, err
	}
	byID := make(map[uint]model.Issue, len(allIssues))
	for _, i := range allIssues {
		byID[i.ID] = i
	}

	enriched := make([]EnrichedLink, len(links))
	for i, l := range links {
		src := byID[l.IssueA]
		tgt := byID[l.IssueB]
		enriched[i] = EnrichedLink{
			ID:   l.ID,
			Type: l.Type,
			Source: LinkedIssue{
				Ref:  fmt.Sprintf("%s-%d", p.Shortcut, src.Number),
				Name: src.Name,
			},
			Target: LinkedIssue{
				Ref:  fmt.Sprintf("%s-%d", p.Shortcut, tgt.Number),
				Name: tgt.Name,
			},
		}
	}
	return enriched, nil
}

func (s *LinkService) Delete(ctx context.Context, id uint) error {
	return s.stores.Links.Delete(ctx, id)
}
