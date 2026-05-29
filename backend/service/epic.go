package service

import (
	"context"
	"fmt"
	"strings"
	"zdzira/backend/model"
	"zdzira/backend/store"
)

type EpicService struct {
	stores *store.Stores
	audit  *AuditService
}

type CreateEpicInput struct {
	ProjectSlug string  `json:"-"`
	Name        string  `json:"name"        doc:"Epic name"                  example:"User authentication overhaul"`
	Description *string `json:"description,omitempty" doc:"Optional epic description"  example:"All issues related to the new OAuth flow"`
}

func setEpicRef(shortcut string, e *model.Epic) *model.Epic {
	e.Ref = fmt.Sprintf("%s-E%d", shortcut, e.Number)
	return e
}

func setEpicRefs(shortcut string, epics []model.Epic) []model.Epic {
	for i := range epics {
		epics[i].Ref = fmt.Sprintf("%s-E%d", shortcut, epics[i].Number)
	}
	return epics
}

func (s *EpicService) Create(ctx context.Context, in CreateEpicInput) (*model.Epic, error) {
	p, err := s.stores.Projects.GetBySlug(ctx, in.ProjectSlug)
	if err != nil {
		return nil, fmt.Errorf("project %q not found", in.ProjectSlug)
	}

	p.EpicCounter++
	if err := s.stores.Projects.Update(ctx, p); err != nil {
		return nil, err
	}

	e := &model.Epic{
		Number:      p.EpicCounter,
		Name:        in.Name,
		Description: in.Description,
		ProjectID:   p.ID,
	}
	if err := s.stores.Epics.Create(ctx, e); err != nil {
		return nil, err
	}
	s.audit.record(ctx, p.ID, "epic", fmt.Sprintf("%s-E%d", p.Shortcut, e.Number), "created", e.Name)
	return setEpicRef(p.Shortcut, e), nil
}

func (s *EpicService) Get(ctx context.Context, projectSlug, ref string) (*model.Epic, error) {
	p, err := s.stores.Projects.GetBySlug(ctx, projectSlug)
	if err != nil {
		return nil, fmt.Errorf("project %q not found", projectSlug)
	}
	number, err := parseEpicRef(p.Shortcut, ref)
	if err != nil {
		return nil, err
	}
	e, err := s.stores.Epics.GetByRef(ctx, p.ID, number)
	if err != nil {
		return nil, err
	}
	return setEpicRef(p.Shortcut, e), nil
}

func (s *EpicService) List(ctx context.Context, projectSlug string) ([]model.Epic, error) {
	p, err := s.stores.Projects.GetBySlug(ctx, projectSlug)
	if err != nil {
		return nil, fmt.Errorf("project %q not found", projectSlug)
	}
	epics, err := s.stores.Epics.ListByProject(ctx, p.ID)
	if err != nil {
		return nil, err
	}
	return setEpicRefs(p.Shortcut, epics), nil
}

type UpdateEpicInput struct {
	ProjectSlug string  `json:"-"`
	EpicRef     string  `json:"-"`
	Name        string  `json:"name"                 doc:"Updated epic name"        example:"Auth overhaul (revised)"`
	Description *string `json:"description,omitempty" doc:"Updated epic description"`
}

func (s *EpicService) Update(ctx context.Context, in UpdateEpicInput) (*model.Epic, error) {
	p, err := s.stores.Projects.GetBySlug(ctx, in.ProjectSlug)
	if err != nil {
		return nil, fmt.Errorf("project %q not found", in.ProjectSlug)
	}
	e, err := s.Get(ctx, in.ProjectSlug, in.EpicRef)
	if err != nil {
		return nil, err
	}
	var changed []string
	if e.Name != in.Name {
		changed = append(changed, "name")
	}
	if !equalStringPtr(e.Description, in.Description) {
		changed = append(changed, "description")
	}
	e.Name = in.Name
	e.Description = in.Description
	if err := s.stores.Epics.Update(ctx, e); err != nil {
		return nil, err
	}
	s.audit.record(ctx, p.ID, "epic", fmt.Sprintf("%s-E%d", p.Shortcut, e.Number), "updated", strings.Join(changed, ", "))
	return setEpicRef(p.Shortcut, e), nil
}

func (s *EpicService) Delete(ctx context.Context, projectSlug, ref string) error {
	p, err := s.stores.Projects.GetBySlug(ctx, projectSlug)
	if err != nil {
		return fmt.Errorf("project %q not found", projectSlug)
	}
	e, err := s.Get(ctx, projectSlug, ref)
	if err != nil {
		return err
	}
	if err := s.stores.Epics.Delete(ctx, e.ID); err != nil {
		return err
	}
	s.audit.record(ctx, p.ID, "epic", fmt.Sprintf("%s-E%d", p.Shortcut, e.Number), "deleted", e.Name)
	return nil
}
