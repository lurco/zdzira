package service

import (
	"context"
	"fmt"
	"zdzira/internal/model"
	"zdzira/internal/store"
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
	s.audit.record(ctx, p.ID, "epic", fmt.Sprintf("%s-E%d", p.Shortcut, e.Number), "created")
	return e, nil
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
	return s.stores.Epics.GetByRef(ctx, p.ID, number)
}

func (s *EpicService) List(ctx context.Context, projectSlug string) ([]model.Epic, error) {
	p, err := s.stores.Projects.GetBySlug(ctx, projectSlug)
	if err != nil {
		return nil, fmt.Errorf("project %q not found", projectSlug)
	}
	return s.stores.Epics.ListByProject(ctx, p.ID)
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
	e.Name = in.Name
	e.Description = in.Description
	if err := s.stores.Epics.Update(ctx, e); err != nil {
		return nil, err
	}
	s.audit.record(ctx, p.ID, "epic", fmt.Sprintf("%s-E%d", p.Shortcut, e.Number), "updated")
	return e, nil
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
	s.audit.record(ctx, p.ID, "epic", fmt.Sprintf("%s-E%d", p.Shortcut, e.Number), "deleted")
	return nil
}
