package mcp

import (
	"net/http"
	"zdzira/internal/service"

	"github.com/mark3labs/mcp-go/server"
)

type Server = server.MCPServer

func NewServer(svcs *service.Services) *Server {
	s := server.NewMCPServer("zdzira", "1.0.0",
		server.WithToolCapabilities(true),
	)

	registerProjectTools(s, svcs)
	registerEpicTools(s, svcs)
	registerUpdateEpicTools(s, svcs)
	registerIssueTools(s, svcs)
	registerUpdateIssueTools(s, svcs)
	registerDeleteIssueTools(s, svcs)
	registerSwimlaneTools(s, svcs)
	registerCommentTools(s, svcs)
	registerLinkTools(s, svcs)

	return s
}

func NewSSEHandler(svcs *service.Services, baseURL string) http.Handler {
	return server.NewSSEServer(NewServer(svcs), server.WithBaseURL(baseURL))
}
