package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"zdzira/backend/api"
	zdziramcp "zdzira/backend/mcp"
	"zdzira/backend/service"
	"zdzira/backend/store"

	"github.com/go-chi/chi/v5"
)

func main() {
	dbPath := flag.String("db", "zdzira.db", "SQLite database path")
	addr := flag.String("addr", ":8080", "Listen address")
	flag.Parse()

	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	db, err := store.Open(*dbPath)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}

	stores := store.New(db)
	svcs := service.New(stores)

	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = fmt.Sprintf("http://localhost%s", *addr)
	}

	r := chi.NewRouter()
	r.Mount("/", api.NewRouter(svcs, logger))
	r.Mount("/mcp", zdziramcp.NewSSEHandler(svcs, baseURL))

	logger.Info("starting", "addr", *addr, "db", *dbPath, "docs", baseURL+"/docs")
	if err := http.ListenAndServe(*addr, r); err != nil {
		log.Fatalf("server: %v", err)
	}
}
