package mcp

import (
	"net/http"
	"zdzira/internal/service"

	"github.com/mark3labs/mcp-go/server"
)

func NewSSEHandler(svcs *service.Services, baseURL string) http.Handler {
	s := server.NewMCPServer("zdzira", "1.0.0",
		server.WithToolCapabilities(true),
	)

	registerProjectTools(s, svcs)
	registerEpicTools(s, svcs)
	registerIssueTools(s, svcs)
	registerCommentTools(s, svcs)
	registerLinkTools(s, svcs)

	return server.NewSSEServer(s, server.WithBaseURL(baseURL))
}
