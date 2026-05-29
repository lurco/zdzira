package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"zdzira/backend/model"
	"zdzira/backend/store"
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
	SwimlaneName string `json:"swimlane,omitempty"    doc:"Target swimlane name" example:"In Progress"`
	SwimlaneID   *uint  `json:"swimlane_id,omitempty" doc:"Target swimlane ID"`
}

func setIssueRef(shortcut string, issue *model.Issue) *model.Issue {
	issue.Ref = fmt.Sprintf("%s-%d", shortcut, issue.Number)
	return issue
}

func setIssueRefs(shortcut string, issues []model.Issue) []model.Issue {
	for i := range issues {
		issues[i].Ref = fmt.Sprintf("%s-%d", shortcut, issues[i].Number)
	}
	return issues
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

	if in.EpicRef != nil && *in.EpicRef != "" {
		epic, err := s.resolveEpicRef(ctx, p, *in.EpicRef)
		if err != nil {
			return nil, err
		}
		issue.EpicID = &epic.ID
	}

	if err := s.stores.Issues.Create(ctx, issue); err != nil {
		return nil, err
	}
	s.audit.record(ctx, p.ID, "issue", fmt.Sprintf("%s-%d", p.Shortcut, issue.Number), "created", issue.Name, "")
	return setIssueRef(p.Shortcut, issue), nil
}

// changedIssueFields lists the fields that differ between two issue snapshots,
// e.g. "priority, epic", for a human-readable audit detail.
func changedIssueFields(before, after *model.Issue) string {
	var changed []string
	if before.Name != after.Name {
		changed = append(changed, "name")
	}
	if !equalStringPtr(before.Description, after.Description) {
		changed = append(changed, "description")
	}
	if before.Type != after.Type {
		changed = append(changed, "type")
	}
	if before.Priority != after.Priority {
		changed = append(changed, "priority")
	}
	if !equalUintPtr(before.EpicID, after.EpicID) {
		changed = append(changed, "epic")
	}
	return strings.Join(changed, ", ")
}

func equalStringPtr(a, b *string) bool {
	if a == nil || b == nil {
		return a == b
	}
	return *a == *b
}

func equalUintPtr(a, b *uint) bool {
	if a == nil || b == nil {
		return a == b
	}
	return *a == *b
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
	issue, err := s.stores.Issues.GetByRef(ctx, p.ID, number)
	if err != nil {
		return nil, err
	}
	setIssueRef(p.Shortcut, issue)
	if issue.EpicID != nil {
		if epic, err := s.stores.Epics.GetByID(ctx, *issue.EpicID); err == nil {
			setEpicRef(p.Shortcut, epic)
			issue.EpicRef = epic.Ref
			issue.EpicName = epic.Name
		}
	}
	return issue, nil
}

func (s *IssueService) List(ctx context.Context, projectSlug string) ([]model.Issue, error) {
	p, err := s.stores.Projects.GetBySlug(ctx, projectSlug)
	if err != nil {
		return nil, fmt.Errorf("project %q not found", projectSlug)
	}
	issues, err := s.stores.Issues.ListByProject(ctx, p.ID)
	if err != nil {
		return nil, err
	}
	return setIssueRefs(p.Shortcut, issues), nil
}

func (s *IssueService) Move(ctx context.Context, in MoveIssueInput) (*model.Issue, error) {
	issue, err := s.Get(ctx, in.ProjectSlug, in.IssueRef)
	if err != nil {
		return nil, err
	}

	var sl *model.Swimlane
	if in.SwimlaneID != nil {
		sl, err = s.stores.Swimlanes.GetByID(ctx, *in.SwimlaneID)
		if err != nil {
			return nil, fmt.Errorf("swimlane %d not found", *in.SwimlaneID)
		}
	} else {
		p, err := s.stores.Projects.GetBySlug(ctx, in.ProjectSlug)
		if err != nil {
			return nil, err
		}
		sl, err = s.stores.Swimlanes.GetByName(ctx, p.ID, in.SwimlaneName)
		if err != nil {
			return nil, fmt.Errorf("swimlane %q not found in project", in.SwimlaneName)
		}
	}

	p, err := s.stores.Projects.GetBySlug(ctx, in.ProjectSlug)
	if err != nil {
		return nil, err
	}

	fromName := ""
	if from, err := s.stores.Swimlanes.GetByID(ctx, issue.SwimlaneID); err == nil {
		fromName = from.Name
	}

	issue.SwimlaneID = sl.ID
	if err := s.stores.Issues.Update(ctx, issue); err != nil {
		return nil, err
	}
	s.audit.record(ctx, p.ID, "issue", fmt.Sprintf("%s-%d", p.Shortcut, issue.Number), "moved", issue.Name, fmt.Sprintf("%s → %s", fromName, sl.Name))
	return setIssueRef(p.Shortcut, issue), nil
}

type UpdateIssueInput struct {
	ProjectSlug string          `json:"-"`
	IssueRef    string          `json:"-"`
	Name        string          `json:"name"        doc:"Updated issue title"    example:"Fix login timeout (revised)"`
	Description *string         `json:"description,omitempty" doc:"Updated description"`
	Type        model.IssueType `json:"type"        doc:"Updated issue type"     example:"BUG"`
	Priority    model.Priority  `json:"priority"    doc:"Updated priority level" example:"IMMEDIATE"`
	EpicRef     *string         `json:"epic_ref,omitempty" doc:"Epic reference to assign, or empty string to unassign, e.g. PROJ-E1" example:"PROJ-E1"`
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
	before := *issue
	issue.Name = in.Name
	issue.Description = in.Description
	issue.Type = in.Type
	issue.Priority = in.Priority
	if in.EpicRef != nil {
		if *in.EpicRef == "" {
			issue.EpicID = nil
		} else {
			epic, err := s.resolveEpicRef(ctx, p, *in.EpicRef)
			if err != nil {
				return nil, err
			}
			issue.EpicID = &epic.ID
		}
	}
	if err := s.stores.Issues.Update(ctx, issue); err != nil {
		return nil, err
	}
	s.audit.record(ctx, p.ID, "issue", fmt.Sprintf("%s-%d", p.Shortcut, issue.Number), "updated", issue.Name, changedIssueFields(&before, issue))
	return setIssueRef(p.Shortcut, issue), nil
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
	issues, err := s.stores.Issues.ListFiltered(ctx, p.ID, store.IssueStoreFilter{
		Type:       in.Type,
		Priority:   in.Priority,
		SwimlaneID: in.SwimlaneID,
		EpicID:     in.EpicID,
	})
	if err != nil {
		return nil, err
	}
	return setIssueRefs(p.Shortcut, issues), nil
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
	s.audit.record(ctx, p.ID, "issue", fmt.Sprintf("%s-%d", p.Shortcut, issue.Number), "deleted", issue.Name, "")
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
