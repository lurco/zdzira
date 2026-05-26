package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"zdzira/internal/model"
	"zdzira/internal/store"
)

type IssueService struct {
	stores *store.Stores
}

type CreateIssueInput struct {
	ProjectSlug string
	Type        model.IssueType
	Priority    model.Priority
	Name        string
	Description *string
	EpicRef     *string
	SwimlaneID  *uint
}

type MoveIssueInput struct {
	ProjectSlug  string
	IssueRef     string
	SwimlaneName string
}

func (s *IssueService) Create(ctx context.Context, in CreateIssueInput) (*model.Issue, error) {
	p, err := s.stores.Projects.GetBySlug(ctx, in.ProjectSlug)
	if err != nil {
		return nil, fmt.Errorf("project %q not found", in.ProjectSlug)
	}

	swimlaneID, err := s.resolveSwimlane(ctx, p.ID, in.SwimlaneID)
	if err != nil {
		return nil, err
	}

	p.IssueCounter++
	if err := s.stores.Projects.Update(ctx, p); err != nil {
		return nil, err
	}

	issue := &model.Issue{
		Number:      p.IssueCounter,
		Type:        in.Type,
		Priority:    in.Priority,
		Name:        in.Name,
		Description: in.Description,
		ProjectID:   p.ID,
		SwimlaneID:  swimlaneID,
		Position:    p.IssueCounter,
	}

	if in.EpicRef != nil {
		epic, err := s.resolveEpicRef(ctx, p, *in.EpicRef)
		if err != nil {
			return nil, err
		}
		issue.EpicID = &epic.ID
	}

	if err := s.stores.Issues.Create(ctx, issue); err != nil {
		return nil, err
	}
	return issue, nil
}

func (s *IssueService) Get(ctx context.Context, projectSlug, ref string) (*model.Issue, error) {
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

func (s *IssueService) List(ctx context.Context, projectSlug string) ([]model.Issue, error) {
	p, err := s.stores.Projects.GetBySlug(ctx, projectSlug)
	if err != nil {
		return nil, fmt.Errorf("project %q not found", projectSlug)
	}
	return s.stores.Issues.ListByProject(ctx, p.ID)
}

func (s *IssueService) Move(ctx context.Context, in MoveIssueInput) (*model.Issue, error) {
	issue, err := s.Get(ctx, in.ProjectSlug, in.IssueRef)
	if err != nil {
		return nil, err
	}

	p, err := s.stores.Projects.GetBySlug(ctx, in.ProjectSlug)
	if err != nil {
		return nil, err
	}

	sl, err := s.stores.Swimlanes.GetByName(ctx, p.ID, in.SwimlaneName)
	if err != nil {
		return nil, fmt.Errorf("swimlane %q not found in project", in.SwimlaneName)
	}

	issue.SwimlaneID = sl.ID
	return issue, s.stores.Issues.Update(ctx, issue)
}

func (s *IssueService) Delete(ctx context.Context, projectSlug, ref string) error {
	issue, err := s.Get(ctx, projectSlug, ref)
	if err != nil {
		return err
	}
	return s.stores.Issues.Delete(ctx, issue.ID)
}

func (s *IssueService) resolveSwimlane(ctx context.Context, projectID uint, id *uint) (uint, error) {
	if id != nil {
		return *id, nil
	}
	swimlanes, err := s.stores.Swimlanes.ListByProject(ctx, projectID)
	if err != nil || len(swimlanes) == 0 {
		return 0, fmt.Errorf("no swimlanes found for project")
	}
	return swimlanes[0].ID, nil
}

func (s *IssueService) resolveEpicRef(ctx context.Context, p *model.Project, ref string) (*model.Epic, error) {
	number, err := parseEpicRef(p.Shortcut, ref)
	if err != nil {
		return nil, err
	}
	return s.stores.Epics.GetByRef(ctx, p.ID, number)
}

func parseIssueRef(shortcut, ref string) (uint, error) {
	prefix := strings.ToUpper(shortcut) + "-"
	if !strings.HasPrefix(strings.ToUpper(ref), prefix) {
		return 0, fmt.Errorf("invalid issue ref %q: expected format %s<number>", ref, prefix)
	}
	n, err := strconv.ParseUint(strings.TrimPrefix(strings.ToUpper(ref), prefix), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid issue ref %q: number part is not an integer", ref)
	}
	return uint(n), nil
}

func parseEpicRef(shortcut, ref string) (uint, error) {
	prefix := strings.ToUpper(shortcut) + "-E"
	if !strings.HasPrefix(strings.ToUpper(ref), prefix) {
		return 0, fmt.Errorf("invalid epic ref %q: expected format %sE<number>", ref, strings.ToUpper(shortcut)+"-")
	}
	n, err := strconv.ParseUint(strings.TrimPrefix(strings.ToUpper(ref), prefix), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid epic ref %q: number part is not an integer", ref)
	}
	return uint(n), nil
}
