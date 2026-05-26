package service

import (
	"context"
	"zdzira/internal/model"
	"zdzira/internal/store"
)

type SwimlaneService struct {
	stores *store.Stores
}

func (s *SwimlaneService) ListByProject(ctx context.Context, projectID uint) ([]model.Swimlane, error) {
	return s.stores.Swimlanes.ListByProject(ctx, projectID)
}
