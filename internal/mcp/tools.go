package mcp

import (
	"context"
	"zdzira/internal/model"
	"zdzira/internal/service"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

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
			return mcp.NewToolResultJSON(projects)
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
			return mcp.NewToolResultJSON(epics)
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
			mcp.WithDescription("List all issues in a project."),
			mcp.WithString("project", mcp.Required(), mcp.Description("Project slug")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			slug, err := req.RequireString("project")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			issues, err := svcs.Issues.List(ctx, slug)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return mcp.NewToolResultJSON(issues)
		},
	)

	s.AddTool(
		mcp.NewTool("get_issue",
			mcp.WithDescription("Get an issue by its ref (e.g. PROJ-42)."),
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
			return mcp.NewToolResultJSON(issue)
		},
	)

	s.AddTool(
		mcp.NewTool("create_issue",
			mcp.WithDescription("Create a new issue in a project. Returns the issue ref (e.g. PROJ-42)."),
			mcp.WithString("project", mcp.Required(), mcp.Description("Project slug")),
			mcp.WithString("name", mcp.Required(), mcp.Description("Issue title")),
			mcp.WithString("type", mcp.Required(), mcp.Description("Issue type: TASK, BUG, or STORY")),
			mcp.WithString("priority", mcp.Required(), mcp.Description("Priority: LOW, HIGH, or IMMEDIATE")),
			mcp.WithString("description", mcp.Description("Optional description")),
			mcp.WithString("epic_ref", mcp.Description("Optional epic ref to associate, e.g. PROJ-E1")),
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
			return mcp.NewToolResultJSON(issue)
		},
	)
}

func registerUpdateIssueTools(s *server.MCPServer, svcs *service.Services) {
	s.AddTool(
		mcp.NewTool("update_issue",
			mcp.WithDescription("Update an issue's name, type, priority, or description."),
			mcp.WithString("project", mcp.Required(), mcp.Description("Project slug")),
			mcp.WithString("issue_ref", mcp.Required(), mcp.Description("Issue reference, e.g. PROJ-42")),
			mcp.WithString("name", mcp.Required(), mcp.Description("Updated issue title")),
			mcp.WithString("type", mcp.Required(), mcp.Description("Issue type: TASK, BUG, or STORY")),
			mcp.WithString("priority", mcp.Required(), mcp.Description("Priority: LOW, HIGH, or IMMEDIATE")),
			mcp.WithString("description", mcp.Description("Optional updated description")),
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
			in := service.UpdateIssueInput{
				ProjectSlug: slug,
				IssueRef:    ref,
				Name:        name,
				Type:        model.IssueType(issueType),
				Priority:    model.Priority(priority),
			}
			if desc := req.GetString("description", ""); desc != "" {
				in.Description = &desc
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
			mcp.WithDescription("List all swimlanes in a project."),
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
			return mcp.NewToolResultJSON(swimlanes)
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
			mcp.WithDescription("Add a comment to an issue."),
			mcp.WithString("project", mcp.Required(), mcp.Description("Project slug")),
			mcp.WithString("issue_ref", mcp.Required(), mcp.Description("Issue reference, e.g. PROJ-42")),
			mcp.WithString("contents", mcp.Required(), mcp.Description("Comment text")),
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
			contents, err := req.RequireString("contents")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			c, err := svcs.Comments.AddToIssue(ctx, slug, ref, contents)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return mcp.NewToolResultJSON(c)
		},
	)

	s.AddTool(
		mcp.NewTool("list_comments",
			mcp.WithDescription("List all comments on an issue."),
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
			comments, err := svcs.Comments.ListForIssue(ctx, slug, ref)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return mcp.NewToolResultJSON(comments)
		},
	)
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
			return mcp.NewToolResultJSON(links)
		},
	)
}
