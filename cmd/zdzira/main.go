package main

import (
	"flag"
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
	dumpOpenAPI := flag.Bool("dump-openapi", false, "Print the OpenAPI spec to stdout and exit")
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

	if *dumpOpenAPI {
		spec, err := api.OpenAPISpec(svcs)
		if err != nil {
			log.Fatalf("dump openapi: %v", err)
		}
		os.Stdout.Write(spec)
		return
	}

	broadcaster := api.NewBroadcaster()

	r := chi.NewRouter()
	r.Mount("/", api.NewRouter(svcs, logger, stores.Ping, broadcaster))
	r.Mount("/mcp", api.NotifyMiddleware(broadcaster)(zdziramcp.NewHandler(svcs)))

	logger.Info("starting", "addr", *addr, "db", *dbPath)
	if err := http.ListenAndServe(*addr, r); err != nil {
		log.Fatalf("server: %v", err)
	}
}
