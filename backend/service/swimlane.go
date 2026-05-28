package service

import (
	"context"
	"fmt"
	"zdzira/backend/model"
	"zdzira/backend/store"
)

type SwimlaneService struct {
	stores *store.Stores
}

type CreateSwimlaneInput struct {
	ProjectSlug string `json:"-"`
	Name        string `json:"name" doc:"Swimlane name" example:"In Review"`
}

type RenameSwimlaneInput struct {
	ProjectSlug string `json:"-"`
	ID          uint   `json:"-"`
	Name        string `json:"name" doc:"New swimlane name" example:"In Review"`
}

func (s *SwimlaneService) ListByProject(ctx context.Context, projectID uint) ([]model.Swimlane, error) {
	return s.stores.Swimlanes.ListByProject(ctx, projectID)
}

func (s *SwimlaneService) ListForProject(ctx context.Context, projectSlug string) ([]model.Swimlane, error) {
	p, err := s.stores.Projects.GetBySlug(ctx, projectSlug)
	if err != nil {
		return nil, fmt.Errorf("project %q not found", projectSlug)
	}
	return s.stores.Swimlanes.ListByProject(ctx, p.ID)
}

func (s *SwimlaneService) Create(ctx context.Context, in CreateSwimlaneInput) (*model.Swimlane, error) {
	p, err := s.stores.Projects.GetBySlug(ctx, in.ProjectSlug)
	if err != nil {
		return nil, fmt.Errorf("project %q not found", in.ProjectSlug)
	}
	existing, err := s.stores.Swimlanes.ListByProject(ctx, p.ID)
	if err != nil {
		return nil, err
	}
	sl := &model.Swimlane{
		ProjectID: p.ID,
		Name:      in.Name,
		Position:  uint(len(existing)) + 1,
	}
	if err := s.stores.Swimlanes.Create(ctx, sl); err != nil {
		return nil, err
	}
	return sl, nil
}

func (s *SwimlaneService) Delete(ctx context.Context, projectSlug string, id uint) error {
	p, err := s.stores.Projects.GetBySlug(ctx, projectSlug)
	if err != nil {
		return fmt.Errorf("project %q not found", projectSlug)
	}
	swimlanes, err := s.stores.Swimlanes.ListByProject(ctx, p.ID)
	if err != nil {
		return err
	}
	for _, sl := range swimlanes {
		if sl.ID == id {
			return s.stores.Swimlanes.Delete(ctx, id)
		}
	}
	return fmt.Errorf("swimlane %d not found in project", id)
}

func (s *SwimlaneService) Rename(ctx context.Context, in RenameSwimlaneInput) (*model.Swimlane, error) {
	p, err := s.stores.Projects.GetBySlug(ctx, in.ProjectSlug)
	if err != nil {
		return nil, fmt.Errorf("project %q not found", in.ProjectSlug)
	}
	swimlanes, err := s.stores.Swimlanes.ListByProject(ctx, p.ID)
	if err != nil {
		return nil, err
	}
	var target *model.Swimlane
	for i := range swimlanes {
		if swimlanes[i].ID == in.ID {
			target = &swimlanes[i]
			break
		}
	}
	if target == nil {
		return nil, fmt.Errorf("swimlane %d not found in project", in.ID)
	}
	target.Name = in.Name
	return target, s.stores.Swimlanes.Update(ctx, target)
}
