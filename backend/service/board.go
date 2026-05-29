package service

import (
	"context"
	"fmt"
	"zdzira/backend/model"
	"zdzira/backend/store"
)

type BoardService struct {
	stores *store.Stores
}

type BoardLane struct {
	model.Swimlane
	Issues []model.Issue `json:"issues"`
}

type BoardView struct {
	Swimlanes []BoardLane  `json:"swimlanes"`
	Epics     []model.Epic `json:"epics"`
	EpicRef   string       `json:"epic_ref,omitempty"`
}

type BoardFilter struct {
	EpicRef string
}

func (s *BoardService) Get(ctx context.Context, projectSlug string, filter BoardFilter) (*BoardView, error) {
	p, err := s.stores.Projects.GetBySlug(ctx, projectSlug)
	if err != nil {
		return nil, fmt.Errorf("project %q not found", projectSlug)
	}

	swimlanes, err := s.stores.Swimlanes.ListByProject(ctx, p.ID)
	if err != nil {
		return nil, err
	}
	issues, err := s.stores.Issues.ListByProject(ctx, p.ID)
	if err != nil {
		return nil, err
	}
	epics, err := s.stores.Epics.ListByProject(ctx, p.ID)
	if err != nil {
		return nil, err
	}

	setIssueRefs(p.Shortcut, issues)
	setEpicRefs(p.Shortcut, epics)

	if filter.EpicRef != "" {
		issues = filterIssuesByEpic(issues, epics, filter.EpicRef)
	}

	issuesByLane := make(map[uint][]model.Issue, len(swimlanes))
	for _, issue := range issues {
		issuesByLane[issue.SwimlaneID] = append(issuesByLane[issue.SwimlaneID], issue)
	}

	lanes := make([]BoardLane, len(swimlanes))
	for i, sl := range swimlanes {
		lanes[i] = BoardLane{Swimlane: sl, Issues: issuesByLane[sl.ID]}
	}

	return &BoardView{Swimlanes: lanes, Epics: epics, EpicRef: filter.EpicRef}, nil
}

func filterIssuesByEpic(issues []model.Issue, epics []model.Epic, epicRef string) []model.Issue {
	var matchedEpicID uint
	for _, e := range epics {
		if e.Ref == epicRef {
			matchedEpicID = e.ID
			break
		}
	}
	if matchedEpicID == 0 {
		return nil
	}
	filtered := issues[:0]
	for _, i := range issues {
		if i.EpicID != nil && *i.EpicID == matchedEpicID {
			filtered = append(filtered, i)
		}
	}
	return filtered
}
