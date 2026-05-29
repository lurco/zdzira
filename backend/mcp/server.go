package mcp

import (
	"net/http"
	"zdzira/backend/service"

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

func NewHandler(svcs *service.Services) http.Handler {
	return server.NewStreamableHTTPServer(NewServer(svcs))
}
