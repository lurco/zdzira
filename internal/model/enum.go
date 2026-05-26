package model

type IssueType string

const (
	IssueTypeTask  IssueType = "TASK"
	IssueTypeBug   IssueType = "BUG"
	IssueTypeStory IssueType = "STORY"
)

type Priority string

const (
	PriorityLow       Priority = "LOW"
	PriorityHigh      Priority = "HIGH"
	PriorityImmediate Priority = "IMMEDIATE"
)

type LinkType string

const (
	LinkTypeDuplicates LinkType = "DUPLICATES"
	LinkTypeIsPartOf   LinkType = "IS_PART_OF"
	LinkTypeBlocks     LinkType = "BLOCKS"
	LinkTypeRelatesTo  LinkType = "RELATES_TO"
)
