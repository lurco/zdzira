package model

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type Project struct {
	ID           uint    `gorm:"primarykey"        json:"id"`
	Name         string  `gorm:"unique;not null"   json:"name"`
	Slug         string  `gorm:"unique;not null"   json:"slug"`
	Shortcut     string  `gorm:"unique;not null"   json:"shortcut"`
	Description  *string `                         json:"description,omitempty"`
	IssueCounter uint    `gorm:"not null;default:0" json:"issue_counter"`
	EpicCounter  uint    `gorm:"not null;default:0" json:"epic_counter"`
	Timestamps
	SoftDelete

	Swimlanes []Swimlane `gorm:"foreignKey:ProjectID" json:"-"`
	Issues    []Issue    `gorm:"foreignKey:ProjectID" json:"-"`
	Epics     []Epic     `gorm:"foreignKey:ProjectID" json:"-"`
}

type Swimlane struct {
	ID        uint   `gorm:"primarykey"     json:"id"`
	ProjectID uint   `gorm:"not null;index" json:"project_id"`
	Name      string `gorm:"not null"       json:"name"`
	Position  uint   `gorm:"not null"       json:"position"`
	SoftDelete
}

type Epic struct {
	ID          uint    `gorm:"primarykey"     json:"id"`
	Number      uint    `gorm:"not null"       json:"number"`
	Ref         string  `gorm:"-"              json:"ref"`
	Name        string  `gorm:"not null"       json:"name"`
	Description *string `                      json:"description,omitempty"`
	ProjectID   uint    `gorm:"not null;index" json:"project_id"`
	Timestamps
	SoftDelete
}

type Issue struct {
	ID          uint      `gorm:"primarykey"     json:"id"`
	Number      uint      `gorm:"not null"       json:"number"`
	Ref         string    `gorm:"-"              json:"ref"`
	Type        IssueType `gorm:"not null"       json:"type"`
	Priority    Priority  `gorm:"not null"       json:"priority"`
	Name        string    `gorm:"not null"       json:"name"`
	Description *string   `                      json:"description,omitempty"`
	ProjectID   uint      `gorm:"not null;index" json:"project_id"`
	EpicID      *uint     `gorm:"index"          json:"epic_id,omitempty"`
	SwimlaneID  uint      `gorm:"not null;index" json:"swimlane_id"`
	Position    uint      `gorm:"not null"       json:"position"`
	Timestamps
	SoftDelete
}

// Link has no soft delete — deletions are permanent.
type Link struct {
	ID     uint     `gorm:"primarykey"                       json:"id"`
	Type   LinkType `gorm:"not null"                         json:"type"`
	IssueA uint     `gorm:"not null;column:issue_a;index"    json:"issue_a"`
	IssueB uint     `gorm:"not null;column:issue_b;index"    json:"issue_b"`
}

type Comment struct {
	ID        uint   `gorm:"primarykey" json:"id"`
	Contents  string `gorm:"not null"   json:"contents"`
	IssueID   *uint  `gorm:"index"      json:"issue_id,omitempty"`
	EpicID    *uint  `gorm:"index"      json:"epic_id,omitempty"`
	ProjectID *uint  `gorm:"index"      json:"project_id,omitempty"`
	Timestamps
	SoftDelete
}

func (c *Comment) BeforeCreate(tx *gorm.DB) error { return c.validateParent() }
func (c *Comment) BeforeSave(tx *gorm.DB) error   { return c.validateParent() }

func (c *Comment) validateParent() error {
	set := 0
	if c.IssueID != nil {
		set++
	}
	if c.EpicID != nil {
		set++
	}
	if c.ProjectID != nil {
		set++
	}
	if set != 1 {
		return errors.New("comment must belong to exactly one of: issue, epic, or project")
	}
	return nil
}

// AuditEntry is append-only — no UpdatedAt, no DeletedAt.
type AuditEntry struct {
	ID         uint      `gorm:"primarykey"              json:"id"`
	ProjectID  uint      `gorm:"not null;index"          json:"project_id"`
	EntityType string    `gorm:"not null"                json:"entity_type"`
	Ref        string    `gorm:"not null"                json:"ref"`
	Action     string    `gorm:"not null"                json:"action"`
	CreatedAt  time.Time `gorm:"not null;autoCreateTime" json:"created_at"`
}
