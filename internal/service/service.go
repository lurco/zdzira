package service

import (
	"zdzira/internal/store"
)

type Services struct {
	Projects  *ProjectService
	Epics     *EpicService
	Issues    *IssueService
	Comments  *CommentService
	Links     *LinkService
	Swimlanes *SwimlaneService
}

func New(stores *store.Stores) *Services {
	return &Services{
		Projects:  &ProjectService{stores: stores},
		Epics:     &EpicService{stores: stores},
		Issues:    &IssueService{stores: stores},
		Comments:  &CommentService{stores: stores},
		Links:     &LinkService{stores: stores},
		Swimlanes: &SwimlaneService{stores: stores},
	}
}
