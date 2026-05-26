package service

import (
	"context"
	"fmt"
	"zdzira/internal/model"
	"zdzira/internal/store"
)

type EpicService struct {
	stores *store.Stores
}

type CreateEpicInput struct {
	ProjectSlug string
	Name        string
	Description *string
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

func (s *EpicService) Delete(ctx context.Context, projectSlug, ref string) error {
	e, err := s.Get(ctx, projectSlug, ref)
	if err != nil {
		return err
	}
	return s.stores.Epics.Delete(ctx, e.ID)
}
