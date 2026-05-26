package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"zdzira/internal/api"
	zdziramcp "zdzira/internal/mcp"
	"zdzira/internal/service"
	"zdzira/internal/store"

	"github.com/go-chi/chi/v5"
)

func main() {
	dbPath := flag.String("db", "zdzira.db", "SQLite database path")
	addr := flag.String("addr", ":8080", "Listen address")
	flag.Parse()

	db, err := store.Open(*dbPath)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}

	stores := store.New(db)
	svcs := service.New(stores)

	baseURL := fmt.Sprintf("http://localhost%s", *addr)

	r := chi.NewRouter()
	r.Mount("/", api.NewRouter(svcs))
	r.Mount("/mcp", zdziramcp.NewSSEHandler(svcs, baseURL))

	log.Printf("zdzira listening on %s (REST + MCP SSE at /mcp)", *addr)
	if err := http.ListenAndServe(*addr, r); err != nil {
		log.Fatalf("server: %v", err)
	}
}
