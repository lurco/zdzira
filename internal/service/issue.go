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
	audit  *AuditService
}

type CreateIssueInput struct {
	ProjectSlug string          `json:"-"`
	Type        model.IssueType `json:"type"        doc:"Issue type"                              example:"BUG"`
	Priority    model.Priority  `json:"priority"    doc:"Priority level"                          example:"HIGH"`
	Name        string          `json:"name"        doc:"Short issue title"                        example:"Fix login timeout"`
	Description *string         `json:"description,omitempty" doc:"Optional description"`
	EpicRef     *string         `json:"epic_ref,omitempty"    doc:"Epic reference to attach to, e.g. PROJ-E1" example:"PROJ-E1"`
	SwimlaneID  *uint           `json:"swimlane_id,omitempty" doc:"Override starting swimlane ID"`
}

type MoveIssueInput struct {
	ProjectSlug  string `json:"-"`
	IssueRef     string `json:"-"`
	SwimlaneName string `json:"swimlane" doc:"Target swimlane name" example:"In Progress"`
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
	s.audit.record(ctx, p.ID, "issue", fmt.Sprintf("%s-%d", p.Shortcut, issue.Number), "created")
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
	if err := s.stores.Issues.Update(ctx, issue); err != nil {
		return nil, err
	}
	s.audit.record(ctx, p.ID, "issue", fmt.Sprintf("%s-%d", p.Shortcut, issue.Number), "moved")
	return issue, nil
}

type UpdateIssueInput struct {
	ProjectSlug string          `json:"-"`
	IssueRef    string          `json:"-"`
	Name        string          `json:"name"        doc:"Updated issue title"    example:"Fix login timeout (revised)"`
	Description *string         `json:"description,omitempty" doc:"Updated description"`
	Type        model.IssueType `json:"type"        doc:"Updated issue type"     example:"BUG"`
	Priority    model.Priority  `json:"priority"    doc:"Updated priority level" example:"IMMEDIATE"`
}

func (s *IssueService) Update(ctx context.Context, in UpdateIssueInput) (*model.Issue, error) {
	p, err := s.stores.Projects.GetBySlug(ctx, in.ProjectSlug)
	if err != nil {
		return nil, fmt.Errorf("project %q not found", in.ProjectSlug)
	}
	issue, err := s.Get(ctx, in.ProjectSlug, in.IssueRef)
	if err != nil {
		return nil, err
	}
	issue.Name = in.Name
	issue.Description = in.Description
	issue.Type = in.Type
	issue.Priority = in.Priority
	if err := s.stores.Issues.Update(ctx, issue); err != nil {
		return nil, err
	}
	s.audit.record(ctx, p.ID, "issue", fmt.Sprintf("%s-%d", p.Shortcut, issue.Number), "updated")
	return issue, nil
}

type IssueFilterInput struct {
	ProjectSlug string           `json:"-"`
	Type        *model.IssueType `json:"type,omitempty"        doc:"Filter by issue type"     example:"BUG"`
	Priority    *model.Priority  `json:"priority,omitempty"    doc:"Filter by priority"       example:"HIGH"`
	SwimlaneID  *uint            `json:"swimlane_id,omitempty" doc:"Filter by swimlane ID"`
	EpicID      *uint            `json:"epic_id,omitempty"     doc:"Filter by epic ID"`
}

func (s *IssueService) Filter(ctx context.Context, in IssueFilterInput) ([]model.Issue, error) {
	p, err := s.stores.Projects.GetBySlug(ctx, in.ProjectSlug)
	if err != nil {
		return nil, fmt.Errorf("project %q not found", in.ProjectSlug)
	}
	return s.stores.Issues.ListFiltered(ctx, p.ID, store.IssueStoreFilter{
		Type:       in.Type,
		Priority:   in.Priority,
		SwimlaneID: in.SwimlaneID,
		EpicID:     in.EpicID,
	})
}

func (s *IssueService) Delete(ctx context.Context, projectSlug, ref string) error {
	p, err := s.stores.Projects.GetBySlug(ctx, projectSlug)
	if err != nil {
		return fmt.Errorf("project %q not found", projectSlug)
	}
	issue, err := s.Get(ctx, projectSlug, ref)
	if err != nil {
		return err
	}
	if err := s.stores.Issues.Delete(ctx, issue.ID); err != nil {
		return err
	}
	s.audit.record(ctx, p.ID, "issue", fmt.Sprintf("%s-%d", p.Shortcut, issue.Number), "deleted")
	return nil
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
