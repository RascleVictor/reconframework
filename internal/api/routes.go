package api

import (
	"net/http"

	"github.com/RascleVictor/recon-engine/internal/api/handlers"
	"github.com/RascleVictor/recon-engine/internal/observation"
	"github.com/RascleVictor/recon-engine/internal/plugin/executor"
	"github.com/RascleVictor/recon-engine/internal/plugin/registry"
	"github.com/RascleVictor/recon-engine/internal/scan"
)

func NewRouter(
	pluginRegistry *registry.Registry,
	dockerExecutor executor.Executor,
) http.Handler {
	mux := http.NewServeMux()

	pluginHandler := handlers.NewPluginHandler(
		pluginRegistry,
		dockerExecutor,
	)

	observationRepository := observation.NewMemoryRepository()

	scanService := scan.NewService(
		dockerExecutor,
		observationRepository,
	)

	scanHandler := handlers.NewScanHandler(
		pluginRegistry,
		scanService,
	)

	mux.HandleFunc("GET /health", handlers.Health)
	mux.HandleFunc("GET /api/v1/plugins", pluginHandler.List)
	mux.HandleFunc("GET /api/v1/plugins/{id}", pluginHandler.Get)
	mux.HandleFunc("POST /api/v1/plugins/{id}/execute", pluginHandler.Execute)

	mux.HandleFunc(
		"POST /api/v1/scans/{scan_id}/plugins/{id}",
		scanHandler.Run,
	)

	return mux
}
