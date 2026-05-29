package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"zdzira/backend/model"
	"zdzira/backend/service"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// newListResult returns a slice as JSON text only. The MCP spec requires
// structured content to be a JSON object, so list tools (which return arrays)
// must not set it — clients reject a top-level array.
func newListResult[T any](items []T) (*mcp.CallToolResult, error) {
	b, err := json.Marshal(items)
	if err != nil {
		return nil, fmt.Errorf("marshal list result: %w", err)
	}
	return mcp.NewToolResultText(string(b)), nil
}

// issueSummary is the AI-facing issue shape: human-readable names instead
// of internal IDs, only the fields an agent actually needs.
type issueSummary struct {
	Ref          string `json:"ref"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	Priority     string `json:"priority"`
	SwimlaneName string `json:"swimlane_name"`
	EpicRef      string `json:"epic_ref,omitempty"`
	Description  string `json:"description,omitempty"`
}

func toIssueSummary(iss model.Issue, slName string) issueSummary {
	desc := ""
	if iss.Description != nil {
		desc = *iss.Description
	}
	return issueSummary{
		Ref:          iss.Ref,
		Name:         iss.Name,
		Type:         string(iss.Type),
		Priority:     string(iss.Priority),
		SwimlaneName: slName,
		EpicRef:      iss.EpicRef,
		Description:  desc,
	}
}

// swimlaneNameMap builds id→name and lowercase-name→id maps for a project.
func swimlaneNameMap(swimlanes []model.Swimlane) (byID map[uint]string, byName map[string]uint) {
	byID = make(map[uint]string, len(swimlanes))
	byName = make(map[string]uint, len(swimlanes))
	for _, sl := range swimlanes {
		byID[sl.ID] = sl.Name
		byName[strings.ToLower(sl.Name)] = sl.ID
	}
	return
}

func registerProjectTools(s *server.MCPServer, svcs *service.Services) {
	s.AddTool(
		mcp.NewTool("list_projects",
			mcp.WithDescription("List all projects."),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			projects, err := svcs.Projects.List(ctx)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return newListResult(projects)
		},
	)

	s.AddTool(
		mcp.NewTool("get_project",
			mcp.WithDescription("Get a project by its slug."),
			mcp.WithString("project", mcp.Required(), mcp.Description("Project slug, e.g. \"my-project\"")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			slug, err := req.RequireString("project")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			p, err := svcs.Projects.Get(ctx, slug)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return mcp.NewToolResultJSON(p)
		},
	)

	s.AddTool(
		mcp.NewTool("create_project",
			mcp.WithDescription("Create a new project. Swimlanes are seeded automatically (Backlog, In Progress, Done)."),
			mcp.WithString("name", mcp.Required(), mcp.Description("Project name")),
			mcp.WithString("shortcut", mcp.Required(), mcp.Description("Uppercase issue prefix, e.g. \"PROJ\"")),
			mcp.WithString("description", mcp.Description("Optional project description")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			name, err := req.RequireString("name")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			shortcut, err := req.RequireString("shortcut")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			desc := req.GetString("description", "")
			in := service.CreateProjectInput{Name: name, Shortcut: shortcut}
			if desc != "" {
				in.Description = &desc
			}
			p, err := svcs.Projects.Create(ctx, in)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return mcp.NewToolResultJSON(p)
		},
	)
}

func registerEpicTools(s *server.MCPServer, svcs *service.Services) {
	s.AddTool(
		mcp.NewTool("list_epics",
			mcp.WithDescription("List all epics in a project."),
			mcp.WithString("project", mcp.Required(), mcp.Description("Project slug")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			slug, err := req.RequireString("project")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			epics, err := svcs.Epics.List(ctx, slug)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return newListResult(epics)
		},
	)

	s.AddTool(
		mcp.NewTool("create_epic",
			mcp.WithDescription("Create a new epic in a project. Returns the epic ref (e.g. PROJ-E1)."),
			mcp.WithString("project", mcp.Required(), mcp.Description("Project slug")),
			mcp.WithString("name", mcp.Required(), mcp.Description("Epic name")),
			mcp.WithString("description", mcp.Description("Optional description")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			slug, err := req.RequireString("project")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			name, err := req.RequireString("name")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			desc := req.GetString("description", "")
			in := service.CreateEpicInput{ProjectSlug: slug, Name: name}
			if desc != "" {
				in.Description = &desc
			}
			e, err := svcs.Epics.Create(ctx, in)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return mcp.NewToolResultJSON(e)
		},
	)

	s.AddTool(
		mcp.NewTool("get_epic",
			mcp.WithDescription("Get an epic by its ref (e.g. PROJ-E1)."),
			mcp.WithString("project", mcp.Required(), mcp.Description("Project slug")),
			mcp.WithString("epic_ref", mcp.Required(), mcp.Description("Epic reference, e.g. PROJ-E1")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			slug, err := req.RequireString("project")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			ref, err := req.RequireString("epic_ref")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			e, err := svcs.Epics.Get(ctx, slug, ref)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return mcp.NewToolResultJSON(e)
		},
	)
}

func registerIssueTools(s *server.MCPServer, svcs *service.Services) {
	s.AddTool(
		mcp.NewTool("list_issues",
			mcp.WithDescription("List issues in a project. Returns human-readable swimlane names. Use optional filters to narrow scope."),
			mcp.WithString("project", mcp.Required(), mcp.Description("Project slug")),
			mcp.WithString("swimlane", mcp.Description("Filter by swimlane name, e.g. \"In Progress\", \"Backlog\", \"Done\"")),
			mcp.WithString("type", mcp.Description("Filter by type: TASK, BUG, or STORY")),
			mcp.WithString("priority", mcp.Description("Filter by priority: LOW, HIGH, or IMMEDIATE")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			slug, err := req.RequireString("project")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			swimlanes, err := svcs.Swimlanes.ListForProject(ctx, slug)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			byID, byName := swimlaneNameMap(swimlanes)

			in := service.IssueFilterInput{ProjectSlug: slug}
			if slName := req.GetString("swimlane", ""); slName != "" {
				id, ok := byName[strings.ToLower(slName)]
				if !ok {
					return mcp.NewToolResultError(fmt.Sprintf("swimlane %q not found", slName)), nil
				}
				in.SwimlaneID = &id
			}
			if t := req.GetString("type", ""); t != "" {
				it := model.IssueType(strings.ToUpper(t))
				in.Type = &it
			}
			if p := req.GetString("priority", ""); p != "" {
				pr := model.Priority(strings.ToUpper(p))
				in.Priority = &pr
			}

			issues, err := svcs.Issues.Filter(ctx, in)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			summaries := make([]issueSummary, len(issues))
			for i, iss := range issues {
				summaries[i] = toIssueSummary(iss, byID[iss.SwimlaneID])
			}
			return newListResult(summaries)
		},
	)

	s.AddTool(
		mcp.NewTool("get_board",
			mcp.WithDescription("Get the full board for a project — all swimlanes with their issues. Best starting point for understanding the current project state."),
			mcp.WithString("project", mcp.Required(), mcp.Description("Project slug")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			slug, err := req.RequireString("project")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			view, err := svcs.Board.Get(ctx, slug, service.BoardFilter{})
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			type laneSummary struct {
				Name   string         `json:"name"`
				Issues []issueSummary `json:"issues"`
			}
			lanes := make([]laneSummary, len(view.Swimlanes))
			for i, lane := range view.Swimlanes {
				issueSums := make([]issueSummary, len(lane.Issues))
				for j, iss := range lane.Issues {
					issueSums[j] = toIssueSummary(iss, lane.Name)
				}
				lanes[i] = laneSummary{Name: lane.Name, Issues: issueSums}
			}
			b, err := json.Marshal(lanes)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return mcp.NewToolResultText(string(b)), nil
		},
	)

	s.AddTool(
		mcp.NewTool("get_issue",
			mcp.WithDescription("Get full details for a single issue by its ref (e.g. PROJ-42)."),
			mcp.WithString("project", mcp.Required(), mcp.Description("Project slug")),
			mcp.WithString("issue_ref", mcp.Required(), mcp.Description("Issue reference, e.g. PROJ-42")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			slug, err := req.RequireString("project")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			ref, err := req.RequireString("issue_ref")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			issue, err := svcs.Issues.Get(ctx, slug, ref)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			swimlanes, _ := svcs.Swimlanes.ListForProject(ctx, slug)
			byID, _ := swimlaneNameMap(swimlanes)
			return mcp.NewToolResultJSON(toIssueSummary(*issue, byID[issue.SwimlaneID]))
		},
	)

	s.AddTool(
		mcp.NewTool("create_issue",
			mcp.WithDescription("Create a new issue. Returns the issue ref (e.g. PROJ-42). New issues land in the first swimlane (usually Backlog)."),
			mcp.WithString("project", mcp.Required(), mcp.Description("Project slug")),
			mcp.WithString("name", mcp.Required(), mcp.Description("Issue title")),
			mcp.WithString("type", mcp.Required(), mcp.Description("TASK, BUG, or STORY")),
			mcp.WithString("priority", mcp.Required(), mcp.Description("LOW, HIGH, or IMMEDIATE")),
			mcp.WithString("description", mcp.Description("Optional description")),
			mcp.WithString("epic_ref", mcp.Description("Optional epic ref, e.g. PROJ-E1")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			slug, err := req.RequireString("project")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			name, err := req.RequireString("name")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			issueType, err := req.RequireString("type")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			priority, err := req.RequireString("priority")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			in := service.CreateIssueInput{
				ProjectSlug: slug,
				Name:        name,
				Type:        model.IssueType(issueType),
				Priority:    model.Priority(priority),
			}
			if desc := req.GetString("description", ""); desc != "" {
				in.Description = &desc
			}
			if epicRef := req.GetString("epic_ref", ""); epicRef != "" {
				in.EpicRef = &epicRef
			}
			issue, err := svcs.Issues.Create(ctx, in)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return mcp.NewToolResultJSON(issue)
		},
	)

	s.AddTool(
		mcp.NewTool("move_issue",
			mcp.WithDescription("Move an issue to a different swimlane (e.g. 'In Progress', 'Done')."),
			mcp.WithString("project", mcp.Required(), mcp.Description("Project slug")),
			mcp.WithString("issue_ref", mcp.Required(), mcp.Description("Issue reference, e.g. PROJ-42")),
			mcp.WithString("swimlane", mcp.Required(), mcp.Description("Target swimlane name, e.g. \"In Progress\"")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			slug, err := req.RequireString("project")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			ref, err := req.RequireString("issue_ref")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			swimlane, err := req.RequireString("swimlane")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			issue, err := svcs.Issues.Move(ctx, service.MoveIssueInput{
				ProjectSlug:  slug,
				IssueRef:     ref,
				SwimlaneName: swimlane,
			})
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return mcp.NewToolResultText(fmt.Sprintf("moved %s to %q", issue.Ref, swimlane)), nil
		},
	)
}

func registerUpdateIssueTools(s *server.MCPServer, svcs *service.Services) {
	s.AddTool(
		mcp.NewTool("update_issue",
			mcp.WithDescription("Update an issue. All fields are optional — omit any field to keep its current value."),
			mcp.WithString("project", mcp.Required(), mcp.Description("Project slug")),
			mcp.WithString("issue_ref", mcp.Required(), mcp.Description("Issue reference, e.g. PROJ-42")),
			mcp.WithString("name", mcp.Description("New title (omit to keep current)")),
			mcp.WithString("type", mcp.Description("New type: TASK, BUG, or STORY (omit to keep current)")),
			mcp.WithString("priority", mcp.Description("New priority: LOW, HIGH, or IMMEDIATE (omit to keep current)")),
			mcp.WithString("description", mcp.Description("New description (omit to keep current)")),
			mcp.WithString("epic_ref", mcp.Description("Epic ref to attach (empty string to detach, omit for no change)")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			slug, err := req.RequireString("project")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			ref, err := req.RequireString("issue_ref")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			// Fetch current state; apply only what was explicitly provided.
			current, err := svcs.Issues.Get(ctx, slug, ref)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			in := service.UpdateIssueInput{
				ProjectSlug: slug,
				IssueRef:    ref,
				Name:        current.Name,
				Type:        current.Type,
				Priority:    current.Priority,
				Description: current.Description,
			}
			if name := req.GetString("name", ""); name != "" {
				in.Name = name
			}
			if t := req.GetString("type", ""); t != "" {
				in.Type = model.IssueType(strings.ToUpper(t))
			}
			if p := req.GetString("priority", ""); p != "" {
				in.Priority = model.Priority(strings.ToUpper(p))
			}
			if desc := req.GetString("description", ""); desc != "" {
				in.Description = &desc
			}
			if args := req.GetArguments(); args != nil {
				if _, provided := args["epic_ref"]; provided {
					er := req.GetString("epic_ref", "")
					in.EpicRef = &er
				}
			}

			issue, err := svcs.Issues.Update(ctx, in)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return mcp.NewToolResultJSON(issue)
		},
	)
}

func registerSwimlaneTools(s *server.MCPServer, svcs *service.Services) {
	s.AddTool(
		mcp.NewTool("list_swimlanes",
			mcp.WithDescription("List all swimlanes (board columns) in a project with their names and colors."),
			mcp.WithString("project", mcp.Required(), mcp.Description("Project slug")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			slug, err := req.RequireString("project")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			swimlanes, err := svcs.Swimlanes.ListForProject(ctx, slug)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return newListResult(swimlanes)
		},
	)
}

func registerUpdateEpicTools(s *server.MCPServer, svcs *service.Services) {
	s.AddTool(
		mcp.NewTool("update_epic",
			mcp.WithDescription("Update an epic's name or description."),
			mcp.WithString("project", mcp.Required(), mcp.Description("Project slug")),
			mcp.WithString("epic_ref", mcp.Required(), mcp.Description("Epic reference, e.g. PROJ-E1")),
			mcp.WithString("name", mcp.Required(), mcp.Description("Updated epic name")),
			mcp.WithString("description", mcp.Description("Optional updated description")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			slug, err := req.RequireString("project")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			ref, err := req.RequireString("epic_ref")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			name, err := req.RequireString("name")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			in := service.UpdateEpicInput{
				ProjectSlug: slug,
				EpicRef:     ref,
				Name:        name,
			}
			if desc := req.GetString("description", ""); desc != "" {
				in.Description = &desc
			}
			e, err := svcs.Epics.Update(ctx, in)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return mcp.NewToolResultJSON(e)
		},
	)
}

func registerDeleteIssueTools(s *server.MCPServer, svcs *service.Services) {
	s.AddTool(
		mcp.NewTool("delete_issue",
			mcp.WithDescription("Delete an issue by its ref. The deletion is soft — data is retained but the issue is no longer accessible."),
			mcp.WithString("project", mcp.Required(), mcp.Description("Project slug")),
			mcp.WithString("issue_ref", mcp.Required(), mcp.Description("Issue reference, e.g. PROJ-42")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			slug, err := req.RequireString("project")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			ref, err := req.RequireString("issue_ref")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			if err := svcs.Issues.Delete(ctx, slug, ref); err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return mcp.NewToolResultText("deleted"), nil
		},
	)
}

func registerCommentTools(s *server.MCPServer, svcs *service.Services) {
	s.AddTool(
		mcp.NewTool("add_comment",
			mcp.WithDescription("Add a comment to an issue or an epic. Provide exactly one of issue_ref or epic_ref."),
			mcp.WithString("project", mcp.Required(), mcp.Description("Project slug")),
			mcp.WithString("issue_ref", mcp.Description("Issue reference, e.g. PROJ-42")),
			mcp.WithString("epic_ref", mcp.Description("Epic reference, e.g. PROJ-E1")),
			mcp.WithString("contents", mcp.Required(), mcp.Description("Comment text")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			slug, err := req.RequireString("project")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			contents, err := req.RequireString("contents")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			issueRef := req.GetString("issue_ref", "")
			epicRef := req.GetString("epic_ref", "")
			if err := requireOneRef(issueRef, epicRef); err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			var c *model.Comment
			if issueRef != "" {
				c, err = svcs.Comments.AddToIssue(ctx, slug, issueRef, contents)
			} else {
				c, err = svcs.Comments.AddToEpic(ctx, slug, epicRef, contents)
			}
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return mcp.NewToolResultJSON(c)
		},
	)

	s.AddTool(
		mcp.NewTool("list_comments",
			mcp.WithDescription("List all comments on an issue or an epic. Provide exactly one of issue_ref or epic_ref."),
			mcp.WithString("project", mcp.Required(), mcp.Description("Project slug")),
			mcp.WithString("issue_ref", mcp.Description("Issue reference, e.g. PROJ-42")),
			mcp.WithString("epic_ref", mcp.Description("Epic reference, e.g. PROJ-E1")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			slug, err := req.RequireString("project")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			issueRef := req.GetString("issue_ref", "")
			epicRef := req.GetString("epic_ref", "")
			if err := requireOneRef(issueRef, epicRef); err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			var comments []model.Comment
			if issueRef != "" {
				comments, err = svcs.Comments.ListForIssue(ctx, slug, issueRef)
			} else {
				comments, err = svcs.Comments.ListForEpic(ctx, slug, epicRef)
			}
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return newListResult(comments)
		},
	)
}

// requireOneRef ensures exactly one of issue_ref / epic_ref was provided.
func requireOneRef(issueRef, epicRef string) error {
	switch {
	case issueRef == "" && epicRef == "":
		return fmt.Errorf("provide either issue_ref or epic_ref")
	case issueRef != "" && epicRef != "":
		return fmt.Errorf("provide only one of issue_ref or epic_ref, not both")
	}
	return nil
}

func registerLinkTools(s *server.MCPServer, svcs *service.Services) {
	s.AddTool(
		mcp.NewTool("link_issues",
			mcp.WithDescription("Create a directed link between two issues. For BLOCKS: source blocks target. For IS_PART_OF: source is part of target."),
			mcp.WithString("project", mcp.Required(), mcp.Description("Project slug")),
			mcp.WithString("source_ref", mcp.Required(), mcp.Description("Source issue ref, e.g. PROJ-38")),
			mcp.WithString("target_ref", mcp.Required(), mcp.Description("Target issue ref, e.g. PROJ-42")),
			mcp.WithString("type", mcp.Required(), mcp.Description("Link type: BLOCKS, IS_PART_OF, DUPLICATES, or RELATES_TO")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			slug, err := req.RequireString("project")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			sourceRef, err := req.RequireString("source_ref")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			targetRef, err := req.RequireString("target_ref")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			linkType, err := req.RequireString("type")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			l, err := svcs.Links.Create(ctx, service.CreateLinkInput{
				ProjectSlug: slug,
				SourceRef:   sourceRef,
				TargetRef:   targetRef,
				Type:        model.LinkType(linkType),
			})
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return mcp.NewToolResultJSON(l)
		},
	)

	s.AddTool(
		mcp.NewTool("list_links",
			mcp.WithDescription("List all links for an issue (as source or target)."),
			mcp.WithString("project", mcp.Required(), mcp.Description("Project slug")),
			mcp.WithString("issue_ref", mcp.Required(), mcp.Description("Issue reference, e.g. PROJ-42")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			slug, err := req.RequireString("project")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			ref, err := req.RequireString("issue_ref")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			links, err := svcs.Links.ListForIssue(ctx, slug, ref)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return newListResult(links)
		},
	)
}
